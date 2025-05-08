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

RUN apt-get update && apt-get install -y wget && \
    wget https://ftp.debian.org/debian/pool/main/g/glibc/libc6_2.34-0ubuntu3.1_amd64.deb && \
    dpkg -i libc6_2.34-0ubuntu3.1_amd64.deb && \
    rm libc6_2.34-0ubuntu3.1_amd64.deb

WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8080
CMD ["./main"]