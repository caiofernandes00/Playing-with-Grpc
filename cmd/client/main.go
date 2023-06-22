package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/caiofernandes00/playing-with-golang/grpc/cmd/client/auth"
	"github.com/caiofernandes00/playing-with-golang/grpc/cmd/util"
	"github.com/caiofernandes00/playing-with-golang/grpc/internal/sample"
	"github.com/caiofernandes00/playing-with-golang/grpc/pkg/proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func testcreateLaptop(laptopClient *auth.LaptopClient) {
	laptopClient.CreateLaptop(sample.NewLaptop())
}

func testSearchLaptop(laptopClient *auth.LaptopClient) {
	for i := 0; i < 10; i++ {
		laptopClient.CreateLaptop(sample.NewLaptop())
	}

	filter := &pb.Filter{
		MaxPriceUsd: 3000,
		MinCpuCores: 4,
		MinCpuGhz:   2.5,
		MinRam:      &pb.Memory{Value: 8, Unit: pb.Memory_GIGABYTE},
	}

	laptopClient.SearchLaptop(filter)
}

func testUploadImage(laptopClient *auth.LaptopClient) {
	laptop := sample.NewLaptop()
	laptopClient.CreateLaptop(laptop)
	laptopClient.UploadImage(laptop, "tmp/angry-cat.jpg")
}

func testRateLaptop(laptopClient *auth.LaptopClient) {
	n := 3
	laptops := make([]string, n)

	for i := 0; i < n; i++ {
		laptop := sample.NewLaptop()
		laptops[i] = laptop.GetId()
		laptopClient.CreateLaptop(laptop)
	}

	scores := make([]float64, n)
	for {
		fmt.Print("rate laptop (y/n)? ")
		var answer string
		fmt.Scan(&answer)

		if strings.ToLower(answer) != "y" {
			break
		}

		for i := 0; i < n; i++ {
			scores[i] = sample.RandomLaptopScore()
		}

		err := laptopClient.RateLaptop(laptops, scores)
		if err != nil {
			log.Fatal("cannot rate laptop: ", err)
		}
	}
}

const (
	username        = "admin1"
	password        = "secret"
	refreshDuration = 30 * time.Second
)

func accessibleRoles() map[string]bool {
	return map[string]bool{
		pb.LaptopService_CreateLaptop_FullMethodName: true,
		pb.LaptopService_UploadImage_FullMethodName:  true,
		pb.LaptopService_RateLaptop_FullMethodName:   true,
	}
}

func loatTLSCertificates() (credentials.TransportCredentials, error) {
	certPool, err := util.LoadCAPool()
	if err != nil {
		return nil, err
	}

	clientCert, err := tls.LoadX509KeyPair("cert/client-cert.pem", "cert/client-key.pem")
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}

	return credentials.NewTLS(config), nil
}

func main() {
	serverAddress := flag.String("address", "", "the server address")
	operation := flag.String("operation", "", "which operation to do")
	enableTLS := flag.Bool("tls", false, "enable SSL/TLS")
	flag.Parse()
	log.Printf("dial server %s, TLS: %t", *serverAddress, *enableTLS)

	transportOption := grpc.WithInsecure()

	if *enableTLS {
		tlsCredentials, err := loatTLSCertificates()
		if err != nil {
			log.Fatal("cannot load TLS certificates: ", err)
		}

		transportOption = grpc.WithTransportCredentials(tlsCredentials)
	}

	cc1, err := grpc.Dial(*serverAddress, transportOption)
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	authClient := auth.NewAuthClient(cc1, username, password)
	interceptor, err := auth.NewAuthInterceptor(authClient, accessibleRoles(), refreshDuration)
	if err != nil {
		log.Fatal("cannot create auth interceptor: ", err)
	}

	cc2, err := grpc.Dial(
		*serverAddress, transportOption,
		grpc.WithUnaryInterceptor(interceptor.Unary()),
		grpc.WithStreamInterceptor(interceptor.Stream()),
	)
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	laptopClient := auth.NewLaptopClient(cc2)

	switch *operation {
	case "create":
		testcreateLaptop(laptopClient)
	case "search":
		testSearchLaptop(laptopClient)
	case "upload":
		testUploadImage(laptopClient)
	case "rate":
		testRateLaptop(laptopClient)
	default:
		log.Fatal("unknown operation")
	}
}
