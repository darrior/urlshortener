// Package config implements parsing configuration options from command line.
package config

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/jackc/pgx/v5"
)

type host string

const (
	_defaultListenAddress  host   = "127.0.0.1:8080"
	_defaultStoragFilePath string = "urls.json"
)

var _defaultBaseAddress = url.URL{
	Scheme: "http",
	Host:   "127.0.0.1:8080",
}

var _defaultDatabaseDSN = func() *pgx.ConnConfig {
	conf, _ := pgx.ParseConfig("postgresql://postgres:postgres@localhost:5432/urls?sslmode=disable")
	return conf
}()

var (
	errorValidateListenAddress = errors.New("listen address must be in form host:port")
	errorValidateBaseAddress   = errors.New("invalid base address")
	errorValidateDatabaseDSN   = errors.New("invalid database DSN")
)

type Config struct {
	ListenAddress host            `env:"LISTEN_ADDRESS"`
	BaseAddress   url.URL         `env:"BASE_ADDRESS"`
	StorageFile   string          `env:"FILE_STORAGE_PATH"`
	DatabaseDSN   *pgx.ConnConfig `env:"DATADBASE_DSN"`
}

func DefaultConfig() Config {
	return Config{
		ListenAddress: _defaultListenAddress,
		BaseAddress:   _defaultBaseAddress,
		StorageFile:   _defaultStoragFilePath,
		DatabaseDSN:   _defaultDatabaseDSN,
	}
}

func ParseConfig() (Config, error) {
	c := DefaultConfig()

	set := flag.NewFlagSet("", flag.ContinueOnError)
	set.Func("a", "listen address for web-server", c.validateListenAddress)
	set.Func("b", "base address for short URL", c.validateBaseAddress)
	set.Func("f", "path to storage file", c.validateSorageFile)
	set.Func("d", "database DSN", c.validateDatabaseDSN)
	if err := set.Parse(os.Args[1:]); err != nil {
		return Config{}, err
	}

	options := env.Options{
		FuncMap: map[reflect.Type]env.ParserFunc{
			reflect.TypeFor[host]():            parseHostEnv,
			reflect.TypeFor[*pgx.ConnConfig](): parseDatabaseDSNEnv,
		},
	}

	if err := env.ParseWithOptions(&c, options); err != nil {
		fmt.Printf("Can not parse env: %s", err.Error())
		return Config{}, err
	}

	return c, nil
}

func (c *Config) validateListenAddress(address string) error {
	if address == "" {
		return nil
	}

	h, err := parseHost(address)
	if err != nil {
		return err
	}

	c.ListenAddress = h

	return nil
}

func (c *Config) validateBaseAddress(address string) error {
	if address == "" {
		return nil
	}

	parsedURL, err := url.Parse(address)
	if err != nil {
		return fmt.Errorf("%s: %w", errorValidateBaseAddress.Error(), err)
	}

	parsedURL.Path = ""
	parsedURL.RawQuery = ""

	c.BaseAddress = *parsedURL

	return nil
}

func (c *Config) validateSorageFile(file string) error {
	c.StorageFile = file
	return nil
}

func (c *Config) validateDatabaseDSN(dsn string) error {
	conf, err := pgx.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("%s: %w", errorValidateDatabaseDSN, err)
	}

	c.DatabaseDSN = conf
	return nil
}

func parseHostEnv(h string) (any, error) {
	return parseHost(h)
}

func parseDatabaseDSNEnv(dsn string) (any, error) {
	conf, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errorValidateDatabaseDSN, err)
	}

	return conf, nil
}

func parseHost(h string) (host, error) {
	splitedHost := strings.Split(h, ":")
	if len(splitedHost) != 2 {
		return "", errorValidateListenAddress
	}

	if number, err := strconv.Atoi(splitedHost[1]); err != nil || number > 65535 {
		return "", fmt.Errorf("%s: %w", errorValidateListenAddress.Error(), err)
	}

	return host(h), nil
}
