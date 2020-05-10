package erutan

import (
	"math"
	"testing"
)

func TestDiamondSquareAlgorithm(t *testing.T) {
	ds := DiamondSquareAlgorithm(int(math.Pow(2, 10))+1, 500, 1)
	t.Logf("%v", ds)
}