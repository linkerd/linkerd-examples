package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	grpcServer "github.com/linkerd/linkerd-examples/docker/helloworld/grpc"
	httpServer "github.com/linkerd/linkerd-examples/docker/helloworld/http"
	proto "github.com/linkerd/linkerd-examples/docker/helloworld/proto"
	"google.golang.org/grpc"
)

const (
	httpTimeout = 10 * time.Second
)

func dieIf(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

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
		if *latency > httpTimeout {
			dieIf(fmt.Errorf("latency cannot exceed %s", httpTimeout))
		}

		server := &http.Server{
			Addr:         *addr,
			ReadTimeout:  httpTimeout + time.Second,
			WriteTimeout: httpTimeout + time.Second,
			Handler:      httpServer.New(serverText, *target, podIp, *latency, *failureRate, *json),
		}

		go func() {
			fmt.Println("starting HTTP server on", *addr)
			server.ListenAndServe()
		}()

		<-stop

		fmt.Println("shutting down HTTP server on", *addr)
		server.Shutdown(context.Background())

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

		go func() {
			fmt.Println("starting gRPC server on", *addr)
			s.Serve(lis)
		}()

		<-stop

		fmt.Println("shutting down gRPC server on", *addr)
		s.GracefulStop()

	default:
		dieIf(fmt.Errorf("unsupported protocol: %s", *protocol))
	}
}
