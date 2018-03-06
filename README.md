# sharexserver [![License Apache 2.0](https://img.shields.io/badge/license-Apache--2.0-blue.svg?maxAge=31622400)](https://www.apache.org/licenses/LICENSE-2.0) [![stability-alpha](https://img.shields.io/badge/stability-alpha-f4d03f.svg)](https://github.com/mkenney/software-guides/blob/master/STABILITY-BADGES.md#alpha) [![GoDoc](https://godoc.org/github.com/mmichaelb/sharexserver?status.svg)](https://godoc.org/github.com/mmichaelb/sharexserver) [![Build Status](https://travis-ci.org/mmichaelb/sharexserver.svg?branch=dev)](https://travis-ci.org/mmichaelb/sharexserver) [![Go Report Card](https://goreportcard.com/badge/github.com/mmichaelb/sharexserver)](https://goreportcard.com/report/github.com/mmichaelb/sharexserver)
Lightweight upload server for the ShareX client (https://getsharex.com/).

# Description
This application can be used as a standalone server side endpoint for your ShareX client. It is written in Go and designed to be lightweight and easy to understand. If you are a Golang developer, you can also use this project as your dependency and use the code in your own project.

# Installation
## Getting the binaries
In order to install the ShareX server you have to get the binaries. There two possible methods of getting them:
- download a release file from the [GitHub releases page](https://github.com/mmichaelb/sharexserver/releases)
- compile the source manually on your own (see [Compilation](https://github.com/mmichaelb/sharexserver#compilation))
## Download default configuration files
In order to adjust values of the application's runtime, you should download the default configurations to get an orientation. The downloads can be found in the [config directory](https://github.com/mmichaelb/sharexserver/tree/master/configs). After downloading the configuration you should rename it and adjust the values according to the [TOML conventions](https://github.com/toml-lang/toml).
## Running the application
At the moment the only parameter which the application accepts on startup is -config - you can specify the path to your configuration file. If you do not specify one, the default path (`./config`) is used. An example of running the application would be:
```bash
./your-executable -config=./my-custom-config.toml
```
Have fun and feel free to open up an issue if you have a problem with running your application. In the future, I hope that I can provide an auto-installation script or provide a custom Docker image.

# Compilation
The compilation of this code was successful with Go 1.8 and 1.9 - newer versions should normally work as well.

In general, there are two ways of building the application:
## Makefile
When using `make`, life is easy and you can just run:
```bash
make build
```
## go build command
When compiling with the standard `go build` command, you can use the extracted command from the Makefile. Because with `make` the ld flags are parsed automatically, you have to replace them on your own when running `go build` manually.
```bash
go build -ldflags '-X "main.applicationName=ShareX Server" -X "main.version=<version>" -X "main.branch=<branch>" -X "main.commit=<commit>"' ./cmd/sharexserver
```

# Using ShareX server as a dependency
To use this project as a dependency for your own project, you can just `go get` the `cmd/sharexserver` package:
```bash
go get -u github.com/mmichaelb/sharexserver/cmd/sharexserver
```
Make sure to check out the [examples package](https://github.com/mmichaelb/sharexserver/tree/master/examples/) for implemented examples and use cases.

# Contribution
Feel free to contribute and help this project to grow. You can also just suggest features/enhancements - for more details check the [contributing file](https://github.com/mmichaelb/sharexserver/tree/master/.github/CONTRIBUTING.md).
