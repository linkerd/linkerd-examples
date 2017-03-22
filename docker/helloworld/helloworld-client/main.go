package main

import (
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
	if len(os.Args) != 2 {
		dieIf(fmt.Errorf("Usage: helloworld-client <host>:<port>"))
	}
	target := os.Args[1]

	conn, err := grpc.Dial(target, grpc.WithInsecure())
	dieIf(err)

	defer conn.Close()
	client := proto.NewSvcClient(conn)

	resp, err := client.Hello(context.Background(), &proto.SvcRequest{})
	dieIf(err)

	fmt.Println(resp.Message)
}
