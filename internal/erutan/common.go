package erutan

import (
	"errors"
	"math/rand"
)

// Copy pasta https://gist.github.com/awilki01/83b65ad852a0ab30192af07cda3d7c0b
////size of grid to generate, note this must be a
////value 2^n+1
func DiamondSquareAlgorithm(terrainPoints int, roughness float64, seed float64) [][]float64 {
	if terrainPoints % 2 != 1 {
		panic(errors.New("must be a power of two plus one"))
	}
	dataSize := terrainPoints // must be a power of two plus one
	data := make([][]float64, dataSize)
	for i := range data {
		data[i] = make([]float64, dataSize)
	}
	rand.Seed(int64(seed))
	data[0][0], data[0][dataSize-1], data[dataSize-1][0], data[dataSize-1][dataSize-1] = seed, seed, seed, seed
	h := roughness //the range (-h -> +h) for the average offset - affects roughness

	for sideLength := dataSize - 1; sideLength >= 2; {

		halfSide := sideLength / 2

		//generate the new square values
		for x := 0; x < dataSize-1; x += sideLength {
			for y := 0; y < dataSize-1; y += sideLength {
				//x, y is upper left corner of square
				//calculate average of existing corners
				avg := data[x][y] + data[x+sideLength][y] + data[x][y+sideLength] + data[x+sideLength][y+sideLength]
				avg /= 4.0

				//center is average plus random offset
				data[x+halfSide][y+halfSide] = avg + (rand.Float64() * 2 * h) - h
			}
		}
		//generate the diamond values
		//since the diamonds are staggered we only move x
		//by half side
		//NOTE: if the data shouldn't wrap then x < DATA_SIZE
		//to generate the far edge values
		for x := 0; x < dataSize-1; x += halfSide {
			//and y is x offset by half a side, but moved by
			//the full side length
			//NOTE: if the data shouldn't wrap then y < DATA_SIZE
			//to generate the far edge values
			for y := (x + halfSide) % sideLength; y < dataSize-1; y += sideLength {
				//x, y is center of diamond
				//note we must use mod  and add DATA_SIZE for subtraction
				//so that we can wrap around the array to find the corners
				avg := data[(x-halfSide+dataSize)%dataSize][y] +
					data[(x+halfSide)%dataSize][y] +
					data[x][(y+halfSide)%dataSize] + //below center
					data[x][(y-halfSide+dataSize)%dataSize] //above center
				avg /= 4.0

				//new value = average plus random offset
				//We calculate random value in range of 2h
				//and then subtract h so the end value is
				//in the range (-h, +h)
				avg = avg + (rand.Float64() * 2 * h) - h
				//update value for center of diamond
				data[x][y] = avg

				//wrap values on the edges, remove
				//this and adjust loop condition above
				//for non-wrapping values.
				if x == 0 {
					data[dataSize-1][y] = avg
				}
				if y == 0 {
					data[x][dataSize-1] = avg
				}
			}
		}
		sideLength /= 2
		h /= 2.0
	}
	return data
}