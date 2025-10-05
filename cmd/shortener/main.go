package main

import (
	"fmt"
	"os"

	"github.com/darrior/urlshortener/internal/handler"
	"github.com/darrior/urlshortener/internal/service"
)

func main() {
	s := service.NewService()
	srv := handler.NewServer("127.0.0.1:8080", s)
	if err := srv.Run(); err != nil {
		fmt.Printf("An error occured: %s\n", err.Error())
		os.Exit(1)
	}
}
