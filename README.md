## Prerequisites

The only requirements to build and use this project are Docker and `make`. The
latter can easily be substituted with your scripting tool of choice.

## Getting started

Building the project will output a static binary in the bin/ folder. The
 platform can be changed using the `PLATFORM` variable:
```console
$ make                        # build for your host OS
$ make PLATFORM=darwin/amd64  # build for macOS
$ make PLATFORM=windows/amd64 # build for Windows x86_64
$ make PLATFORM=linux/amd64   # build for Linux x86_64
$ make PLATFORM=linux/arm     # build for Linux ARM
```

You can then run the binary as follows:

```console
$ ./bin/certificate-bootstrapper
```

To run the linter and unit tests:

```console
$ make test
```

## Containerized go development environment

The Dockerfile and Makefile in this project were inspired on: ![containerized-go-dev](https://github.com/chris-crone/containerized-go-dev/)
