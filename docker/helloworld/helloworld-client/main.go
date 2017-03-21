package main

import (
	"flag"
	"fmt"
	"os"

	proto "github.com/buoyantio/linkerd-examples/docker/helloworld/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func dieIf(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func main() {
	target := flag.String("target", "localhost:7777", "address of the hello gRPC service")
	flag.Parse()

	conn, err := grpc.Dial(*target, grpc.WithInsecure())
	dieIf(err)

	defer conn.Close()
	client := proto.NewSvcClient(conn)

	resp, err := client.Hello(context.Background(), &proto.SvcRequest{})
	dieIf(err)

	fmt.Println(resp.Message)
}
