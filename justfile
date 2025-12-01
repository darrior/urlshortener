# Build urlshortener binary.
build:
    go build -o build/urlshortener cmd/shortener/main.go

# Run urlshortener.
run *args:
    go run cmd/shortener/main.go {{args}}
    
# Run unit-tests only
test path="./...":
    go test {{path}}

# Calculate test coverage.
cover path="./...":
    @go test ./... -coverprofile /tmp/cover.out > /dev/null && go tool cover -func /tmp/cover.out
    @rm /tmp/cover.out

cover-html:
    @go test ./... -coverprofile /tmp/cover.out > /dev/null && go tool cover -html /tmp/cover.out
    @rm /tmp/cover.out

# Run golangci-lint.
lint path="./...":
    golangci-lint run {{path}}

# Run linter and unit-tests.
check: lint test

# Generaet mock files from all intefaces.
generate-mocks: generate-repository-mock generate-service-mock

# Generate Repository mock.
generate-repository-mock: (_generate_mock "mock_repository.go" "internal/repository" "Repository")
    
# Generate Service mock.
generate-service-mock: (_generate_mock "mock_service.go" "internal/service" "IService")

# Common generate-XXX-mock implementation
_generate_mock dest-file package interface:
    go tool mockgen -destination=internal/mocks/{{dest-file}} -package=mocks github.com/darrior/urlshortener/{{package}} {{interface}}
