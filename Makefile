export DOCKER_BUILDKIT=1

.PHONY: all
all: agents image

.PHONY: lint
lint: lint-cert-bootstrapper lint-agents

# Agent binaries are required for various tasks
bins: agents
image: agents
lint-cert-bootstrapper: agents

.PHONY: bins
bins:
	@docker build . \
		--file Dockerfile.cert-bootstrapper \
		--target bin \
		--output bin/ \
		--platform linux

	@docker build . \
		--file Dockerfile.cert-bootstrapper \
		--target bin \
		--output bin/ \
		--platform windows

.PHONY: image
image:
	@docker build . \
		--file Dockerfile.cert-bootstrapper \
		--target image \
		--platform linux \
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

.PHONY: lint-cert-bootstrapper
lint-cert-bootstrapper:
	@docker build . \
		--file Dockerfile.cert-bootstrapper \
		--target lint\
		--platform linux
	@docker build . \
		--file Dockerfile.cert-bootstrapper \
		--target lint \
		--platform linux

.PHONY: lint-agents
lint-agents:
	@docker build . \
		--file Dockerfile.agent \
		--target lint \
		--platform linux
