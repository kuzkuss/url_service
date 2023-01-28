FROM golang:latest AS builder

COPY . /server/

WORKDIR /server/

RUN go mod tidy
RUN go mod download
RUN go build -o main cmd/main.go

FROM ubuntu:20.04

COPY --from=builder /server/main .
COPY --from=builder /server/config/config.toml .

CMD ./main

