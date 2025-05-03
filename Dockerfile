FROM golang:1.22.0-alpine

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o main ./cmd/main.go

ENTRYPOINT ["./main"]