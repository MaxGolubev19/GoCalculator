FROM golang:1.23.1

WORKDIR /app

COPY go.mod ./
RUN go mod tidy

COPY . .

RUN go build -o orchestrator ./cmd/orchestrator
RUN go build -o agent ./cmd/agent