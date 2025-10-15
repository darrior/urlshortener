# URL Shortener

Learning project for Yandex Practicum cources.

## Build

```sh
go build -o urlshortener cmd/shortener/main.go
```

## Run tests

```sh
go test ./...
```

## Run linter

```sh
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.5.0
golangci-lint run
```
