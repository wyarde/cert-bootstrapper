export DOCKER_BUILDKIT=1

.PHONY: all
all: agents image

.PHONY: agents
agents: agent-linux agent-windows

.PHONY: lint
lint: lint-cert-bootstrapper lint-agent

bins: agents
image: agents

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

.PHONY: agent-linux
agent-linux:
	@docker build . \
		--file Dockerfile.agent \
		--target bin \
		--output cmd/cert-bootstrapper/bin/ \
		--platform linux

.PHONY: agent-windows
agent-windows:
	@docker build . \
		--file Dockerfile.agent \
		--target bin \
		--output cmd/cert-bootstrapper/bin/ \
		--platform windows

.PHONY: lint-cert-bootstrapper
lint-cert-bootstrapper:
	@docker build . \
		--file Dockerfile.cert-bootstrapper \
		--target lint

.PHONY: lint-agent
lint-agent:
	@docker build . \
		--file Dockerfile.agent \
		--target lint
