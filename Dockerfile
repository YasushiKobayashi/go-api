FROM ptpadan1246/go-glide-docker
MAINTAINER Yasushi kobayashi <ptpadan@gmail.com>

RUN curl https://glide.sh/get | sh

COPY ./src/app /go/src/app

WORKDIR /go/src/app
RUN glide install && \
  go build main.go
EXPOSE 5000

