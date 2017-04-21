package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	grpcServer "github.com/linkerd/linkerd-examples/docker/helloworld/grpc"
	httpServer "github.com/linkerd/linkerd-examples/docker/helloworld/http"
	proto "github.com/linkerd/linkerd-examples/docker/helloworld/proto"
	"google.golang.org/grpc"
)

func dieIf(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	addr := flag.String("addr", ":7777", "address to serve on")
	text := flag.String("text", "Hello", "text to serve")
	target := flag.String("target", "", "target service to call before returning")
	latency := flag.Duration("latency", 0, "time to sleep before processing request")
	failureRate := flag.Float64("failure-rate", 0.0, "rate of error responses to return")
	json := flag.Bool("json", false, "return JSON instead of plaintext responses (HTTP only)")
	protocol := flag.String("protocol", "http", "API protocol: http or grpc")
	flag.Parse()

	serverText := *text
	if envText := os.Getenv("TARGET_WORLD"); envText != "" {
		serverText = envText
	}

	podIp := os.Getenv("POD_IP")

	switch *protocol {
	case "http":
		server := &http.Server{
			Addr:         *addr,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			Handler:      httpServer.New(serverText, *target, podIp, *latency, *failureRate, *json),
		}

		fmt.Println("starting HTTP server on", *addr)

		err := server.ListenAndServe()
		dieIf(err)

	case "grpc":
		lis, err := net.Listen("tcp", *addr)
		dieIf(err)

		s := grpc.NewServer()
		server, err := grpcServer.New(serverText, *target, podIp, *latency, *failureRate)
		dieIf(err)

		if strings.ToLower(serverText) == "hello" {
			proto.RegisterHelloServer(s, server)
		} else {
			proto.RegisterWorldServer(s, server)
		}

		fmt.Println("starting gRPC server on", *addr)

		err = s.Serve(lis)
		dieIf(err)
	default:
		dieIf(fmt.Errorf("unsupported protocol: %s", *protocol))
	}
}
