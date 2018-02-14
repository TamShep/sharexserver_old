APPLICATION_NAME = ShareX Server
VERSION = 0.1.0
BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
COMMIT = $(shell git rev-parse HEAD)

LD_FLAGS = -X "main.applicationName=${APPLICATION_NAME}" -X "main.version=${VERSION}" -X "main.branch=${BRANCH}" -X "main.commit=${COMMIT}"

# general go get function
define goget
	@echo -ne "go get >> $1\r"
	@go get -t github.com/
	@echo "go get << $1"
endef

# "go gets" all needed dependencies with the default built-in standard tool
deps:
# 	MongoDB Golang library
	$(call goget,gopkg.in/mgo.v2)
#	Powerful router and dispatcher for Golang
	$(call goget,github.com/gorilla/mux)

# formats the *.go files with the built-in go fmt tool
format:
	@go fmt ./...

# builds and formats the project with the built-in Golang tool
build: format
	@go build -ldflags '${LD_FLAGS}' ./cmd/sharexserver

# installs and formats the project with the built-in Golang tool
install: format
	@go install -ldflags '${LD_FLAGS}' ./cmd/sharexserver
