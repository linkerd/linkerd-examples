FROM golang:1.10.3 AS builder
LABEL maintainer="linkerd-users@googlegroups.com"

ADD https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

WORKDIR $GOPATH/src/github.com/linkerd/linkerd-examples/perf-baseline/reporter
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /reporter .

FROM scratch
COPY --from=builder /reporter ./
ENTRYPOINT ["./reporter"]
