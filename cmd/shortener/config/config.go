// Package config implements parsing configuration options from command line.
package config

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"strings"
)

const (
	_defaultListenAddress = "127.0.0.1:8080"
	_defaultBaseAddress   = "http://127.0.0.1:8080"
)

var (
	errorValidateListenAddress = errors.New("listen address must be in form host:port")
	errorValidateBaseAddress   = errors.New("invalid base address")
)

type Config struct {
	ListenAddress string
	BaseAddress   string
}

func ParseConfig() Config {
	c := Config{
		ListenAddress: _defaultListenAddress,
		BaseAddress:   _defaultBaseAddress,
	}
	flag.Func("a", "listen address for web-server", c.validateListenAddress)
	flag.Func("b", "base address for short URL", c.validateBaseAddress)
	flag.Parse()

	return c
}

func (c *Config) validateListenAddress(address string) error {
	if address == "" {
		c.ListenAddress = _defaultListenAddress
		return nil
	}

	splitedAddress := strings.Split(address, ":")
	if len(splitedAddress) != 2 {
		return errorValidateListenAddress
	}

	c.ListenAddress = address

	return nil
}

func (c *Config) validateBaseAddress(address string) error {
	if address == "" {
		c.BaseAddress = _defaultBaseAddress
		return nil
	}

	parsedURL, err := url.Parse(address)
	if err != nil {
		return fmt.Errorf("%s: %w", errorValidateBaseAddress.Error(), err)
	}

	parsedURL.Path = ""
	parsedURL.RawQuery = ""

	c.BaseAddress = parsedURL.String()

	return nil
}
