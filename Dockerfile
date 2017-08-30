FROM ptpadan1246/go-glide-docker
MAINTAINER Yasushi kobayashi <ptpadan@gmail.com>

ENV DOCUMENT_ENV production
EXPOSE 5000

COPY ./src/app /go/src/app
WORKDIR /go/src/app
RUN glide install && \
  go build main.go

