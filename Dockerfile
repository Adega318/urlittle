# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.27.1 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
  go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/app .

FROM alpine:3.24

RUN addgroup -S -g 65532 app \
  && adduser -S -D -H -u 65532 -G app app \
  && apk add --no-cache wget

COPY --from=builder /out/app /usr/local/bin/app

ENV PORT=8080
EXPOSE 8080

USER app

ENTRYPOINT ["/usr/local/bin/app"]
