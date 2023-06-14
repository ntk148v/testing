FROM golang:1.20 AS build-stage
WORKDIR /app
RUN go mod init test
COPY *.go ./
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /rediscli

FROM alpine AS release-stage
WORKDIR /
COPY --from=build-stage /rediscli /rediscli
CMD ["/rediscli"]
