package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/darrior/urlshortener/cmd/shortener/config"
	"github.com/darrior/urlshortener/internal/handler"
	"github.com/darrior/urlshortener/internal/repository"
	"github.com/darrior/urlshortener/internal/service"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
)

func main() {
	c, err := config.ParseConfig()
	if err != nil {
		c = config.DefaultConfig()
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-sig
		cancel()
	}()

	r, err := initRepository(c)
	if err != nil {
		log.Fatal().Err(err).Msg("Can not initialize repository")
		os.Exit(1)
	}
	defer func() {
		if err := r.Close(); err != nil {
			log.Error().Err(err).Msg("Can not close repository")
		}
	}()

	s := service.NewService(r, c.BaseAddress.String(), c.AuthKey)

	srv := handler.NewServer(string(c.ListenAddress), s)
	go func() {
		if err := srv.Stop(ctx); err != nil {
			log.Error().Err(err).Msg("Can not stop server properly")
			return
		}
		log.Info().Msg("Shutting down server gracefuly")
	}()

	if err := srv.Run(); err != nil {
		log.Error().Err(err).Msg("Unexpected server error")
		os.Exit(1)
	}
}

func initRepository(cfg config.Config) (repository.Repository, error) {
	if cfg.DatabaseDSN != nil {
		db, err := sql.Open("pgx", cfg.DatabaseDSN.ConnString())
		if err != nil {
			return nil, fmt.Errorf("can not open db connection: %w", err)
		}
		r, err := repository.NewDBRepository(db)
		if err != nil {
			return nil, fmt.Errorf("can not create DBRepository: %w", err)
		}

		return r, nil
	}

	if cfg.StorageFile != "" {
		f, err := os.OpenFile(cfg.StorageFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			return nil, fmt.Errorf("can not open storage file: %w", err)
		}

		r, err := repository.NewFSRepository(f)
		if err != nil {
			return nil, fmt.Errorf("can not create FSRepository: %w", err)
		}

		return r, nil
	}

	return repository.NewMapRepository(), nil
}
