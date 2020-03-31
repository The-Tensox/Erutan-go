package utils

import (
	"github.com/The-Tensox/protometry"
	erutan "github.com/user/erutan/protos/realtime"
)

// CreateCube return a shape representing a cube
// Based on http://ilkinulas.github.io/development/unity/2016/04/30/cube-mesh-in-unity3d.html
func CreateCube(sideLength float64) *erutan.Shape {
	vertices := []*protometry.VectorN{
		protometry.NewVectorN(0, 0, 0),
		protometry.NewVectorN(sideLength, 0, 0),
		protometry.NewVectorN(sideLength, sideLength, 0),
		protometry.NewVectorN(0, sideLength, 0),
		protometry.NewVectorN(0, sideLength, sideLength),
		protometry.NewVectorN(sideLength, sideLength, sideLength),
		protometry.NewVectorN(sideLength, 0, sideLength),
		protometry.NewVectorN(0, 0, sideLength),
	}

	tris := []int32{
		0, 2, 1, //face front
		0, 3, 2,
		2, 3, 4, //face top
		2, 4, 5,
		1, 2, 5, //face right
		1, 5, 6,
		0, 7, 4, //face left
		0, 4, 3,
		5, 4, 7, //face back
		5, 7, 6,
		0, 6, 7, //face bottom
		0, 1, 6,
	}

	return &erutan.Shape{Vertices: vertices, Tris: tris}
}

// CreateShapeWithMutation takes two shapes and apply random mutations somewhere
func CreateShapeWithMutation(a erutan.Shape, b erutan.Shape) *erutan.Shape {
	return nil
}
