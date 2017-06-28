FROM ptpadan1246/go-glide-docker
MAINTAINER Yasushi kobayashi <ptpadan@gmail.com>
COPY ./src/app /go/src/app
WORKDIR /go/src/app
RUN glide install && \
  go build main.go
EXPOSE 5000

