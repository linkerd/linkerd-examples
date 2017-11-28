package grpc

import (
	"context"
	"strings"

	"google.golang.org/grpc/metadata"
)

func linkerdContext(ctx context.Context) context.Context {
	pairs := make([]string, 0)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		for key, values := range md {
			if strings.HasPrefix(strings.ToLower(key), "l5d-ctx") {
				for _, value := range values {
					pairs = append(pairs, key, value)
				}
			}
		}
	}
	return metadata.NewOutgoingContext(context.Background(), metadata.Pairs(pairs...))
}
