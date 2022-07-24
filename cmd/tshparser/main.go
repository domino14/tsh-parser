package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/justinas/alice"
	_ "github.com/mattn/go-sqlite3"
	"github.com/namsral/flag"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"

	"github.com/domino14/tshparser/pkg/parser"
	proto "github.com/domino14/tshparser/rpc"
)

type Config struct {
	DBMigrationsPath  string
	DBPath            string
	TourneySchemaPath string
	SecretKey         string
}

func (c *Config) Load(args []string) error {
	fs := flag.NewFlagSet("tshparser", flag.ContinueOnError)
	fs.StringVar(&c.DBMigrationsPath, "db-migrations-path", "", "the path where migrations are stored")
	fs.StringVar(&c.DBPath, "db-path", "", "the path of the sqlite3 database")
	fs.StringVar(&c.TourneySchemaPath, "tourney-schema-path", "", "the path of the tournament schema, with points/division breakdowns")
	fs.StringVar(&c.SecretKey, "secret-key", "", "a secret key for signing stuff")
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

	if cfg.SecretKey == "" {
		panic("no secret key provided")
	}

	log.Info().Interface("config", cfg).Msg("loaded-config")
	ensureMigrations(cfg)

	router := http.NewServeMux()

	middlewares := alice.New(
		hlog.NewHandler(log.With().Str("service", "liwords").Logger()),
		hlog.AccessHandler(func(r *http.Request, status int, size int, d time.Duration) {
			path := strings.Split(r.URL.Path, "/")
			method := path[len(path)-1]
			hlog.FromRequest(r).Info().Str("method", method).Int("status", status).Dur("duration", d).Msg("")
		}),
		parser.ExposeResponseWriterMiddleware,
		parser.AuthenticationMiddleware,
		parser.JWTMiddleware,
	)

	store, err := parser.NewSqliteStore(cfg.DBPath)
	if err != nil {
		panic(err)
	}

	tourneyservice := parser.NewService(store, cfg.TourneySchemaPath, cfg.SecretKey)
	authservice := parser.NewAuthService(store, cfg.SecretKey)

	router.Handle(proto.TournamentRankerServicePathPrefix,
		middlewares.Then(proto.NewTournamentRankerServiceServer(tourneyservice)))

	router.Handle(proto.AuthenticationServicePathPrefix,
		middlewares.Then(proto.NewAuthenticationServiceServer(authservice)))

	srv := &http.Server{
		Addr:         ":8082",
		Handler:      router,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second}

	idleConnsClosed := make(chan struct{})
	sig := make(chan os.Signal, 1)

	// router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/static"))))
	router.Handle("/", http.FileServer(http.Dir("./ui")))

	go func() {
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		log.Info().Msg("got quit signal...")
		ctx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := srv.Shutdown(ctx); err != nil {
			// Error from closing listeners, or context timeout:
			log.Error().Msgf("HTTP server Shutdown: %v", err)
		}
		shutdownCancel()
		close(idleConnsClosed)
	}()
	log.Info().Msg("starting listening...")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("")
	}
	<-idleConnsClosed
	log.Info().Msg("server gracefully shutting down")
}
