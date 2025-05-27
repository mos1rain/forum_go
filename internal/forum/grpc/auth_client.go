package grpc

import (
	"context"
	"time"

	"github.com/mos1rain/forum_go/proto/auth"
	"google.golang.org/grpc"
)

type AuthGRPCClient struct {
	client auth.AuthServiceClient
}

func NewAuthGRPCClient(addr string) (*AuthGRPCClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(3*time.Second))
	if err != nil {
		return nil, err
	}
	return &AuthGRPCClient{client: auth.NewAuthServiceClient(conn)}, nil
}

func (c *AuthGRPCClient) ValidateToken(token string) (*auth.ValidateTokenResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return c.client.ValidateToken(ctx, &auth.ValidateTokenRequest{Token: token})
}
