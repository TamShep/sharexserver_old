APPLICATION_NAME = ShareX Server
VERSION = 0.3.0
BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
COMMIT = $(shell git rev-parse HEAD)

LD_FLAGS = -X "main.applicationName=${APPLICATION_NAME}" -X "main.version=${VERSION}" -X "main.branch=${BRANCH}" -X "main.commit=${COMMIT}"

# builds and formats the project with the built-in Golang tool
build:
	@go build -ldflags '${LD_FLAGS}' ./cmd/sharexserver

# installs and formats the project with the built-in Golang tool
install:
	@go install -ldflags '${LD_FLAGS}' ./cmd/sharexserver

# tests the project by running all test go files
test:
	@go test -race $(go list ./... | grep -v /vendor/ | grep -v /cmd/)
