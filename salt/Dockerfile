FROM golang:alpine

WORKDIR /salt
COPY . /salt

RUN go get -d -v
RUN go install -v

EXPOSE $SALT_PORT

ENTRYPOINT /go/bin/salt