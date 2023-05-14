FROM golang:alpine

WORKDIR /user
ADD . /user

RUN ls -lha
RUN go get -d -v
RUN go install -v

EXPOSE 3001

ENTRYPOINT /go/bin/user
