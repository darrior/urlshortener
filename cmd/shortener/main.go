package main

import (
	"fmt"
	"os"

	"github.com/darrior/urlshortener/cmd/shortener/config"
	"github.com/darrior/urlshortener/internal/handler"
	"github.com/darrior/urlshortener/internal/repository"
	"github.com/darrior/urlshortener/internal/service"
)

func main() {
	c := config.ParseConfig()

	r := repository.NewMapRepository()

	s := service.NewService(r, c.BaseAddress)

	srv := handler.NewServer(c.ListenAddress, s)
	if err := srv.Run(); err != nil {
		fmt.Printf("An error occured: %s\n", err.Error())
		os.Exit(1)
	}
}
