package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/caiofernandes00/playing-with-golang/grpc/cmd/util"
	"github.com/caiofernandes00/playing-with-golang/grpc/internal/entity"
	"github.com/caiofernandes00/playing-with-golang/grpc/internal/repository"
	"github.com/caiofernandes00/playing-with-golang/grpc/internal/service"
	"github.com/caiofernandes00/playing-with-golang/grpc/internal/service/interceptor"
	"github.com/caiofernandes00/playing-with-golang/grpc/pkg/proto/pb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

const (
	secretKey      = "secret"
	tokenDuration  = 15 * time.Minute
	serverCertFile = "cert/server-cert.pem"
	serverKeyFile  = "cert/server-key.pem"
)

func createUser(userStore repository.UserStore, username, pasword, role string) error {
	user, err := entity.NewUser(username, pasword, role)
	if err != nil {
		log.Fatal("cannot create user: ", err)
	}

	return userStore.Save(user)
}

func seedUsers(userStore repository.UserStore) error {
	err := createUser(userStore, "admin1", "secret", "admin")
	if err != nil {
		return err
	}

	return createUser(userStore, "user1", "secret", "user")
}

func accessibleRoles() map[string][]string {
	return map[string][]string{
		pb.LaptopService_CreateLaptop_FullMethodName: {"admin"},
		pb.LaptopService_UploadImage_FullMethodName:  {"admin"},
		pb.LaptopService_RateLaptop_FullMethodName:   {"admin", "user"},
	}
}

func loatTLSCredentials() (credentials.TransportCredentials, error) {
	certPool, err := util.LoadCAPool()
	if err != nil {
		return nil, err
	}

	serverCert, err := tls.LoadX509KeyPair(serverCertFile, serverKeyFile)
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	return credentials.NewTLS(config), nil
}

func runGRPCServer(
	laptopServer *service.LaptopServer,
	authServer *service.AuthServer,
	jwtManager *service.JWTManager,
	listener net.Listener,
	enableTLS bool,
) error {
	authInteceptor := interceptor.NewAuthInterceptor(jwtManager, accessibleRoles())
	serverOptions := []grpc.ServerOption{
		grpc.UnaryInterceptor(authInteceptor.Unary()),
		grpc.StreamInterceptor(authInteceptor.Stream()),
	}

	if enableTLS {
		tlsCredentials, err := loatTLSCredentials()
		if err != nil {
			return fmt.Errorf("cannot load TLS credentials: %w", err)
		}
		serverOptions = append(serverOptions, grpc.Creds(tlsCredentials))
	}

	grpcServer := grpc.NewServer(serverOptions...)
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)
	pb.RegisterAuthServiceServer(grpcServer, authServer)
	reflection.Register(grpcServer)

	err := grpcServer.Serve(listener)
	if err != nil {
		return fmt.Errorf("cannot start server: %w", err)
	}

	return nil
}

func runRESTServer(
	laptopServer *service.LaptopServer,
	authServer *service.AuthServer,
	jwtManager *service.JWTManager,
	listener net.Listener,
	enableTLS bool,
) error {
	mux := runtime.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := pb.RegisterAuthServiceHandlerServer(ctx, mux, authServer)
	if err != nil {
		return fmt.Errorf("cannot register auth server: %w", err)
	}

	err = pb.RegisterLaptopServiceHandlerServer(ctx, mux, laptopServer)
	if err != nil {
		return fmt.Errorf("cannot register laptop server: %w", err)
	}

	log.Printf("start REST server")
	if enableTLS {
		return http.ServeTLS(listener, mux, serverCertFile, serverKeyFile)
	}

	return http.Serve(listener, mux)
}

func main() {
	port := flag.Int("port", 0, "the server port")
	enableTLS := flag.Bool("tls", false, "enable SSL/TLS")
	serverType := flag.String("type", "grpc", "type of server (grpc/rest)")
	flag.Parse()
	log.Printf("start server on port: %d, TLS: %t", *port, *enableTLS)

	laptopStore := repository.NewInMemoryLaptopStore()
	imageStore := repository.NewDiskImageStore("tmp/")
	ratingStore := repository.NewInMemoryRatingStore()
	userStore := repository.NewInMemoryUserStore()
	jwtManager := service.NewJWTManager(secretKey, tokenDuration)

	laptopServer := service.NewLaptopServer(laptopStore, imageStore, ratingStore)
	authServer := service.NewAuthServer(userStore, jwtManager)

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

	err = seedUsers(userStore)
	if err != nil {
		log.Fatal("cannot seed users: %w", err)
	}

	if *serverType == "grpc" {
		err = runGRPCServer(laptopServer, authServer, jwtManager, listener, *enableTLS)
	} else {
		err = runRESTServer(laptopServer, authServer, jwtManager, listener, *enableTLS)
	}
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
