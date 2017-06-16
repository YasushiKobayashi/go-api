FROM golang:1.8
MAINTAINER Yasushi kobayashi <ptpadan@gmail.com>

RUN curl https://glide.sh/get | sh

COPY ./src/document /go/src/document

WORKDIR /go/src/document
RUN glide install && \
  go build main.go

CMD ./main &

EXPOSE 5000
VOLUME /go/src/document/log
