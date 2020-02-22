package erutan

// NewNetVector3 constructs a NetVector3 (skip unkeyed thing)
func NewNetVector3(x, y, z float64) *NetVector3 {
	return &NetVector3{X: x, Y: y, Z: z}
}

// TODO: move all vector stuff here or make another proto special "vector/math"
