FROM golang:1.22.5-alpine as builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main cmd/main.go


FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

# Configuration files
RUN mkdir -p /root/configs
COPY configs/config.yml /root/configs/
COPY .env .

# Migrations

RUN mkdir -p /root/db/migrations
COPY db/migrations /root/db/migrations/

EXPOSE 8080

CMD ["./main"]