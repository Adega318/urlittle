# syntax=docker/dockerfile:1

# ---- Build stage ----
FROM --platform=$BUILDPLATFORM docker.io/library/golang:1.27.1-alpine AS builder

ARG TARGETOS
ARG TARGETARCH
ARG VERSION=dev

WORKDIR /src

RUN --mount=type=bind,source=go.mod,target=go.mod \
  --mount=type=bind,source=go.sum,target=go.sum \
  --mount=type=cache,target=/go/pkg/mod \
  go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
  go build -trimpath -ldflags="-s -w -X main.version=${VERSION}" \
  -o /out/app ./cmd/urlittle

# ---- Runtime stage ----
FROM docker.io/library/alpine:3.24

RUN apk add --no-cache ca-certificates tzdata \
  && addgroup -S -g 65532 app \
  && adduser -S -H -u 65532 -G app app

COPY --from=builder /out/app /usr/local/bin/app

ENV PORT=8080
EXPOSE 8080

USER 65532:65532

ENTRYPOINT ["/usr/local/bin/app"]
