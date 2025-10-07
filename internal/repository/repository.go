// Package repository provides interface over data sources.
package repository

type Repository interface {
	AddURL(id, url string) (err error)
	GetURL(id string) (url string, err error)
}
