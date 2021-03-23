# Certificate Bootstrapper

The certificate bootstrapper will monitor for new Docker containers, and then bootstrap them with a custom CA certificate.

## Getting started

See below for instructions on how to get started on both Linux and Windows.

### Linux

You can run the certificate bootstrapper as follows:

```shell
docker run -d --restart unless-stopped \
  -v /var/run.docker.sock:/var/run/docker.sock \
  -v /path/to/my_cert/my_cert.pem:/ssl/cert.pem \
  wyarde/cert-bootstrapper
```

### Windows

Make sure your certificate is named `cert.pem`. You can then run the certificate bootstrapper as follows:

```shell
docker run -d --restart unless-stopped --isolation process ``
 -v \\.\pipe\docker_engine:\\.\pipe\docker_engine ``
 -v c:/path/to/cert_pem/:c:/ssl/ ``
 wyarde/cert-bootstrapper
```

## Build it yourself

The only requirements to build and use this project are Docker, `make`, and `sed`.

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
