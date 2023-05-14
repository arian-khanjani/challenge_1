FROM golang:alpine

WORKDIR /salt
COPY . /salt
ENV PORT $SALT_PORT

RUN go get -d -v
RUN go install -v

EXPOSE $PORT

ENTRYPOINT /go/bin/salt