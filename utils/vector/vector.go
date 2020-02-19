package vector

import (
	"fmt"
	"log"
	"math"

	erutan "github.com/user/erutan/protos/realtime"
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
	//DebugLogf("sub:%v||%v|| %v - %v: %v", a, b, a.X, b.X, a.X-b.X)
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

// Min Returns the a vector where each component is the lesser of the
// corresponding component in this and the specified vector.
func Min(v *erutan.NetVector3, other *erutan.NetVector3) erutan.NetVector3 {
	return erutan.NetVector3{
		X: math.Min(v.X, other.X),
		Y: math.Min(v.Y, other.Y),
		Z: math.Min(v.Z, other.Z),
	}
}

// Max Returns the a vector where each component is the greater of the
// corresponding component in this and the specified vector.
func Max(v *erutan.NetVector3, other *erutan.NetVector3) erutan.NetVector3 {
	return erutan.NetVector3{
		X: math.Max(v.X, other.X),
		Y: math.Max(v.Y, other.Y),
		Z: math.Max(v.Z, other.Z),
	}
}

// Lerp Returns the linear interpolation between two erutan.NetVector3(s).
func Lerp(v *erutan.NetVector3, other *erutan.NetVector3, f float64) erutan.NetVector3 {
	return erutan.NetVector3{
		X: (other.X-v.X)*f + v.X,
		Y: (other.Y-v.Y)*f + v.Y,
		Z: (other.Z-v.Z)*f + v.Z,
	}
}

// ToString Get a human readable representation of the state of
// this vector.
func ToString(v *erutan.NetVector3) string {
	return fmt.Sprintf("erutan.NetVector3{%f, %f, %f}", v.X, v.Y, v.Z)
}

///////// BOX /////////////////

// Box Defines an axis aligned rectangular solid.
type Box struct {
	Min erutan.NetVector3
	Max erutan.NetVector3
}

// Size Returns the dimensions of the Box.
func (b *Box) Size() erutan.NetVector3 {
	return Sub(b.Max, b.Min)
}

// ContainsPoint Returns whether the specified point is contained in this box.
func (b *Box) ContainsPoint(v *erutan.NetVector3) bool {
	return (b.Min.X <= v.X &&
		b.Max.X >= v.X &&
		b.Min.Y <= v.Y &&
		b.Max.Y >= v.Y &&
		b.Min.Z <= v.Z &&
		b.Max.Z >= v.Z)
}

// Contains Returns whether the specified box is contained in this box.
func (b *Box) Contains(o *Box) bool {
	return (b.Min.X <= o.Min.X &&
		b.Max.X >= o.Max.X &&
		b.Min.Y <= o.Min.Y &&
		b.Max.Y >= o.Max.Y &&
		b.Min.Z <= o.Min.Z &&
		b.Max.Z >= o.Max.Z)
}

// IsContainedIn Returns whether the specified box contains this box.
func (b *Box) IsContainedIn(o *Box) bool {
	return o.Contains(b)
}

// Intersects Returns whether any portion of this box intersects with
// the specified box.
func (b *Box) Intersects(o *Box) bool {
	return !(b.Max.X < o.Min.X ||
		o.Max.X < b.Min.X ||
		b.Max.Y < o.Min.Y ||
		o.Max.Y < b.Min.Y ||
		b.Max.Z < o.Min.Z ||
		o.Max.Z < b.Min.Z)
}

// MakeSubBoxes ...
func (b *Box) MakeSubBoxes() [8]Box {
	// gets the child boxes (octants) of the box.
	center := Lerp(&b.Min, &b.Max, 0.5)

	return [8]Box{
		Box{erutan.NetVector3{X: b.Min.Y, Y: b.Min.Y, Z: b.Min.Z}, erutan.NetVector3{X: center.X, Y: center.Y, Z: center.Z}},
		Box{erutan.NetVector3{X: center.Y, Y: b.Min.Y, Z: b.Min.Z}, erutan.NetVector3{X: b.Max.X, Y: center.Y, Z: center.Z}},
		Box{erutan.NetVector3{X: b.Min.Y, Y: center.Y, Z: b.Min.Z}, erutan.NetVector3{X: center.X, Y: b.Max.Y, Z: center.Z}},
		Box{erutan.NetVector3{X: center.Y, Y: center.Y, Z: b.Min.Z}, erutan.NetVector3{X: b.Max.X, Y: b.Max.Y, Z: center.Z}},
		Box{erutan.NetVector3{X: b.Min.Y, Y: b.Min.Y, Z: center.Z}, erutan.NetVector3{X: center.X, Y: center.Y, Z: b.Max.Z}},
		Box{erutan.NetVector3{X: center.Y, Y: b.Min.Y, Z: center.Z}, erutan.NetVector3{X: b.Max.X, Y: center.Y, Z: b.Max.Z}},
		Box{erutan.NetVector3{X: b.Min.Y, Y: center.Y, Z: center.Z}, erutan.NetVector3{X: center.X, Y: b.Max.Y, Z: b.Max.Z}},
		Box{erutan.NetVector3{X: center.Y, Y: center.Y, Z: center.Z}, erutan.NetVector3{X: b.Max.X, Y: b.Max.Y, Z: b.Max.Z}},
	}
}

// GetBox return a box based on position and scale of the object, not really correct,
// should use size instead of scale
func GetBox(position, scale erutan.NetVector3) Box {
	return Box{
		Min: erutan.NetVector3{X: position.X - scale.X/2, Y: position.Y - scale.Y/2, Z: position.Z - scale.Z/2},
		Max: erutan.NetVector3{X: position.X + scale.X/2, Y: position.Y + scale.Y/2, Z: position.Z + scale.Z/2},
	}
}

// MinimumTranslation tells how much an entity has to move to no longer overlap another entity.
// TODO: 3D
func MinimumTranslation(rect1, rect2 Box) erutan.NetVector3 {
	mtd := erutan.NetVector3{}

	left := rect2.Min.X - rect1.Max.X
	right := rect2.Max.X - rect1.Min.X
	top := rect2.Min.Y - rect1.Max.Y
	bottom := rect2.Max.Y - rect1.Min.Y

	if left > 0 || right < 0 {
		log.Println("Box aint intercepting")
		return mtd
		//box doesn't intercept
	}

	if top > 0 || bottom < 0 {
		log.Println("Box aint intercepting")
		return mtd
		//box doesn't intercept
	}
	if math.Abs(left) < right {
		mtd.X = left
	} else {
		mtd.X = right
	}

	if math.Abs(top) < bottom {
		mtd.Y = top
	} else {
		mtd.Y = bottom
	}

	if math.Abs(mtd.X) < math.Abs(mtd.Y) {
		mtd.Y = 0
	} else {
		mtd.X = 0
	}

	return mtd
}

// ToString Get a human readable representation of the state of
// this box.
func (b *Box) ToString() string {
	return fmt.Sprintf("Box{min: %v, max: %v}", ToString(&b.Min), ToString(&b.Max))
}
