package parser

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteStore struct {
	db *sql.DB
}

func NewSqliteStore(dbName string) (*SqliteStore, error) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, err
	}
	return &SqliteStore{db: db}, nil
}

func (s *SqliteStore) AddTournament(ctx context.Context,
	ttype TournamentType, name string, date time.Time,
	tfileContents []byte) (int64, error) {

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	datestr := date.Format(time.RFC3339)
	// I can technically insert the date directly into sqlite
	// but then it removes the `T` from the format string.
	res, err := tx.ExecContext(ctx, `
		INSERT INTO tournaments(type, name, date, contents)
		VALUES(?, ?, ?, ?)
	`, ttype, name, datestr, tfileContents)
	if err != nil {
		return 0, err
	}
	// later parse out players from tfile and insert those etc.
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *SqliteStore) RemoveTournament(ctx context.Context, id int) error {
	res, err := s.db.ExecContext(ctx, `
		DELETE FROM tournaments
		WHERE id = ?
	`, id)
	if err != nil {
		return err
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if ra == 0 {
		return errors.New("no tournament was found with that id")
	}
	return nil
}

func (s *SqliteStore) GetTournaments(ctx context.Context, begin, end time.Time) ([]Tournament, error) {

	rows, err := s.db.QueryContext(ctx, `
		SELECT id, type, name, date, contents
		FROM tournaments
		WHERE date >= ? AND date <= ?
	`, begin, end)
	if err != nil {
		return nil, err
	}
	tourneys := []Tournament{}
	defer rows.Close()
	for rows.Next() {
		t := Tournament{}
		var sdate string
		err = rows.Scan(&t.ID, &t.TType, &t.Name, &sdate, &t.Contents)
		if err != nil {
			return nil, err
		}
		t.Date, err = time.Parse(time.RFC3339, sdate)
		if err != nil {
			return nil, err
		}

		tourneys = append(tourneys, t)
	}
	return tourneys, nil
}

func (s *SqliteStore) AddPlayerAlias(ctx context.Context, origPlayer, alias string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// An original name should not be allowed to be an alias to
	// another original name.
	var count int
	err = s.db.QueryRowContext(ctx, `
		SELECT count(*) FROM player_aliases WHERE alias = ?
	`, origPlayer).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("the original player name %s is already an alias of another name", origPlayer)
	}
	_, err = tx.ExecContext(ctx, `
		INSERT INTO player_aliases(original_player, alias)
		VALUES(?, ?)
	`, origPlayer, alias)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (s *SqliteStore) RemovePlayerAlias(ctx context.Context, origPlayer, alias string) error {
	res, err := s.db.ExecContext(ctx, `
		DELETE FROM player_aliases WHERE original_player = ? AND alias = ?
	`, origPlayer, alias)
	if err != nil {
		return err
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if ra == 0 {
		return errors.New("that alias was not found")
	}
	return nil
}

func (s *SqliteStore) GetAllAliases(ctx context.Context) (map[string]string, error) {
	aliases := map[string]string{} // map of alias to player

	rows, err := s.db.QueryContext(ctx, `
		SELECT alias, original_player FROM player_aliases
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var alias, origPlayer string
		err = rows.Scan(&alias, &origPlayer)
		if err != nil {
			return nil, err
		}
		aliases[alias] = origPlayer
	}
	return aliases, nil
}
