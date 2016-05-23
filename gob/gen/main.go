package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	pb "github.com/buoyantio/linkerd-examples/gob/proto"
	"google.golang.org/grpc"
)

type genSvc struct{}

func (s *genSvc) Gen(req *pb.GenRequest, stream pb.GenSvc_GenServer) error {
	if err := stream.Send(&pb.GenResponse{req.Text}); err != nil {
		return err
	}
	doWrite := func() bool {
		err := stream.Send(&pb.GenResponse{" " + req.Text})
		return err == nil
	}
	if req.Limit == 0 {
		for {
			if ok := doWrite(); !ok {
				break
			}
		}
	} else {
		for i := uint(1); i != uint(req.Limit); i++ {
			if ok := doWrite(); !ok {
				break
			}
		}
	}
	return nil
}

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
