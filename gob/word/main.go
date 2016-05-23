package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"

	pb "github.com/buoyantio/linkerd-examples/gob/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
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
	//
	// Parse flags
	//
	addr := flag.String("srv", ":8282", "TCP address to listen on (in host:port form)")
	flag.Parse()
	if flag.NArg() != 0 {
		dieIf(fmt.Errorf("expecting zero arguments but got %d", flag.NArg()))
	}

	lis, err := net.Listen("tcp", *addr)
	dieIf(err)

	//
	// Setup Http server
	//
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

	s := grpc.NewServer()
	pb.RegisterWordSvcServer(s, svc)
	s.Serve(lis)
}
