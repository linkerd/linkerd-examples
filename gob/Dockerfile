FROM library/golang:1.6.0

ADD . /go/src/github.com/linkerd/linkerd-examples/gob

RUN ["go", "build", "-o", "/usr/bin/web", "/go/src/github.com/linkerd/linkerd-examples/gob/src/web/main.go"]
RUN ["go", "build", "-o", "/usr/bin/gen", "/go/src/github.com/linkerd/linkerd-examples/gob/src/gen/main.go"]
RUN ["go", "build", "-o", "/usr/bin/word", "/go/src/github.com/linkerd/linkerd-examples/gob/src/word/main.go"]
