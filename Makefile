PLATFORM=local
export DOCKER_BUILDKIT=1

.PHONY: all
all: cert-bootstrapper agent-linux agent-windows 

.PHONY: lint
lint: lint-cert-bootstrapper lint-agent

# Certificate bootstrapper depends on the agent
cert-bootstrapper: agent-linux agent-windows 

.PHONY: cert-bootstrapper
cert-bootstrapper:
	@docker build . --target bin \
	--output bin/ \
	--platform ${PLATFORM} \
	--build-arg COMPONENT=cert-bootstrapper

.PHONY: agent-linux
agent-linux:
	@docker build . --target bin \
	--output cmd/cert-bootstrapper/bin/ \
	--platform linux \
	--build-arg COMPONENT=agent

.PHONY: agent-windows
agent-windows:
	@docker build . --target bin \
	--output cmd/cert-bootstrapper/bin/ \
	--platform windows \
	--build-arg COMPONENT=agent

.PHONY: lint-cert-bootstrapper
lint-cert-bootstrapper:
	@docker build . --target lint \
	--build-arg COMPONENT=cert-bootstrapper

.PHONY: lint-agent
lint-agent:
	@docker build . --target lint \
	--build-arg COMPONENT=agent
