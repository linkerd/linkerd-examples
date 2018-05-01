FROM golang:1.10.1-alpine3.7
WORKDIR /go/src/github.com/linkerd/linkerd-examples/add-steps/
RUN apk update && apk add git
RUN go get -d -v github.com/prometheus/client_golang/prometheus
COPY server.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM scratch
COPY --from=0 /go/src/github.com/linkerd/linkerd-examples/add-steps/app /app
ENTRYPOINT ["/app"]
