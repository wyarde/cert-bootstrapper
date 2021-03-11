FROM --platform=${BUILDPLATFORM} golang:1.16 AS base
WORKDIR /src
ENV CGO_ENABLED=0
COPY src/go.* .
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

FROM base AS build
ARG TARGETOS
ARG TARGETARCH
RUN --mount=source=src,target=. \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/certificate-bootstrapper .

FROM golangci/golangci-lint:v1.31.0-alpine AS lint-base

FROM base AS lint
RUN --mount=source=src,target=. \
    --mount=from=lint-base,src=/usr/bin/golangci-lint,target=/usr/bin/golangci-lint \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/.cache/golangci-lint \
    golangci-lint run --timeout 10m0s ./...

FROM scratch AS bin-unix
COPY --from=build /out/certificate-bootstrapper /certificate-bootstrapper-Linux-x86_64

FROM bin-unix AS bin-linux

FROM scratch AS bin-windows
COPY --from=build /out/certificate-bootstrapper /certificate-bootstrapper-Windows-x86_64.exe

FROM bin-${TARGETOS} as bin
