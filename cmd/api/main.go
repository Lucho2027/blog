package main

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lucho2027/blog/internal/api"
	"github.com/lucho2027/blog/internal/database"
	"github.com/rs/zerolog"
)

const (
	maxConn = 10
)

func main() {
	env := os.Getenv("ENV")

	if env == "" {
		env = "dev"
	}
	logger := loggerInit(env, "api")
	err := run(logger)
	if err != nil {
		os.Exit(1)
	}
}

func run(l zerolog.Logger) error {
	l.Debug().Msg("starting up")
	dbUrl := os.Getenv("DATABASE_URL")

	if dbUrl == "" {
		const msg = "DATABASE_URL is not set"
		err := errors.New(msg)
		l.Err(err).Msg(msg)
		return err
	}
	dbConfig := database.DBConfig{
		DatabaseURL: dbUrl,
		MaxConn:     maxConn,
	}
	l.Debug().Msg("loading config")
	startupContext := context.Background()
	ctx, cancel := context.WithTimeout(startupContext, time.Second*5)
	defer cancel()
	l.Debug().Msg("starting db")
	db, err := database.New(ctx, dbConfig)
	if err != nil {
		l.Err(err).Msg("db unable to load")
		return err
	}
	l.Debug().Msg("db connected")
	defer db.Close()

	srv := api.NewServer(db, l)
	handler := srv.Handler()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	s := &http.Server{
		Handler: handler,
		Addr:    ":" + port,
	}

	l.Debug().Str("addr", s.Addr).Msg("listening on port")

	sigCh := make(chan os.Signal, 1)

	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	servErr := make(chan error, 1)
	go func() {
		err := s.ListenAndServe()
		servErr <- err
	}()

	select {
	case sig := <-sigCh:
		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		l.Info().Interface("signal", sig).Msg("shutting down")
		err := s.Shutdown(shutdownCtx)
		if err != nil {
			l.Err(err).Msg("problem shutting down")
			return err
		}
		l.Info().Interface("signal", sig).Msg("shutting down success")
		return nil

	case err := <-servErr:
		if err == http.ErrServerClosed {
			l.Info().Msg("server stopped")
			return nil
		} else if err != nil {
			l.Err(err).Msg("http server failed")
			return err
		}
	}
	return nil
}

func loggerInit(env, service string) zerolog.Logger {
	var writer io.Writer

	if env == "dev" || env == "local" {
		writer = zerolog.ConsoleWriter{
			Out: os.Stdout,
		}
	} else {
		writer = os.Stdout
	}
	var logLevel zerolog.Level
	if env == "prod" {
		logLevel = zerolog.InfoLevel
	} else {
		logLevel = zerolog.DebugLevel
	}
	zerologInstance := zerolog.New(writer).Level(logLevel).With().Str("service", service).Timestamp().Logger()

	return zerologInstance
}
