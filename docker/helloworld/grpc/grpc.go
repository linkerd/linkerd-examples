package grpc

import (
	"fmt"
	"math/rand"
	"time"

	proto "github.com/buoyantio/linkerd-examples/docker/helloworld/proto"
	"github.com/buoyantio/linkerd-examples/docker/helloworld/redis"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Server struct {
	text        string
	target      proto.SvcClient
	podIp       string
	latency     time.Duration
	failureRate float64
	redis       *redis.Client
}

func New(text, target, podIp string, latency time.Duration, failureRate float64, redisClient *redis.Client) (*Server, error) {
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
		redis:       redisClient,
	}, nil
}

func (s *Server) Hello(ctx context.Context, req *proto.SvcRequest) (*proto.SvcResponse, error) {
	return s.respond(ctx, req)
}

func (s *Server) World(ctx context.Context, req *proto.SvcRequest) (*proto.SvcResponse, error) {
	return s.respond(ctx, req)
}

func (s *Server) respond(_ context.Context, _ *proto.SvcRequest) (*proto.SvcResponse, error) {
	if s.redis != nil {
		if text, err := s.redis.Get(); err == nil {
			return &proto.SvcResponse{Message: text}, nil
		}
	}

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

	text += "!"

	if s.redis != nil {
		s.redis.Set(text)
	}

	return &proto.SvcResponse{Message: text}, nil
}

func (s *Server) callTarget() (string, error) {
	resp, err := s.target.World(context.Background(), &proto.SvcRequest{})
	if err != nil {
		return "", err
	}

	return resp.GetMessage(), nil
}
