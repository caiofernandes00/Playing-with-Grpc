package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/caiofernandes00/playing-with-golang/grpc/internal/entity"
	"github.com/caiofernandes00/playing-with-golang/grpc/internal/repository"
	"github.com/caiofernandes00/playing-with-golang/grpc/internal/service"
	"github.com/caiofernandes00/playing-with-golang/grpc/internal/service/interceptor"
	"github.com/caiofernandes00/playing-with-golang/grpc/pkg/proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	secretKey     = "secret"
	tokenDuration = 15 * time.Minute
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
	const laptopServicePath = "/playingwithgolang.grpc.LaptopService/"
	return map[string][]string{
		laptopServicePath + "CreateLaptop": {"admin"},
		laptopServicePath + "UploadImage":  {"admin"},
		laptopServicePath + "RateLaptop":   {"admin", "user"},
	}
}

func main() {
	port := flag.Int("port", 0, "the server port")
	flag.Parse()
	log.Printf("start server on port: %d", *port)

	laptopStore := repository.NewInMemoryLaptopStore()
	imageStore := repository.NewDiskImageStore("tmp/")
	ratingStore := repository.NewInMemoryRatingStore()
	userStore := repository.NewInMemoryUserStore()
	jwtManager := service.NewJWTManager(secretKey, tokenDuration)

	laptopServer := service.NewLaptopServer(laptopStore, imageStore, ratingStore)
	authServer := service.NewAuthServer(userStore, jwtManager)

	authInteceptor := interceptor.NewAuthInterceptor(jwtManager, accessibleRoles())

	err := seedUsers(userStore)
	if err != nil {
		log.Fatal("cannot seed users: ", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInteceptor.Unary()),
		grpc.StreamInterceptor(authInteceptor.Stream()),
	)
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)
	pb.RegisterAuthServiceServer(grpcServer, authServer)
	reflection.Register(grpcServer)

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
