package sample

import (
	"github.com/caiofernandes00/playing-with-golang/grpc/pkg/proto/pb"
	"github.com/golang/protobuf/ptypes"
)

func NewKeyboard() *pb.Keyboard {
	return &pb.Keyboard{
		Layout:  randomKeyboardLayout(),
		Backlit: randomBool(),
	}
}

func NewCpu() *pb.CPU {
	brand := randomCPUBrand()
	name := randomCpuName(brand)

	numberOfCores := randomInt(2, 8)
	numberOfThreads := randomInt(numberOfCores, 12)

	minGz := randomFloat64(2.0, 3.5)
	maxGz := randomFloat64(minGz, 5.0)

	return &pb.CPU{
		Brand:         brand,
		Name:          name,
		NumberCores:   uint32(numberOfCores),
		NumberThreads: uint32(numberOfThreads),
		MinGhz:        minGz,
		MaxGhz:        maxGz,
	}
}

func NewGpu() *pb.GPU {
	brand := randomGpuBrand()
	name := randomGpuName(brand)

	minGz := randomFloat64(1.0, 1.5)
	maxGz := randomFloat64(minGz, 2.0)

	memory := &pb.Memory{
		Value: uint64(randomInt(2, 8)),
		Unit:  pb.Memory_GIGABYTE,
	}

	return &pb.GPU{
		Brand:  brand,
		Name:   name,
		MinGhz: minGz,
		MaxGhz: maxGz,
		Memory: memory,
	}
}

func NewRam() *pb.Memory {
	return &pb.Memory{
		Value: uint64(randomInt(4, 64)),
		Unit:  pb.Memory_GIGABYTE,
	}
}

func NewSSD() *pb.Storage {
	return &pb.Storage{
		Driver: pb.Storage_SSD,
		Memory: &pb.Memory{
			Value: uint64(randomInt(64, 128)),
			Unit:  pb.Memory_GIGABYTE,
		},
	}
}

func NewHDD() *pb.Storage {
	return &pb.Storage{
		Driver: pb.Storage_HDD,
		Memory: &pb.Memory{
			Value: uint64(randomInt(1, 6)),
			Unit:  pb.Memory_TERABYTE,
		},
	}
}

func NewScreen() *pb.Screen {
	return &pb.Screen{
		SizeInch:   randomFloat32(13, 17),
		Resolution: randomScreenResolution(),
		Painel:     randomScreenPainel(),
	}
}

func NewLaptop() *pb.Laptop {
	brand := randomLaptopBrand()
	name := randomLaptopName(brand)

	return &pb.Laptop{
		Id:       randomID(),
		Brand:    brand,
		Name:     name,
		Cpu:      NewCpu(),
		Ram:      NewRam(),
		Gpus:     []*pb.GPU{NewGpu()},
		Storages: []*pb.Storage{NewSSD(), NewHDD()},
		Screen:   NewScreen(),
		Keyboard: NewKeyboard(),
		Weight: &pb.Laptop_WeightKg{
			WeightKg: randomFloat64(1.0, 3.0),
		},
		PriceUsd:    randomFloat64(1500, 3000),
		ReleaseYear: uint32(randomInt(2015, 2021)),
		UpdatedAt:   ptypes.TimestampNow(),
	}
}

func RandomLaptopScore() float64 {
	return float64(randomInt(1, 10))
}
