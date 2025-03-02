ARG GO_BUILDER=docker.io/library/golang:1.23
ARG BASE_IMAGE=docker.io/library/alpine:3.12

FROM ${GO_BUILDER} AS builder

ARG VERSION
ARG GOPROXY
WORKDIR /workspace
COPY . .

RUN GOPROXY=${GOPROXY} go mod download
RUN GOPROXY=${GOPROXY} CGO_ENABLED=0 go build -ldflags "-w -s -X github.com/linuxsuren/api-testing/pkg/version.version=${VERSION}\
    -X github.com/linuxsuren/api-testing/pkg/version.date=$(date +%Y-%m-%d)" -o atest-store-redis .

FROM ${BASE_IMAGE}

LABEL org.opencontainers.image.source=https://github.com/LinuxSuRen/atest-ext-store-redis
LABEL org.opencontainers.image.description="Redis database Store Extension of the API Testing."

COPY --from=builder /workspace/atest-store-redis /usr/local/bin/atest-store-redis

CMD [ "atest-store-redis" ]
