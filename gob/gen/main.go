package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	pb "github.com/buoyantio/linkerd-examples/gob/proto/gen"
	"google.golang.org/grpc"
)

type genSvc struct{}

func (s *genSvc) Gen(req *pb.GenRequest, stream pb.GenSvc_GenServer) error {
	return nil
}

// Writes to the stream until `limit` writes have been completed or
// the stream is closed (i.e. because the client disconnects).

// func (svc *GenSvc) generate(text string, limit uint, stream io.Writer) error {
// 	if _, err := stream.Write([]byte(text)); err != nil {
// 		return err
// 	}
// 	doWrite := func() bool {
// 		_, err := stream.Write([]byte(" " + text))
// 		return err == nil
// 	}
// 	if limit == 0 {
// 		for {
// 			if ok := doWrite(); !ok {
// 				break
// 			}
// 		}
// 	} else {
// 		// start at 1 because we've already written 1
// 		for i := uint(1); i != limit; i++ {
// 			if ok := doWrite(); !ok {
// 				break
// 			}
// 		}
// 	}
// 	return nil
// }

func dieIf(err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "Error: %s. Try --help for help.\n", err)
	os.Exit(-1)
}

func main() {
	addr := flag.String("srv", ":8181", "TCP address to listen on (in host:port form)")
	flag.Parse()
	if flag.NArg() != 0 {
		dieIf(fmt.Errorf("expecting zero arguments but got %d", flag.NArg()))
	}

	lis, err := net.Listen("tcp", *addr)
	dieIf(err)

	s := grpc.NewServer()
	pb.RegisterGenSvcServer(s, &genSvc{})
	s.Serve(lis)
}
