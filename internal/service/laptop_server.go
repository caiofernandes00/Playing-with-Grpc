package service

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"

	"github.com/caiofernandes00/playing-with-golang/grpc/internal/repository"
	"github.com/caiofernandes00/playing-with-golang/grpc/internal/utils"
	"github.com/caiofernandes00/playing-with-golang/grpc/pkg/proto/pb"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const MAX_IMAGE_SIZE = 1 << 20 // 1 MB

type LaptopServer struct {
	laptopStore repository.LaptopStore
	imageStore  repository.ImageStore
	ratingStore repository.RatingStore
}

func NewLaptopServer(laptopStore repository.LaptopStore, imageStore repository.ImageStore, ratingStore repository.RatingStore) *LaptopServer {
	return &LaptopServer{laptopStore, imageStore, ratingStore}
}

func (server *LaptopServer) CreateLaptop(ctx context.Context, req *pb.CreateLaptopRequest) (*pb.CreateLaptopResponse, error) {
	laptop := req.GetLaptop()
	log.Printf("receive a create-laptop request with id: %s", laptop.GetId())

	if len(laptop.GetId()) > 0 {
		_, err := uuid.Parse(laptop.GetId())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "laptop ID is not a valid UUID: %v", err)
		}
	} else {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot generate a new laptop ID: %v", err)
		}
		laptop.Id = id.String()
	}

	if err := utils.ContextError(ctx); err != nil {
		return nil, err
	}

	err := server.laptopStore.Save(laptop)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, repository.ErrAlreadyExists) {
			code = codes.AlreadyExists
		}
		return nil, status.Errorf(code, "cannot save laptop to the store: %v", err)
	}

	log.Printf("saved laptop with id: %s", laptop.GetId())

	return &pb.CreateLaptopResponse{
		Id: laptop.GetId(),
	}, nil
}

func (server *LaptopServer) SearchLaptop(req *pb.SearchLaptopRequest, stream pb.LaptopService_SearchLaptopServer) error {
	filter := req.GetFilter()
	log.Printf("receive a search-laptop request with filter: %v", filter)

	err := server.laptopStore.Search(stream.Context(), filter, func(laptop *pb.Laptop) error {
		res := &pb.SearchLaptopResponse{Laptop: laptop}

		err := stream.Send(res)
		if err != nil {
			return err
		}

		log.Printf("sent laptop with id: %s", laptop.GetId())
		return nil
	})

	if err != nil {
		return status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	return nil
}

func (server *LaptopServer) UploadImage(stream pb.LaptopService_UploadImageServer) error {
	req, err := stream.Recv()
	if err != nil {
		return utils.LogError(status.Errorf(codes.Unknown, "cannot receive image info: %v", err))
	}

	laptopID := req.GetInfo().GetLaptopId()
	imageType := req.GetInfo().GetImageType()
	log.Printf("receive an upload-image request for laptop %s with image type %s", laptopID, imageType)

	laptop, err := server.laptopStore.Find(laptopID)
	if err != nil {
		return utils.LogError(status.Errorf(codes.Internal, "cannot find laptop: %v", err))
	}

	if laptop == nil {
		return utils.LogError(status.Errorf(codes.NotFound, "laptop ID %s doesn't exist", laptopID))
	}

	imageData := bytes.Buffer{}
	imageSize := 0

	for {
		if err := utils.ContextError(stream.Context()); err != nil {
			return err
		}

		log.Print("waiting to receive more data")

		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				log.Print("no more data")
				break
			}

			return utils.LogError(status.Errorf(codes.Unknown, "cannot receive chunk data: %v", err))
		}

		chunk := req.GetChunkData()
		size := len(chunk)

		log.Printf("receive a chunk with size: %d", size)

		imageSize += size
		if imageSize > MAX_IMAGE_SIZE {
			return utils.LogError(status.Errorf(codes.InvalidArgument, "image is too large: %d > %d", imageSize, MAX_IMAGE_SIZE))
		}

		_, err = imageData.Write(chunk)
		if err != nil {
			return utils.LogError(status.Errorf(codes.Internal, "cannot write chunk data: %v", err))
		}
	}

	imageID, err := server.imageStore.Save(laptopID, imageType, imageData)
	if err != nil {
		return utils.LogError(status.Errorf(codes.Internal, "cannot save image to the store: %v", err))
	}

	res := &pb.UploadImageResponse{
		Id:   imageID,
		Size: uint32(imageSize),
	}

	err = stream.SendAndClose(res)
	if err != nil {
		return utils.LogError(status.Errorf(codes.Unknown, "cannot send response: %v", err))
	}

	return nil
}

func (server *LaptopServer) RateLaptop(stream pb.LaptopService_RateLaptopServer) error {
	for {
		err := utils.ContextError(stream.Context())
		if err != nil {
			return err
		}

		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				log.Print("no more data")
				break
			}

			return utils.LogError(status.Errorf(codes.Unknown, "cannot receive stream request: %v", err))
		}

		laptopID := req.GetLaptopId()
		score := req.GetScore()

		log.Printf("receive a rate-laptop request: id = %s, score = %.2f", laptopID, score)

		found, err := server.laptopStore.Find(laptopID)
		if err != nil {
			return utils.LogError(status.Errorf(codes.Internal, "cannot find laptop: %v", err))
		}
		if found == nil {
			return utils.LogError(status.Errorf(codes.NotFound, "laptop ID %s doesn't exist", laptopID))
		}

		rating, err := server.ratingStore.Add(laptopID, score)
		if err != nil {
			return utils.LogError(status.Errorf(codes.Internal, "cannot add rating to the store: %v", err))
		}

		res := &pb.RateLaptopResponse{
			LaptopId:     laptopID,
			RatedCount:   rating.Count,
			AverageScore: rating.Sum / float64(rating.Count),
		}

		err = stream.Send(res)
		if err != nil {
			return utils.LogError(status.Errorf(codes.Internal, "cannot send stream response: %v", err))
		}
	}

	return nil
}
