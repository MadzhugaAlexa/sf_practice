FROM golang:latest as source
LABEL version="1.0.0"
LABEL maintainer="Alexandra Madzhuga <madzhav@me.com>"
RUN mkdir -p /go/src
WORKDIR /go/src
ADD main.go .
ADD go.mod .
RUN go install .

FROM alpine:latest
LABEL version="1.0.0"
LABEL maintainer="Alexandra Madzhuga <madzhav@me.com>"
WORKDIR /root/
COPY --from=source /go/bin/hw .
ENTRYPOINT ./hw

