# build
FROM golang:alpine as builder
WORKDIR /app
#RUN go env -w GOPROXY=https://goproxy.io,direct
ADD go.mod .
ADD go.sum .
RUN go mod download -x
ADD . .
RUN go build -v -o sync-up cmd/sync-up/main.go

FROM alpine
MAINTAINER "catfishlty"
WORKDIR /app
COPY --from=builder /app/sync-up /app/sync-up
ENTRYPOINT /app/sync-up
