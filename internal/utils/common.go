package utils

import (
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"math"
	"math/rand"
)

// RandomString generates a random string of 4 bytes
func RandomString() string {
	str := make([]byte, 4)
	rand.Read(str)
	return fmt.Sprintf("%x", str)
}

func RandFloats(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func GetProtoTime() float64 {
	return float64(ptypes.TimestampNow().Seconds)*math.Pow(10, 9) + float64(ptypes.TimestampNow().Nanos)
}
