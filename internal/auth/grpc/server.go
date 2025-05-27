package grpc

import (
	"context"
	"log"
	"net"

	"github.com/mos1rain/forum_go/internal/auth/repository"
	"github.com/mos1rain/forum_go/pkg/jwt"
	"github.com/mos1rain/forum_go/proto/auth"
	"google.golang.org/grpc"
)

type AuthGRPCServer struct {
	auth.UnimplementedAuthServiceServer
	repo      *repository.UserRepository
	tokenMngr *jwt.TokenManager
}

func NewAuthGRPCServer(repo *repository.UserRepository, tokenMngr *jwt.TokenManager) *AuthGRPCServer {
	return &AuthGRPCServer{repo: repo, tokenMngr: tokenMngr}
}

func (s *AuthGRPCServer) ValidateToken(ctx context.Context, req *auth.ValidateTokenRequest) (*auth.ValidateTokenResponse, error) {
	claims, err := s.tokenMngr.Parse(req.Token)
	if err != nil {
		return &auth.ValidateTokenResponse{Valid: false, Error: err.Error()}, nil
	}
	return &auth.ValidateTokenResponse{
		UserId:   int32(claims.UserID),
		Username: claims.Username,
		Valid:    true,
	}, nil
}

func (s *AuthGRPCServer) GetUserByID(ctx context.Context, req *auth.GetUserByIDRequest) (*auth.GetUserByIDResponse, error) {
	user, err := s.repo.GetByID(int(req.UserId))
	if err != nil {
		return &auth.GetUserByIDResponse{Error: err.Error()}, nil
	}
	if user == nil {
		return &auth.GetUserByIDResponse{Error: "user not found"}, nil
	}
	return &auth.GetUserByIDResponse{
		UserId:   int32(user.ID),
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func RunGRPCServer(repo *repository.UserRepository, tokenMngr *jwt.TokenManager, addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	auth.RegisterAuthServiceServer(grpcServer, NewAuthGRPCServer(repo, tokenMngr))
	log.Printf("gRPC Auth server started on %s", addr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
