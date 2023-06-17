package sample

import (
	"math/rand"
	"time"

	"github.com/caiofernandes00/playing-with-golang/grpc/pkg/proto/pb"
	"github.com/google/uuid"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randomKeyboardLayout() pb.Keyboard_Layout {
	switch rand.Intn(3) {
	case 1:
		return pb.Keyboard_QWERTY
	case 2:
		return pb.Keyboard_QWERTZ
	default:
		return pb.Keyboard_AZERTY
	}
}

func randomLaptopBrand() string {
	return randomStringFromSet("Apple", "Dell", "Lenovo")
}

func randomLaptopName(brand string) string {
	switch brand {
	case "Apple":
		return randomStringFromSet("Macbook Air", "Macbook Pro")
	case "Dell":
		return randomStringFromSet("Latitude", "Vostro", "XPS")
	default:
		return randomStringFromSet("Thinkpad X1", "Thinkpad P1")
	}
}

func randomCpuName(brand string) string {
	if brand == "Intel" {
		return randomStringFromSet("Core i3", "Core i5", "Core i7")
	} else {
		return randomStringFromSet("Ryzen 3", "Ryzen 5", "Ryzen 7")
	}
}

func randomCPUBrand() string {
	return randomStringFromSet("Intel", "AMD")
}

func randomGpuBrand() string {
	return randomStringFromSet("Nvidia", "AMD")
}

func randomScreenPainel() pb.Screen_Painel {
	switch rand.Intn(3) {
	case 1:
		return pb.Screen_IPS
	case 2:
		return pb.Screen_OLED
	default:
		return pb.Screen_UNKNOWN
	}
}

func randomGpuName(brand string) string {
	if brand == "Nvidia" {
		return randomStringFromSet("GTX 1060", "GTX 1070", "GTX 1080")
	} else {
		return randomStringFromSet("RX 570", "RX 580", "RX 590")
	}
}

func randomStringFromSet(a ...string) string {
	n := len(a)

	if n == 0 {
		return ""
	}

	return a[rand.Intn(n)]
}

func randomScreenResolution() *pb.Screen_Resolution {
	height := randomInt(1080, 4320)
	width := height * 16 / 9

	return &pb.Screen_Resolution{
		Height: uint32(height),
		Width:  uint32(width),
	}
}

func randomBool() bool {
	return rand.Intn(2) == 1
}

func randomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func randomFloat64(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func randomFloat32(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

func randomID() string {
	return uuid.New().String()
}
