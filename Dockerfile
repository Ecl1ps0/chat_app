FROM golang:1.22 AS builder
WORKDIR /app

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go env

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download && go mod verify

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o main ./cmd/main.go

FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y libc6

WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8080
CMD ["./main"]