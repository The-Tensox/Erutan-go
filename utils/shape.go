package utils

import (
	erutan "github.com/The-Tensox/erutan/protobuf"
	"github.com/The-Tensox/protometry"
)

// TODO: SHOULD THIS BE MOVED TO PROTOMETRY INSTEAD ?
func cubeTris() []int32 {
	return []int32{
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
}

// CreateCube return a shape representing a cube
// Based on http://ilkinulas.github.io/development/unity/2016/04/30/cube-mesh-in-unity3d.html
func CreateCube(sideLength float64) *erutan.Mesh {
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

	return &erutan.Mesh{Vertices: vertices, Tris: cubeTris()}
}

// CreateCubeCenterBased create a cube with it's pivot placed on the center instead of bottom-left
func CreateCubeCenterBased(sideLength float64) *erutan.Mesh {
	halfSide := sideLength / 2
	vertices := []*protometry.VectorN{
		protometry.NewVectorN(-halfSide, -halfSide, -halfSide),
		protometry.NewVectorN(halfSide, -halfSide, -halfSide),
		protometry.NewVectorN(halfSide, halfSide, -halfSide),
		protometry.NewVectorN(-halfSide, halfSide, -halfSide),

		protometry.NewVectorN(-halfSide, halfSide, halfSide),
		protometry.NewVectorN(halfSide, halfSide, halfSide),
		protometry.NewVectorN(halfSide, -halfSide, halfSide),
		protometry.NewVectorN(-halfSide, -halfSide, halfSide),
	}

	return &erutan.Mesh{Vertices: vertices, Tris: cubeTris()}
}


// CreateShapeWithMutation takes two shapes and apply random mutations somewhere
func CreateShapeWithMutation(a erutan.Mesh, b erutan.Mesh) *erutan.Mesh {
	return nil
}
