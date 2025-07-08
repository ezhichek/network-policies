FROM golang:1.23-alpine

RUN apk add --no-cache curl

WORKDIR /app
COPY . .
RUN go build -o frontend .
CMD ["./frontend"]
