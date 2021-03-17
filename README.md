## Prerequisites

The only requirements to build and use this project are Docker and `make`. The
latter can easily be substituted with your scripting tool of choice.

## Getting started

You can run the certificate bootstrapper as follows:

```console
# linux
$ docker run -d --restart unless-stopped \
  -v /var/run.docker.sock:/var/run/docker.sock \
  -v /path/to/my_cert.pem:/cert.pem \
  wyarde/cert-bootstrapper

# windows
$ docker run -d --restart unless-stopped \
  -v //var/run/docker.sock:/var/run/docker.sock \
  -v /path/to/my_cert.pem:/cert.pem \
  wyarde/cert-bootstrapper
```

To build the Docker image yourself:
```console
$ make  
```

If needed, the build can also output linux and windows binaries in bin/:
```
$ make bins
```

To run the linter:

```console
$ make lint
```

## Containerized go development environment

The Dockerfile and Makefile in this project were inspired on: ![containerized-go-dev](https://github.com/chris-crone/containerized-go-dev/)
