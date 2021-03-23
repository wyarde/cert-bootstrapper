FROM mcr.microsoft.com/windows/nanoserver:1809-amd64 as runtime-Windows
COPY bin/cert-bootstrapper-Windows-x86_64.exe /
USER ContainerAdministrator
ENTRYPOINT "c:\cert-bootstrapper-Windows-x86_64.exe"

FROM --platform=${BUILDPLATFORM} golang:1.16 AS base
WORKDIR /project
ENV CGO_ENABLED=0
COPY go.* .
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download


FROM base AS build
ARG TARGETOS
ARG TARGETARCH
RUN --mount=source=src,target=src,rw \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o src/bin/agent ./src/cmd/agent && \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/cert-bootstrapper ./src && \
    cp src/bin/agent /out/

FROM golangci/golangci-lint:v1.31.0-alpine AS lint-base

FROM build AS lint
RUN --mount=source=src,target=src,rw \
    --mount=from=build,src=/out,target=/out \
    --mount=from=lint-base,src=/usr/bin/golangci-lint,target=/usr/bin/golangci-lint \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/.cache/golangci-lint \
    cp /out/agent src/bin/ && \
    golangci-lint run --timeout 10m0s ./src

FROM scratch AS bin-linux
COPY --from=build /out/cert-bootstrapper /cert-bootstrapper-Linux-x86_64

FROM scratch AS bin-windows
COPY --from=build /out/cert-bootstrapper /cert-bootstrapper-Windows-x86_64.exe

FROM bin-${TARGETOS} as bin

FROM bin-linux as runtime-linux
ENTRYPOINT [ "/cert-bootstrapper-Linux-x86_64" ]
