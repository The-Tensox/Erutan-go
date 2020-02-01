package utils

import (
	"fmt"
	"math"

	erutan "github.com/user/erutan_two/protos/realtime"
)

// ApproxEqual reports whether a and b are equal within a small epsilon.
func ApproxEqual(a, b erutan.NetVector3) bool {
	const epsilon = 1e-16
	return math.Abs(a.X-b.X) < epsilon && math.Abs(a.Y-b.Y) < epsilon && math.Abs(a.Z-b.Z) < epsilon
}

// String returns the vector to string
func String(a erutan.NetVector3) string { return fmt.Sprintf("(%0.24f, %0.24f, %0.24f)", a.X, a.Y, a.Z) }

// Norm returns the vector's norm.
func Norm(a erutan.NetVector3) float64 { return math.Sqrt(Dot(a, a)) }

// Norm2 returns the square of the norm.
func Norm2(a erutan.NetVector3) float64 { return Dot(a, a) }

// Normalize returns a unit vector in the same direction as a.
func Normalize(a erutan.NetVector3) erutan.NetVector3 {
	n2 := Norm2(a)
	if n2 == 0 {
		return erutan.NetVector3{X: 0, Y: 0, Z: 0}
	}
	return Mul(a, 1/math.Sqrt(n2))
}

// Abs returns the vector with nonnegative components.
func Abs(a erutan.NetVector3) erutan.NetVector3 {
	return erutan.NetVector3{X: math.Abs(a.X), Y: math.Abs(a.Y), Z: math.Abs(a.Z)}
}

// Add returns the standard vector sum of a and b.
func Add(a, b erutan.NetVector3) erutan.NetVector3 {
	return erutan.NetVector3{X: a.X + b.X, Y: a.Y + b.Y, Z: a.Z + b.Z}
}

// Sub returns the standard vector difference of a and b.
func Sub(a, b erutan.NetVector3) erutan.NetVector3 {
	return erutan.NetVector3{X: a.X - b.X, Y: a.Y - b.Y, Z: a.Z - b.Z}
}

// Mul returns the standard scalar product of a and m.
func Mul(a erutan.NetVector3, m float64) erutan.NetVector3 {
	return erutan.NetVector3{X: m * a.X, Y: m * a.Y, Z: m * a.Z}
}

// Div returns the standard scalar division of a and m.
func Div(a erutan.NetVector3, m float64) erutan.NetVector3 {
	return erutan.NetVector3{X: a.X / m, Y: a.Y / m, Z: a.Z / m} // TODO: check 0
}

// Dot returns the standard dot product of a and b.
func Dot(a, b erutan.NetVector3) float64 { return a.X*b.X + a.Y*b.Y + a.Z*b.Z }

// Cross returns the standard cross product of a and b.
func Cross(a, b erutan.NetVector3) erutan.NetVector3 {
	return erutan.NetVector3{
		X: a.Y*b.Z - a.Z*b.Y,
		Y: a.Z*b.X - a.X*b.Z,
		Z: a.X*b.Y - a.Y*b.X,
	}
}

// Distance returns the Euclidean distance between a and b.
func Distance(a, b erutan.NetVector3) float64 { return Norm(Sub(a, b)) }

// Angle returns the angle between a and b.
func Angle(a, b erutan.NetVector3) float64 {
	return math.Atan2(Norm(Cross(a, b)), Dot(a, b))
}

// LookAt return a quaternion corresponding to the rotation required to look at the other Vector3
func LookAt(a, b erutan.NetVector3) erutan.NetQuaternion {
	angle := Angle(a, b)
	return erutan.NetQuaternion{X: 0, Y: angle, Z: 0, W: angle}
}

// LookAtTwo ...
func LookAtTwo(from, to erutan.NetVector3) [][]float64 {
	tmp := erutan.NetVector3{X: 0, Y: 1, Z: 0}
	forward := Normalize(Sub(from, to))
	right := Cross(Normalize(tmp), forward)
	up := Cross(forward, right)

	a := make([][]float64, 4)
	for i := range a {
		a[i] = make([]float64, 3)
	}

	a[0][0] = right.X
	a[0][1] = right.Y
	a[0][2] = right.Z
	a[1][0] = up.X
	a[1][1] = up.Y
	a[1][2] = up.Z
	a[2][0] = forward.X
	a[2][1] = forward.Y
	a[2][2] = forward.Z

	a[3][0] = from.X
	a[3][1] = from.Y
	a[3][2] = from.Z

	return a
}

// ToQuaternion ... yaw (Z), pitch (Y), roll (X)
func ToQuaternion(yaw, pitch, roll float64) erutan.NetQuaternion {
	// Abbreviations for the various angular functions
	cy := math.Cos(yaw * 0.5)
	sy := math.Sin(yaw * 0.5)
	cp := math.Cos(pitch * 0.5)
	sp := math.Sin(pitch * 0.5)
	cr := math.Cos(roll * 0.5)
	sr := math.Sin(roll * 0.5)

	var q erutan.NetQuaternion
	q.W = cy*cp*cr + sy*sp*sr
	q.X = 0 //cy*cp*sr - sy*sp*cr
	q.Y = sy*cp*sr + cy*sp*cr
	q.Z = 0 //sy*cp*cr - cy*sp*sr

	return q
}
