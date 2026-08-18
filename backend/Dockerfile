FROM golang:1.26-alpine3.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o tasker ./cmd/server

FROM alpine:3.24

RUN adduser -D tasker
USER tasker

WORKDIR /app

COPY --from=builder /app/tasker .

EXPOSE 8080

CMD ["./tasker"]
