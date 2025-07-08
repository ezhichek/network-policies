FROM golang:1.23-alpine

WORKDIR /app

COPY . .

RUN apk add --no-cache curl && GO111MODULE=off go build -o backend-tester .


CMD ["./backend-tester"]
