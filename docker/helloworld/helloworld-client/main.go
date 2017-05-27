package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"

	proto "github.com/linkerd/linkerd-examples/docker/helloworld/proto"
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
	streaming := flag.Bool("streaming", false, "send streaming requests")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <host>:<port> [flags]\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s <host>:<port> [flags]\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(1)
	}
	target := flag.Arg(0)

	conn, err := grpc.Dial(target, grpc.WithInsecure())
	dieIf(err)
	defer conn.Close()

	client := proto.NewHelloClient(conn)
	req := &proto.SvcRequest{}

	if *streaming {
		stream, err := client.StreamGreeting(context.Background(), req)
		dieIf(err)
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			dieIf(err)
			fmt.Println(resp.Message)
		}
	} else {
		resp, err := client.Greeting(context.Background(), req)
		dieIf(err)
		fmt.Println(resp.Message)
	}
}
