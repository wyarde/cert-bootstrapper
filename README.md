# Certificate Bootstrapper

The certificate bootstrapper will monitor for new Docker containers, and then bootstrap them with a custom CA certificate.

## Prerequisites

The only requirements to build and use this project are Docker and `make`. The
latter can easily be substituted with your scripting tool of choice.

## Getting started

### Linux

You can run the certificate bootstrapper as follows:

```shell
docker run -d --restart unless-stopped \
  -v /var/run.docker.sock:/var/run/docker.sock \
  -v /path/to/my_cert/my_cert.pem:/ssl/cert.pem \
  wyarde/cert-bootstrapper
```

### Windows

Running a Linux container on Windows still requires some trickery. Follow following steps:

1. Enable experimental mode to allow Linux Containers on Windows (LCOW)
2. Get ip address of the docker nat interface

```shell
  docker network inspect -f '{{range .IPAM.Config}}{{.Gateway}}{{end}}' nat
```

3. Add configuration in `daemon.json` to listen to the docker nat interface address
   docker run --rm --platform=linux -e DOCKER_HOST=tcp://<docker_nat_network_ip>:2375 -v c:/tmp/:/ssl/ wyarde/cert-bootstrapper ssl/cert.pem

4. Start the bootstrapper

```shell
docker run --platform=linux -d --restart unless-stopped \
 -e DOCKER_HOST=tcp://<docker_nat_network_ip>:2375
-v c:/path/to/my_cert_pem/:/ssl/ \
 wyarde/cert-bootstrapper
```

## Build it yourself

To build the Docker image yourself:

```shell
make
```

If needed, the build can also output the binary:

```shell
make bin
```

To run the linter:

```shell
make lint
```

```

```
