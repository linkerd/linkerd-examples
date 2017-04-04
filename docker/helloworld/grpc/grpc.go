package grpc

import (
	"fmt"
	"math/rand"
	"time"

	proto "github.com/linkerd/linkerd-examples/docker/helloworld/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Server struct {
	text        string
	target      proto.SvcClient
	podIp       string
	latency     time.Duration
	failureRate float64
}

func New(text, target, podIp string, latency time.Duration, failureRate float64) (*Server, error) {
	var client proto.SvcClient
	if target != "" {
		conn, err := grpc.Dial(target, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		client = proto.NewSvcClient(conn)
	}

	return &Server{
		text:        text,
		target:      client,
		podIp:       podIp,
		latency:     latency,
		failureRate: failureRate,
	}, nil
}

func (s *Server) Hello(ctx context.Context, req *proto.SvcRequest) (*proto.SvcResponse, error) {
	return s.respond(ctx, req)
}

func (s *Server) World(ctx context.Context, req *proto.SvcRequest) (*proto.SvcResponse, error) {
	return s.respond(ctx, req)
}

func (s *Server) respond(_ context.Context, _ *proto.SvcRequest) (*proto.SvcResponse, error) {
	time.Sleep(s.latency)
	if rand.Float64() < s.failureRate {
		return nil, fmt.Errorf("server error")
	}

	text := s.text
	if s.podIp != "" {
		text += fmt.Sprintf(" (%s)", s.podIp)
	}

	if s.target != nil {
		targetText, err := s.callTarget()
		if err != nil {
			return nil, err
		}
		text += fmt.Sprintf(" %s", targetText)
	}

	return &proto.SvcResponse{Message: text + "!"}, nil
}

func (s *Server) callTarget() (string, error) {
	resp, err := s.target.World(context.Background(), &proto.SvcRequest{})
	if err != nil {
		return "", err
	}

	return resp.GetMessage(), nil
}
