# Build urlshortener binary.
build:
    go build -o build/urlshortener cmd/shortener/main.go

# Run urlshortener.
run *args:
    go run cmd/shortener/main.go {{args}}

# Run test DB.
start-db:
    podman run --publish 5432:5432 \
               --rm \
               --env POSTGRES_PASSWORD=123 \
               --detach \
               --name shortener-db \
               docker.io/library/postgres

# Stop test DB.
stop-db:
    podman stop shortener-db

# Run urlshortener with DB.
run-with-db: start-db && (run "-d" "postgres://postgres:123@localhost:5432/postgres?sslmode=disable")
    # HACK: time to postgres startup.
    @sleep 5
    
# Run unit-tests only
test path="./...":
    go test {{path}}

# Calculate test coverage.
cover path="./...":
    @go test ./... -coverprofile /tmp/cover.out > /dev/null && go tool cover -func /tmp/cover.out
    @rm /tmp/cover.out

# Calculate test coverage and create html page with visualization.
cover-html:
    @go test ./... -coverprofile /tmp/cover.out > /dev/null && go tool cover -html /tmp/cover.out
    @rm /tmp/cover.out

# Run golangci-lint.
lint path="./...":
    golangci-lint run {{path}}

# Run linter and unit-tests.
check: lint test

# Create new SQL migration.
[working-directory: 'migrations']
new-migration name:
    goose create {{name}} sql

# Generaet mock files from all intefaces.
generate-mocks: generate-repository-mock generate-service-mock

# Generate Repository mock.
generate-repository-mock: (_generate_mock "mock_repository.go" "internal/repository" "Repository")
    
# Generate Service mock.
generate-service-mock: (_generate_mock "mock_service.go" "internal/service" "IService")

# Common generate-XXX-mock implementation
_generate_mock dest-file package interface:
    go tool mockgen -destination=internal/mocks/{{dest-file}} -package=mocks github.com/darrior/urlshortener/{{package}} {{interface}}
