FROM library/golang:1.6.0
RUN ["mkdir", "-p", "/src"]
COPY ./main.go /src/
RUN ["go", "build", "-o", "/usr/bin/gobwebsvc", "/src/main.go"]
ENTRYPOINT ["/usr/bin/gobwebsvc"]
