FROM golang:latest AS buildstage

ENV HTTP_PROXY http://127.0.0.1:10809
ENV HTTPS_PROXY http://127.0.0.1:10809

WORKDIR /go/src

ADD . /go/src

RUN cd /go/src/basic_go && go build -o main ./webook/


FROM ubuntu:22.04

WORKDIR /app

COPY --from=buildstage /go/src/basic_go/main /app/

ENTRYPOINT ["./main"]