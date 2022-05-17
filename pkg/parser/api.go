package parser

import (
	"context"
	"time"
)

type TournamentType string

type Standing struct {
	PlayerName        string  `json:"player_name"`
	Points            int     `json:"points"`
	Wins              float32 `json:"wins"`
	Spread            int     `json:"spread"`
	TournamentsPlayed int     `json:"tournaments_played"`
}

type Tournament struct {
	ID        int            `json:"id"`
	TType     TournamentType `json:"ttype"`
	Name      string         `json:"name"`
	Date      time.Time      `json:"date"`
	Contents  []byte         `json:"contents"`
	Standings []Standing     `json:"standings"`
}

type Service struct {
	store      *SqliteStore
	schemaPath string
}

func NewService(store *SqliteStore, schemaPath string) *Service {
	return &Service{store: store, schemaPath: schemaPath}
}

// AddTournament adds a tournament of type ttype, with the given name
// and .t file contents. A .t file is a TSH data file.
// The date is used for computing YTD standings.
func (s *Service) AddTournament(ctx context.Context, ttype TournamentType, name string, date time.Time, contents []byte) (int64, error) {
	return s.store.AddTournament(ctx, ttype, name, date, contents)
}

// RemoveTournament deletes the tournament from the database.
func (s *Service) RemoveTournament(ctx context.Context, id int) error {
	return s.store.RemoveTournament(ctx, id)
}

// ComputeStandings computes the standings between a certain date range.
func (s *Service) ComputeStandings(ctx context.Context, beginDate time.Time, endDate time.Time) ([]Standing, error) {
	tourneys, err := s.store.GetTournaments(ctx, beginDate, endDate)
	if err != nil {
		return nil, err
	}

	return computeStandings(tourneys, s.schemaPath)
}

func (s *Service) GetTournaments(ctx context.Context, beginDate time.Time, endDate time.Time) ([]Tournament, error) {
	return s.store.GetTournaments(ctx, beginDate, endDate)
}

// AddPlayerAlias adds an alias to an original player. This can be used for cases
// where names might not exactly match across different tournaments.
func (s *Service) AddPlayerAlias(ctx context.Context, origPlayer, alias string) error {

	return nil
}

func (s *Service) RemovePlayerAlias(ctx context.Context, alias string) error {
	return nil
}
