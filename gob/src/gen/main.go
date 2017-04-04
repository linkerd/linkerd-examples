package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	pb "github.com/linkerd/linkerd-examples/gob/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
	certFile := flag.String("cert", "", "Path to PEM-encoded certificate")
	keyFile := flag.String("key", "", "Path to PEM-encoded secret key")
	flag.Parse()
	if flag.NArg() != 0 {
		dieIf(fmt.Errorf("expecting zero arguments but got %d", flag.NArg()))
	}

	svc := &genSvc{}

	var server *grpc.Server
	if *keyFile == "" && *certFile == "" {
		server = grpc.NewServer()
	} else if *certFile == "" {
		dieIf(fmt.Errorf("key specified with no cert"))
	} else if *keyFile == "" {
		dieIf(fmt.Errorf("cert specified with no keey"))
	} else {
		pair, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		dieIf(err)
		creds := grpc.Creds(pair)
		server = grpc.NewServer(creds)
	}
	lis, err := net.Listen("tcp", *addr)
	dieIf(err)
	pb.RegisterGenSvcServer(server, svc)
	server.Serve(lis)
}
