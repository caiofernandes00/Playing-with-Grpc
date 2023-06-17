package auth

import (
	"context"
	"time"

	"github.com/caiofernandes00/playing-with-golang/grpc/pkg/proto/pb"
	"google.golang.org/grpc"
)

type AuthClient struct {
	service  pb.AuthServiceClient
	username string
	password string
}

func NewAuthClient(conn *grpc.ClientConn, username, password string) *AuthClient {
	service := pb.NewAuthServiceClient(conn)
	return &AuthClient{service, username, password}
}

func (client *AuthClient) Login() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.LoginRequest{
		Username: client.username,
		Password: client.password,
	}

	resp, err := client.service.Login(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.GetAccessToken(), nil
}
