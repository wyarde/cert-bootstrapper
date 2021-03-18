export DOCKER_BUILDKIT=1

.PHONY: all
all: agents image

# Agent binaries are required for various tasks
bin: agents
image: agents

.PHONY: bin
bin:
	@docker build . \
		--file Dockerfile.cert-bootstrapper \
		--output bin/

.PHONY: image
image:
	@docker build . \
		--file Dockerfile.cert-bootstrapper \
		--tag wyarde/cert-bootstrapper:latest

.PHONY: agents
agents:
	@docker build . \
		--file Dockerfile.agent \
		--target bin \
		--output cmd/cert-bootstrapper/bin/ \
		--platform linux
	@docker build . \
		--file Dockerfile.agent \
		--target bin \
		--output cmd/cert-bootstrapper/bin/ \
		--platform windows
