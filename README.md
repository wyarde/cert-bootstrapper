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

Running a Linux container on Windows still requires some trickery.

**Warning**: The trickery below includes attaching the docker daemon to the docker nat interface, which means **any container** will be able to access it.
Follow following steps if you know what you're doing and still would like to proceed:

1. Get ip address of the docker nat interface

   ```shell
     docker network inspect -f '{{range .IPAM.Config}}{{.Gateway}}{{end}}' nat
   ```

2. Update/create `c:\ProgramData\Docker\config\daemon.json` to enable experimental mode to allow Linux Containers on Windows (LCOW) and listen to the docker nat interface address

   ```json
   {
     "experimental": true,
     "hosts": [
       "npipe:////./pipe/docker_engine",
       "tcp://<docker_nat_network_ip>"
     ]
   }
   ```

3. Rename your certificate to `cert.pem`

4. Start the bootstrapper

   ```shell
   docker run --platform=linux -d --restart unless-stopped ``
     -e DOCKER_HOST=tcp://<docker_nat_network_ip>:2375 ``
     -v c:/path/to/cert_pem/:/ssl/ ``
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

```

```
