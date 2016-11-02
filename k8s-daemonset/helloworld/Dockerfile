FROM python:2-onbuild

RUN apt-get update && apt-get install -y --no-install-recommends \
        jq \
    && rm -rf /var/lib/apt/lists/*

RUN curl -O https://storage.googleapis.com/golang/go1.6.linux-amd64.tar.gz
RUN tar -xvf go1.6.linux-amd64.tar.gz
RUN mv go /usr/local

ENV GOPATH /go
RUN mkdir $GOPATH
ENV PATH /usr/local/go/bin:/go/bin:$PATH
RUN go get github.com/BuoyantIO/namerctl && go install github.com/BuoyantIO/namerctl
