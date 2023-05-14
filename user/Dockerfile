FROM golang:alpine

WORKDIR /user
COPY . /user

RUN go get -d -v
RUN go install -v

EXPOSE $USER_PORT

ENTRYPOINT /go/bin/user