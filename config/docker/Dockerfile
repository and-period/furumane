# ビルド用コンテナ
FROM golang:1.21.1-alpine as builder

ENV LANG C.UTF-8
ENV TZ Asia/Tokyo

ARG SERVICE

WORKDIR /go/src/github.com/and-period/furumane

RUN apk add --update --no-cache \
      git \
      make \
      tzdata

ADD . .

RUN make build SERVICE=${SERVICE}

# 実行用コンテナ
FROM alpine

ENV LANG C.UTF-8
ENV TZ Asia/Tokyo

WORKDIR /var/api

RUN apk add --update --no-cache \
      ca-certificates \
      curl

COPY --from=builder /go/src/github.com/and-period/furumane/app ./app
COPY --from=builder /usr/share/zoneinfo/Asia/Tokyo /usr/share/zoneinfo/Asia/Tokyo

EXPOSE 8080 9090
CMD ["./app"]
