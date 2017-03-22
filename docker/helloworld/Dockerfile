FROM alpine:3.5

RUN apk add --update curl jq && rm -rf /var/cache/apk/*

ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH

RUN mkdir -p $GOPATH/bin
COPY ./helloworld $GOPATH/bin
COPY ./helloworld-client/helloworld-client $GOPATH/bin
COPY ./hostIP.sh /usr/local/bin

ENTRYPOINT ["helloworld"]
