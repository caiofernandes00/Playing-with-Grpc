package service

import (
	"context"

	"github.com/caiofernandes00/playing-with-golang/grpc/internal/repository"
	"github.com/caiofernandes00/playing-with-golang/grpc/pkg/proto/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	userStore  repository.UserStore
	jwtManager *JWTManager
}

func NewAuthServer(userStore repository.UserStore, jwtManager *JWTManager) *AuthServer {
	return &AuthServer{
		userStore:  userStore,
		jwtManager: jwtManager,
	}
}

func (server *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := server.userStore.Find(req.Username)
	if err != nil {
		return nil, err
	}

	if user == nil || !user.IsCorrectPassword(req.GetPassword()) {
		return nil, status.Errorf(codes.NotFound, "incorrect username/password")
	}

	token, err := server.jwtManager.Generate(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot generate access token")
	}

	res := &pb.LoginResponse{
		AccessToken: token,
	}
	return res, nil
}
