package parser

import (
	"context"
	"database/sql"
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
	tfileContents []byte) error {

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `
		INSERT INTO tournaments(type, name, date, contents)
		VALUES(?, ?, ?, ?)
	`, ttype, name, date, tfileContents)
	if err != nil {
		return err
	}
	// later parse out players from tfile and insert those etc.
	err = tx.Commit()
	return err
}

func (s *SqliteStore) RemoveTournament(ctx context.Context, id int) error {
	_, err := s.db.ExecContext(ctx, `
		DELETE FROM tournaments
		WHERE id = ?
	`, id)
	return err
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
		err = rows.Scan(&t.ID, &t.TType, &t.Name, &t.Date, &t.Contents)
		if err != nil {
			return nil, err
		}
		tourneys = append(tourneys, t)
	}
	return tourneys, nil
}
