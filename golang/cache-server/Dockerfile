#FROM golang:1.16-alpine as builder
#LABEL maintainer="Kien Nguyen-Tuan <kiennt2609@gmail.com>"
#ENV GO111MODULE=on
#ADD main.go main.go
#RUN go build -ldflags "-s -w" -o /bin/simple-cache-server main.go && \
#    chmod a+x /bin/simple-cache-server
#
FROM alpine:3.12
LABEL maintainer="Kien Nguyen <kiennt2609@gmail.com>"
RUN apk add --no-cache bash
COPY simple-cache-server /bin/simple-cache-server
EXPOSE 9010
ENTRYPOINT ["/bin/simple-cache-server"]
