FROM golang:1.21 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go build -o srvr ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/srvr .
COPY --from=builder /app/.env .

EXPOSE 8787

CMD ["./srvr"]
