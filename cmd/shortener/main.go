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
	c, err := config.ParseConfig()
	if err != nil {
		c = config.DefaultConfig()
	}

	r := repository.NewMapRepository()

	s := service.NewService(r, c.BaseAddress.String())

	srv := handler.NewServer(string(c.ListenAddress), s)
	if err := srv.Run(); err != nil {
		fmt.Printf("An error occured: %s\n", err.Error())
		os.Exit(1)
	}
}
