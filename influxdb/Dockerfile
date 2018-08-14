FROM golang:1.10.3 as builder
WORKDIR /go/src/github.com/linkerd/linkerd-examples/influxdb/
RUN go get -d -v github.com/prometheus/client_golang/prometheus
COPY app.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM alpine:3.8
RUN apk --no-cache add curl
WORKDIR /root/
COPY --from=builder /go/src/github.com/linkerd/linkerd-examples/influxdb/app .
ENTRYPOINT ["./app"]
