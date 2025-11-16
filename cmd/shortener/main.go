package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/darrior/urlshortener/cmd/shortener/config"
	"github.com/darrior/urlshortener/internal/handler"
	"github.com/darrior/urlshortener/internal/repository"
	"github.com/darrior/urlshortener/internal/service"
	_ "github.com/jackc/pgx/v5"
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

	f, err := os.OpenFile(c.StorageFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal().Err(err).Msg("Can not open storage file")
	}

	r, err := repository.NewFSRepository(f)
	if err != nil {
		log.Fatal().Err(err).Msg("Can not initialize repository")
		os.Exit(-1)
	}

	defer func() {
		if err := r.Close(); err != nil {
			log.Error().Err(err).Msg("Can not close repository")
		}
	}()

	s := service.NewService(r, c.BaseAddress.String())

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
