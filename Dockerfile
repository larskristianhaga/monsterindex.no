FROM golang:1.23-bookworm AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

ENV CGO_ENABLED=1

COPY . .

RUN go build -v -o /run-app .

FROM debian:bookworm-slim

COPY --from=builder /run-app /usr/local/bin/

CMD ["run-app"]
