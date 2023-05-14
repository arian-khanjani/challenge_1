FROM golang:alpine

WORKDIR /user
COPY . /user
ENV PORT $USER_PORT

RUN go get -d -v
RUN go install -v

EXPOSE $PORT

ENTRYPOINT /go/bin/user