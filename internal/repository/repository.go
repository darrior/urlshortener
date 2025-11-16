// Package repository provides interface over data sources.
package repository

import (
	"errors"
)

type Repository interface {
	AddURL(id, url string) (err error)
	GetURL(id string) (url string, err error)
	Close() (err error)
}

var ErrorNotFound = errors.New("not found")
