### builder
FROM golang:1.24-alpine AS builder
WORKDIR /build

ARG BUILD_PATH=./apps/gateway/cmd/main.go
ARG BINARY_NAME=gateway
ARG VERSION=dev
ARG CGO_ENABLED=0

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=${CGO_ENABLED} GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w -X main.version=${VERSION}" \
    -o /out/${BINARY_NAME} ${BUILD_PATH}

### runtime
FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /app
ARG BINARY_NAME=gateway
COPY --from=builder /out/${BINARY_NAME} /app/${BINARY_NAME}

USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/app/gateway"]