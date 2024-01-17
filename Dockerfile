FROM golang:1-alpine AS build

RUN apk add --no-cache git

WORKDIR /go/src/github.com/7Maliko7/telegram-bot

COPY ./ /go/src/github.com/7Maliko7/telegram-bot
COPY ./config.yml /bin/config.yml
COPY ./scenario.yml /bin/scenario.yml

RUN go build -v -o /bin/app /go/src/github.com/7Maliko7/telegram-bot

FROM alpine:3.17.5

RUN apk add ca-certificates

COPY --from=build /bin/app /bin/app
COPY --from=build /bin/config.yml /bin/config.yml
COPY --from=build /bin/scenario.yml /bin/scenario.yml

RUN cat /bin/config.yml && cat /bin/scenario.yml

EXPOSE 2112

ENTRYPOINT ["/bin/app", "-c", "/bin/config.yml", "-s", "/bin/scenario.yml"]