FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o comments-graphql cmd/server/server.go

# RUN go build -o main .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/comments-graphql .
COPY config.yaml .

EXPOSE 8080

CMD ["./comments-graphql"]
