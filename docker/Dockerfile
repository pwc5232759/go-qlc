# Build gqlc in a stock Go builder container
FROM golang:1.12.1-alpine as builder

ARG BUILD_ACT=build

RUN apk add --no-cache make gcc musl-dev linux-headers git

ADD . /qlcchain/go-qlc
RUN cd /qlcchain/go-qlc && make clean ${BUILD_ACT}

# Pull gqlc into a second stage deploy alpine container
FROM alpine:latest

COPY --from=builder /qlcchain/go-qlc/build/gqlc .
COPY --from=builder /qlcchain/go-qlc/docker/entrypoint.sh .

RUN chmod +x entrypoint.sh

EXPOSE 9734 9735 9736

ENV seed=""

ENTRYPOINT [ "/bin/sh", "entrypoint.sh" ]
