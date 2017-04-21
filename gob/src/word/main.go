package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"

	pb "github.com/linkerd/linkerd-examples/gob/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type wordSvc struct {
	words []string
}

func (svc *wordSvc) GetWord(ctx context.Context, in *pb.WordRequest) (*pb.WordResponse, error) {
	word := svc.randomWord()
	if word == "" {
		return nil, fmt.Errorf("empty word")
	}
	return &pb.WordResponse{word}, nil
}

func (svc *wordSvc) randomWord() string {
	n := len(svc.words)
	switch n {
	case 0:
		return ""
	default:
		idx := rand.Int() % n
		return svc.words[idx]
	}
}

func dieIf(err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "Error: %s. Try --help for help.\n", err)
	os.Exit(-1)
}

func main() {
	addr := flag.String("srv", ":8282", "TCP address to listen on (in host:port form)")
	certFile := flag.String("cert", "", "Path to PEM-encoded certificate")
	keyFile := flag.String("key", "", "Path to PEM-encoded secret key")
	flag.Parse()
	if flag.NArg() != 0 {
		dieIf(fmt.Errorf("expecting zero arguments but got %d", flag.NArg()))
	}

	svc := &wordSvc{
		words: []string{
			"banana",
			"bees",
			"cmon",
			"gob",
			"illusion",
			"same",
		},
	}

	var server *grpc.Server
	if *keyFile == "" && *certFile == "" {
		server = grpc.NewServer()
	} else if *certFile == "" {
		dieIf(fmt.Errorf("key specified with no cert"))
	} else if *keyFile == "" {
		dieIf(fmt.Errorf("cert specified with no key"))
	} else {
		pair, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		dieIf(err)
		creds := grpc.Creds(pair)
		server = grpc.NewServer(creds)
	}
	lis, err := net.Listen("tcp", *addr)
	dieIf(err)
	pb.RegisterWordSvcServer(server, svc)
	server.Serve(lis)
}
