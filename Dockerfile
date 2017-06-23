FROM golang:1.8
MAINTAINER Yasushi kobayashi <ptpadan@gmail.com>

RUN curl https://glide.sh/get | sh

COPY ./src/app /go/src/app

WORKDIR /go/src/app
RUN glide install && \
  go build main.go
EXPOSE 5000
VOLUME /go/src/app/log
VOLUME /go/src/app/static

# CMD /go/src/app/main &
