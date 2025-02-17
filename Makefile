
.DEFAULT_GOAL := build

brew:
	brew install golangci-lint
	brew install staticcheck
	brew install gofumpt

install:
	go install github.com/google/addlicense@latest

clean:
	rm -rf dist/
	rm -rf tmp/
	rm -f coverage.out
	rm -f result.json

init-dependency:
	go get -u github.com/davecgh/go-spew
	go get -u github.com/jinzhu/copier
	go get -u github.com/google/uuid
	go get -u github.com/pelletier/go-toml"
	go get -u github.com/go-playground/validator/v10
	go get -u github.com/stretchr/testify

mod:
	go mod download
	go mod tidy

check:
	staticcheck  ./...

lint:
	go vet ./...
	gofmt -s -w **/**.go
	gofumpt -l -w .
	golangci-lint run --disable-all --enable staticcheck

lint-fix:
	gofmt -s -w **/**.go
	go vet ./...
	gofumpt -l -w .
	golangci-lint run ./... --fix

test:
	go test ./...

teste2e:
	export E2E="TRUE" && GOFLAGS="-count=1" go test ./e2e/...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out

coverage-plugin:
	go test -coverprofile=coverage.out ./plugin/...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out

converage-%:
	go test -coverprofile=coverage.out ./...

converage-json:
	go test -json -coverprofile=coverage.out ./... > result.json

run-release:
	go run ./cmd/server-all-in-one

build:  clean mod

run:  clean mod lint-fix run-release

# disallow any parallelism (-j) for Make. This is necessary since some
# commands during the build process create temporary files that collide
# under parallel conditions.
.NOTPARALLEL:

.PHONY: clean mod lint lint-fix release all
