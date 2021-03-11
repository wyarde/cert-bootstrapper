all: bin/certificate-bootstrapper

PLATFORM=local
export DOCKER_BUILDKIT=1

.PHONY: bin/certificate-bootstrapper
bin/certificate-bootstrapper:
	@docker build . --target bin \
	--output bin/ \
	--platform ${PLATFORM}

.PHONY: lint
lint:
	@docker build . --target lint
