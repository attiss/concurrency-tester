# syntax=docker/dockerfile:1
FROM golang:1.16
COPY . /concurrency-tester
WORKDIR /concurrency-tester
RUN go build .
CMD ["./concurrency-tester"]
