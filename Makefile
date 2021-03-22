define get_docker_server_os
$(shell docker version -f "{{.Server.Os}}" | sed "s/./\U&/")
endef

ifeq ($(OS),Windows_NT)
	export PLATFORM=Windows
	ORIGINAL_DOCKER_SERVER_OS := $(call get_docker_server_os)
else
	export PLATFORM=Linux
	ORIGINAL_DOCKER_SERVER_OS := Linux
endif

.PHONY: image
image: image-$(PLATFORM)

.PHONY: image-Linux
image-Linux:
	$(MAKE) switchToLinux
	docker buildx build . \
		--target runtime-Linux \
	  --platform Linux \
	  --tag wyarde/cert-bootstrapper:latest
	$(MAKE) switchBack ORIGINAL_DOCKER_SERVER_OS=${ORIGINAL_DOCKER_SERVER_OS}

.PHONY: image-Windows
image-Windows: bin
	$(MAKE) switchToWindows
	docker build . \
		--target runtime-Windows \
		--tag wyarde/cert-bootstrapper:latest
	$(MAKE) switchBack ORIGINAL_DOCKER_SERVER_OS=${ORIGINAL_DOCKER_SERVER_OS}

image-%:
	@echo "Invalid platform specified: '$*'. Use either 'Linux' or 'Windows' (case sensitive!)"

.PHONY: bin
bin:
	$(MAKE) switchToLinux
	docker buildx build . \
	--target bin \
	--platform $(PLATFORM) \
	--output bin/
	$(MAKE) switchBack ORIGINAL_DOCKER_SERVER_OS=${ORIGINAL_DOCKER_SERVER_OS}

switchTo%:
ifeq ($(OS),Windows_NT)
	$(if $(subst $(call get_docker_server_os),,$*),\
		@echo Switching Docker engine to $*... && \
		'C:\Program Files\Docker\Docker\DockerCli.exe' -Switch$*Engine && \
		sleep 1\
		,\
		@echo No need to switch\
	)
else
	$(if $(subst $(call get_docker_server_os),,$*),\
		$(error Please switch Docker Engine to $* first) \
	)
	@echo
endif

switchBack:
ifeq ($(OS),Windows_NT)
	$(if $(subst $(ORIGINAL_DOCKER_SERVER_OS),,$(call get_docker_server_os)),\
		@echo Switching back Docker engine to $(ORIGINAL_DOCKER_SERVER_OS)... && \
		'C:\Program Files\Docker\Docker\DockerCli.exe' -Switch$(ORIGINAL_DOCKER_SERVER_OS)Engine && \
		sleep 1\
		,\
		@echo No need to switch back\
	)
else
# Do nothing in Linux
	@echo 
endif