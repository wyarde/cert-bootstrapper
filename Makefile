PLATFORM=local
export DOCKER_BUILDKIT=1

.PHONY: all
all: certificate-bootstrapper agent-linux agent-windows 

.PHONY: lint
lint: lint-certificate-bootstrapper lint-agent

# Certificate bootstrapper depends on the agent
certificate-bootstrapper: agent-linux agent-windows 

.PHONY: certificate-bootstrapper
certificate-bootstrapper:
	@docker build . --target bin \
	--output bin/ \
	--platform ${PLATFORM} \
	--build-arg COMPONENT=certificate-bootstrapper

.PHONY: agent-linux
agent-linux:
	@docker build . --target bin \
	--output cmd/certificate-bootstrapper/bin/ \
	--platform linux \
	--build-arg COMPONENT=agent

.PHONY: agent-windows
agent-windows:
	@docker build . --target bin \
	--output cmd/certificate-bootstrapper/bin/ \
	--platform windows \
	--build-arg COMPONENT=agent

.PHONY: lint-certificate-bootstrapper
lint-certificate-bootstrapper:
	@docker build . --target lint \
	--build-arg COMPONENT=certificate-bootstrapper

.PHONY: lint-agent
lint-agent:
	@docker build . --target lint \
	--build-arg COMPONENT=agent
