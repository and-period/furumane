# ビルド用コンテナ
FROM golang:1.21.1-alpine as builder

ENV LANG C.UTF-8
ENV TZ Asia/Tokyo

ARG SERVICE

WORKDIR /go/src/github.com/and-period/furumane

RUN apk add --update --no-cache \
      git \
      make

ADD . .

RUN make install
RUN make proto
RUN make build SERVICE=${SERVICE}

# 実行用コンテナ
FROM alpine

ENV LANG C.UTF-8
ENV TZ Asia/Tokyo

WORKDIR /var/api

RUN apk add --update --no-cache \
      ca-certificates

COPY --from=builder /go/src/github.com/and-period/furumane/app ./app

EXPOSE 9000 9001
CMD ["./app"]