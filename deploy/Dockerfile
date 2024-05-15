FROM golang:1.19-alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOOS linux

RUN apk update --no-cache && apk add --no-cache tzdata
RUN apk add make

COPY . /yadroclub/

WORKDIR /yadroclub/

RUN make build
RUN chown -R root ./bin/events

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /yadroclub/bin .
COPY --from=builder /yadroclub/testdata ./testdata
ARG DATA_PATH="./testdata/test.txt"

CMD [ "sh", "-c", "./events -path ${DATA_PATH}" ]