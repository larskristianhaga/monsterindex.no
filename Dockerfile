FROM golang:1.20-buster AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN apt-get update && apt-get install -y libsqlite3-dev

RUN go mod download && go mod verify

ENV CGO_ENABLED=1

COPY . .

RUN go build -v -o /run-app .

FROM debian:buster-slim

COPY --from=builder /run-app /usr/local/bin/

EXPOSE 8080
CMD ["run-app"]
