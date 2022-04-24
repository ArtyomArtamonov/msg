package auth

import (
	"context"

	pb "github.com/ArtyomArtamonov/msg/pkg/auth/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer

	userStore  UserStore
	jwtManager *JWTManager
}

func NewAuthService(userStore UserStore, jwtManager *JWTManager) *AuthServer {
	return &AuthServer{
		userStore:  userStore,
		jwtManager: jwtManager,
	}
}

func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.TokenResponse, error) {
	user, err := s.userStore.Find(req.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot find user %v", err)
	}

	if user == nil || !user.IsCorrectPassword(req.Password) {
		return nil, status.Error(codes.NotFound, "incorrect username or password")
	}

	token, err := s.jwtManager.Generate(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot generate access token")
	}

	res := pb.TokenResponse{
		AccessToken: token,
	}

	return &res, nil
}
