FROM library/golang:1.6.0

ADD . /go/src/github.com/buoyantio/linkerd-examples/gob

RUN ["go", "build", "-o", "/usr/bin/web", "/go/src/github.com/buoyantio/linkerd-examples/gob/web/main.go"]
RUN ["go", "build", "-o", "/usr/bin/gen", "/go/src/github.com/buoyantio/linkerd-examples/gob/gen/main.go"]
RUN ["go", "build", "-o", "/usr/bin/word", "/go/src/github.com/buoyantio/linkerd-examples/gob/word/main.go"]
