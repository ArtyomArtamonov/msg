package auth

import (
	"context"
	"time"
	"unicode/utf8"

	pb "github.com/ArtyomArtamonov/msg/internal/auth/proto"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer

	userStore         UserStore
	refreshTokenStore RefreshTokenStore
	jwtManager        *JWTManager
}

func NewAuthService(userStore UserStore, refreshTokenStore RefreshTokenStore, jwtManager *JWTManager) *AuthServer {
	return &AuthServer{
		userStore:         userStore,
		jwtManager:        jwtManager,
		refreshTokenStore: refreshTokenStore,
	}
}

func (s *AuthServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.TokenResponse, error) {
	if user, _ := s.userStore.FindByUsername(req.Username); user != nil {
		return nil, status.Errorf(codes.AlreadyExists, "user already exists")
	}

	if utf8.RuneCountInString(req.Username) > 15 {
		return nil, status.Error(codes.InvalidArgument, "username could not be more than 15 characters")
	}

	// TODO: add more password validation
	if utf8.RuneCountInString(req.Password) < 6 {
		return nil, status.Error(codes.InvalidArgument, "password could not be less than 6 characters")
	}

	user, err := NewUser(req.Username, req.Password, USER_ROLE)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	tokenPair, err := s.jwtManager.Generate(user)
	if err != nil {
		return nil, err
	}

	if err := s.userStore.Save(user); err != nil {
		return nil, status.Errorf(codes.Internal, "could not save user: %v", err)
	}

	if err := s.refreshTokenStore.Add(tokenPair.RefreshToken); err != nil {
		return nil, status.Errorf(codes.Internal, "could not add refresh key to database: %v", err)
	}

	response := pb.TokenResponse{
		AccessToken:  tokenPair.JwtToken,
		RefreshToken: tokenPair.RefreshToken.Token.String(),
	}

	return &response, nil
}

func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.TokenResponse, error) {
	user, err := s.userStore.FindByUsername(req.Username)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "incorrect username or password")
	}

	if user == nil || !user.IsCorrectPassword(req.Password) {
		return nil, status.Error(codes.NotFound, "incorrect username or password")
	}

	tokenPair, err := s.jwtManager.Generate(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not generate access token")
	}

	if err := s.refreshTokenStore.Add(tokenPair.RefreshToken); err != nil {
		return nil, status.Errorf(codes.Internal, "could not save refresh token to database: %v", err)
	}

	res := pb.TokenResponse{
		AccessToken:  tokenPair.JwtToken,
		RefreshToken: tokenPair.RefreshToken.Token.String(),
	}

	return &res, nil
}

func (s *AuthServer) Refresh(ctx context.Context, req *pb.RefreshRequest) (*pb.TokenResponse, error) {
	refreshUUID, err := uuid.Parse(req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "could not parse refresh token")
	}

	token, err := s.refreshTokenStore.Get(refreshUUID)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "refresh token does not exists")
	}

	if token.ExpiresAt.Unix() < time.Now().Unix() {
		if err := s.refreshTokenStore.Delete(refreshUUID); err != nil {
			logrus.Error("could not delete old refresh token: %v", err)
		}
		return nil, status.Error(codes.Unauthenticated, "refresh token is expired")
	}

	user, err := s.userStore.Find(token.UserId)
	if err != nil {
		logrus.Error("needs to be investigated: %v", err)
		return nil, status.Error(codes.Internal, "hmm... this is strange. That could not possibly happen")
	}

	tokenPair, err := s.jwtManager.Generate(user)
	if err != nil {
		return nil, status.Error(codes.Internal, "could not generate token pair")
	}

	if err := s.refreshTokenStore.Delete(refreshUUID); err != nil {
		return nil, status.Errorf(codes.Internal, "could not delete old refresh token: %v", err)
	}

	if err := s.refreshTokenStore.Add(tokenPair.RefreshToken); err != nil {
		return nil, status.Errorf(codes.Internal, "could not save new refresh token: %v", err)
	}

	return &pb.TokenResponse{
		AccessToken:  tokenPair.JwtToken,
		RefreshToken: tokenPair.RefreshToken.Token.String(),
	}, nil
}