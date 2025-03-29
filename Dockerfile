FROM golang:1.23-alpine3.21

LABEL maintainer="maksimacx50@gmail.com"

WORKDIR /tgbot

COPY . .

RUN go mod tidy

WORKDIR /tgbot/cmd/app

VOLUME /tgbot/logs

ENTRYPOINT ["go", "run", "main.go"]
