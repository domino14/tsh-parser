package main

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/domino14/tshparser/pkg/parser"
	"github.com/domino14/tshparser/pkg/server"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"github.com/namsral/flag"
	"github.com/rs/zerolog/log"
)

type Config struct {
	DBMigrationsPath  string
	DBPath            string
	TourneySchemaPath string
}

func (c *Config) Load(args []string) error {
	fs := flag.NewFlagSet("tshparser", flag.ContinueOnError)
	fs.StringVar(&c.DBMigrationsPath, "db-migrations-path", "", "the path where migrations are stored")
	fs.StringVar(&c.DBPath, "db-path", "", "the path of the sqlite3 database")
	fs.StringVar(&c.TourneySchemaPath, "tourney-schema-path", "", "the path of the tournament schema, with points/division breakdowns")
	err := fs.Parse(args)
	return err
}

func ensureMigrations(cfg *Config) {
	sqliteDb, err := sql.Open("sqlite3", cfg.DBPath)
	if err != nil {
		panic(err)
	}
	driver, err := sqlite3.WithInstance(sqliteDb, &sqlite3.Config{})
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(cfg.DBMigrationsPath, cfg.DBPath, driver)
	if err != nil {
		panic(err)
	}
	log.Info().Msg("bringing up migration")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		panic(err)
	}
	e1, e2 := m.Close()
	log.Err(e1).Msg("close-source")
	log.Err(e2).Msg("close-database")
}

func main() {
	cfg := &Config{}
	cfg.Load(os.Args[1:])
	log.Info().Interface("config", cfg).Msg("loaded-config")
	ensureMigrations(cfg)

	store, err := parser.NewSqliteStore(cfg.DBPath)
	if err != nil {
		panic(err)
	}
	service := parser.NewService(store, cfg.TourneySchemaPath)
	srv := server.NewHTTPServer(service)

	http.Handle("/api", srv)
	http.ListenAndServe(":8082", nil)

}
