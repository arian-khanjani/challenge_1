FROM golang:alpine

WORKDIR /salt
COPY . /salt

RUN ls -lha
RUN go get -d -v
RUN go install -v

EXPOSE 3000

ENTRYPOINT /go/bin/salt