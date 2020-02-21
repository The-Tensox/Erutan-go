package octree

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	erutan "github.com/user/erutan/protos/realtime"
	"github.com/user/erutan/utils/vector"
)

// From https://github.com/benbjohnson/testing
// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func Test_BoxContainsPoints(t *testing.T) {
	b := vector.Box{
		Min: erutan.NetVector3{X: 0, Y: 0, Z: 0},
		Max: erutan.NetVector3{X: 1, Y: 1, Z: 1},
	}

	equals(t, true, b.ContainsPoint(&erutan.NetVector3{X: 0, Y: 0, Z: 0}))
	equals(t, true, b.ContainsPoint(&erutan.NetVector3{X: 1, Y: 0, Z: 0}))
	equals(t, true, b.ContainsPoint(&erutan.NetVector3{X: 0, Y: 0, Z: 1}))
	equals(t, true, b.ContainsPoint(&erutan.NetVector3{X: 0.5, Y: 0.5, Z: 0.5}))
	equals(t, true, b.ContainsPoint(&erutan.NetVector3{X: 0, Y: 0, Z: 0}))
	equals(t, false, b.ContainsPoint(&erutan.NetVector3{X: -0.000001, Y: 0.5, Z: 0.5}))
	equals(t, false, b.ContainsPoint(&erutan.NetVector3{X: 0.5, Y: -0.000001, Z: 0.5}))
	equals(t, false, b.ContainsPoint(&erutan.NetVector3{X: 0.5, Y: 0.5, Z: -0.000001}))
	equals(t, false, b.ContainsPoint(&erutan.NetVector3{X: 1.000001, Y: 0.5, Z: 0.5}))
	equals(t, false, b.ContainsPoint(&erutan.NetVector3{X: 0.5, Y: 1.000001, Z: 0.5}))
	equals(t, false, b.ContainsPoint(&erutan.NetVector3{X: 0.5, Y: 0.5, Z: 1.000001}))
}

func Test_BoxContainsBox(t *testing.T) {
	b := vector.Box{
		Min: erutan.NetVector3{X: 0, Y: 0, Z: 0},
		Max: erutan.NetVector3{X: 1, Y: 1, Z: 1},
	}
	var b2 vector.Box

	// contains equal vector.Box, symmetrically
	b2 = vector.Box{erutan.NetVector3{X: 0, Y: 0, Z: 0}, erutan.NetVector3{X: 1, Y: 1, Z: 1}}
	equals(t, true, b2.Contains(&b))
	equals(t, true, b.Contains(&b2))
	equals(t, true, b2.IsContainedIn(&b))
	equals(t, true, b.IsContainedIn(&b2))

	// contained on edge
	b2 = vector.Box{erutan.NetVector3{X: 0, Y: 0, Z: 0}, erutan.NetVector3{X: 0.5, Y: 1, Z: 1}}
	equals(t, true, b.Contains(&b2))
	equals(t, false, b2.Contains(&b))
	equals(t, false, b.IsContainedIn(&b2))
	equals(t, true, b2.IsContainedIn(&b))

	// contained away from edges
	b2 = vector.Box{erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1}, erutan.NetVector3{X: 0.9, Y: 0.9, Z: 0.9}}
	equals(t, true, b.Contains(&b2))
	equals(t, false, b2.Contains(&b))
	equals(t, false, b.IsContainedIn(&b2))
	equals(t, true, b2.IsContainedIn(&b))

	// 1 corner inside
	b2 = vector.Box{erutan.NetVector3{X: -0.1, Y: -0.1, Z: -0.1}, erutan.NetVector3{X: 0.9, Y: 0.9, Z: 0.9}}
	equals(t, false, b.Contains(&b2))
	equals(t, false, b2.Contains(&b))
	equals(t, false, b.IsContainedIn(&b2))
	equals(t, false, b2.IsContainedIn(&b))

	b2 = vector.Box{erutan.NetVector3{X: 0.9, Y: 0.9, Z: 0.9}, erutan.NetVector3{X: 1.1, Y: 1.1, Z: 1.1}}
	equals(t, false, b.Contains(&b2))
	equals(t, false, b2.Contains(&b))
	equals(t, false, b.IsContainedIn(&b2))
	equals(t, false, b2.IsContainedIn(&b))
}

func Test_BoxIntersectsBox(t *testing.T) {
	b := vector.Box{
		Min: erutan.NetVector3{X: 0, Y: 0, Z: 0},
		Max: erutan.NetVector3{X: 1, Y: 1, Z: 1},
	}
	var b2 vector.Box

	// not intersecting vector.Box above or below in each dimension
	b2 = vector.Box{erutan.NetVector3{X: 1.1, Y: 0, Z: 0}, erutan.NetVector3{X: 2, Y: 1, Z: 1}}
	equals(t, false, b.Intersects(&b2))
	b2 = vector.Box{erutan.NetVector3{X: -1, Y: 0, Z: 0}, erutan.NetVector3{X: -0.1, Y: 1, Z: 1}}
	equals(t, false, b.Intersects(&b2))
	b2 = vector.Box{erutan.NetVector3{X: 0, Y: 1.1, Z: 0}, erutan.NetVector3{X: 1, Y: 2, Z: 1}}
	equals(t, false, b.Intersects(&b2))
	b2 = vector.Box{erutan.NetVector3{X: 0, Y: -1, Z: 0}, erutan.NetVector3{X: 1, Y: -0.1, Z: 1}}
	equals(t, false, b.Intersects(&b2))
	b2 = vector.Box{erutan.NetVector3{X: 0, Y: 0, Z: 1.1}, erutan.NetVector3{X: 1, Y: 1, Z: 2}}
	equals(t, false, b.Intersects(&b2))
	b2 = vector.Box{erutan.NetVector3{X: 0, Y: 0, Z: -1}, erutan.NetVector3{X: 1, Y: 1, Z: -0.1}}
	equals(t, false, b.Intersects(&b2))

	// intersects equal vector.Box, symmetrically
	b2 = vector.Box{erutan.NetVector3{X: 0, Y: 0, Z: 0}, erutan.NetVector3{X: 1, Y: 1, Z: 1}}
	equals(t, true, b.Intersects(&b2))
	equals(t, true, b2.Intersects(&b))

	// intersects containing and contained
	b2 = vector.Box{erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1}, erutan.NetVector3{X: 0.9, Y: 0.9, Z: 0.9}}
	equals(t, true, b.Intersects(&b2))
	equals(t, true, b2.Intersects(&b))

	// intersects partial containment on each corner
	b2 = vector.Box{erutan.NetVector3{X: 0.9, Y: 0.9, Z: 0.9}, erutan.NetVector3{X: 2, Y: 2, Z: 2}}
	equals(t, true, b.Intersects(&b2))
	b2 = vector.Box{erutan.NetVector3{X: -1, Y: 0.9, Z: 0.9}, erutan.NetVector3{X: 0.1, Y: 2, Z: 2}}
	equals(t, true, b.Intersects(&b2))
	b2 = vector.Box{erutan.NetVector3{X: 0.9, Y: -1, Z: 0.9}, erutan.NetVector3{X: 2, Y: 0.1, Z: 2}}
	equals(t, true, b.Intersects(&b2))
	b2 = vector.Box{erutan.NetVector3{X: -1, Y: -1, Z: 0.9}, erutan.NetVector3{X: 0.1, Y: 0.1, Z: 2}}
	equals(t, true, b.Intersects(&b2))
	b2 = vector.Box{erutan.NetVector3{X: 0.9, Y: 0.9, Z: -1}, erutan.NetVector3{X: 2, Y: 2, Z: 0.1}}
	equals(t, true, b.Intersects(&b2))
	b2 = vector.Box{erutan.NetVector3{X: -1, Y: 0.9, Z: -1}, erutan.NetVector3{X: 0.1, Y: 2, Z: 0.1}}
	equals(t, true, b.Intersects(&b2))
	b2 = vector.Box{erutan.NetVector3{X: 0.9, Y: -1, Z: -1}, erutan.NetVector3{X: 2, Y: 0.1, Z: 0.1}}
	equals(t, true, b.Intersects(&b2))
	b2 = vector.Box{erutan.NetVector3{X: -1, Y: -1, Z: -1}, erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1}}
	equals(t, true, b.Intersects(&b2))

	// intersects 'beam'; where no corners inside
	// other but some contained
	b2 = vector.Box{erutan.NetVector3{X: -1, Y: 0.1, Z: 0.1}, erutan.NetVector3{X: 2, Y: 0.9, Z: 0.9}}
	equals(t, true, b.Intersects(&b2))
	b2 = vector.Box{erutan.NetVector3{X: 0.1, Y: -1, Z: 0.1}, erutan.NetVector3{X: 0.9, Y: 2, Z: 0.9}}
	equals(t, true, b.Intersects(&b2))
	b2 = vector.Box{erutan.NetVector3{X: 0.1, Y: 0.1, Z: -1}, erutan.NetVector3{X: 0.9, Y: 0.9, Z: 2}}
	equals(t, true, b.Intersects(&b2))
}

func TestInitializesRoot(t *testing.T) {
	o := NewOctree(vector.Box{erutan.NetVector3{X: 0, Y: 0, Z: 0}, erutan.NetVector3{X: 1, Y: 1, Z: 1}})

	//o.Add(99, erutan.NetVector3{10, 0, 0})

	equals(t, true, o.root.point == nil)
	equals(t, erutan.NetVector3{X: 0, Y: 0, Z: 0}, o.root.box.Min)
	equals(t, erutan.NetVector3{X: 1, Y: 1, Z: 1}, o.root.box.Max)
	equals(t, false, o.root.hasChildren)
	equals(t, 0, len(o.root.children))
}

func TestInsertsContainedElements(t *testing.T) {
	o := NewOctree(vector.Box{erutan.NetVector3{X: 0, Y: 0, Z: 0}, erutan.NetVector3{X: 1, Y: 1, Z: 1}})

	equals(t, true, o.Add(99, erutan.NetVector3{X: 1.00000000001, Y: 1, Z: 1}) == nil)
	equals(t, false, o.root.hasChildren)
	equals(t, true, o.root.children == nil)
	equals(t, true, o.root.point == nil)

	equals(t, true, o.Add(99, erutan.NetVector3{X: -0.0000000001, Y: 0, Z: 0}) == nil)
	equals(t, false, o.root.hasChildren)
	equals(t, true, o.root.children == nil)
	equals(t, true, o.root.point == nil)

	equals(t, false, o.Add(88, erutan.NetVector3{X: 0.5, Y: 0, Z: 0}) == nil)
	equals(t, false, o.root.hasChildren)
	equals(t, true, o.root.children == nil)
	equals(t, false, o.root.point == nil)
}

func TestEqualPointsSubdivide(t *testing.T) {
	o := NewOctree(vector.Box{erutan.NetVector3{X: 0, Y: 0, Z: 0}, erutan.NetVector3{X: 1, Y: 1, Z: 1}})

	o.Add(1, erutan.NetVector3{X: 0, Y: 0, Z: 0})
	o.Add(1, erutan.NetVector3{X: 0, Y: 0, Z: 0})
	equals(t, false, o.root.hasChildren)
	equals(t, true, o.root.children == nil)
	equals(t, true, vector.ApproxEqual(*o.root.point, erutan.NetVector3{X: 0, Y: 0, Z: 0}))
	o.Add(1, erutan.NetVector3{X: 1, Y: 1, Z: 1})
	equals(t, true, o.root.hasChildren)
	equals(t, false, o.root.children == nil)
	equals(t, true, o.root.point == nil)
}

func TestRetrievesElementsIn(t *testing.T) {
	o := NewOctree(vector.Box{erutan.NetVector3{X: 0, Y: 0, Z: 0}, erutan.NetVector3{X: 1, Y: 1, Z: 1}})

	o.Add(11, erutan.NetVector3{X: 0, Y: 0, Z: 0})
	// contains point
	equals(t, 1, len(o.ElementsIn(vector.Box{erutan.NetVector3{X: -1, Y: -1, Z: -1}, erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1}})))
	// 0 size at point
	equals(t, 1, len(o.ElementsIn(vector.Box{erutan.NetVector3{X: 0, Y: 0, Z: 0}, erutan.NetVector3{X: 0, Y: 0, Z: 0}})))
	// contains vector.Box
	equals(t, 1, len(o.ElementsIn(vector.Box{erutan.NetVector3{X: -1, Y: -1, Z: -1}, erutan.NetVector3{X: 2, Y: 2, Z: 2}})))

	// coincident point
	o.Add(12, erutan.NetVector3{X: 0, Y: 0, Z: 0})
	// contains point
	equals(t, 2, len(o.ElementsIn(vector.Box{erutan.NetVector3{X: -1, Y: -1, Z: -1}, erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1}})))
	// 0 size at point
	equals(t, 2, len(o.ElementsIn(vector.Box{erutan.NetVector3{X: 0, Y: 0, Z: 0}, erutan.NetVector3{X: 0, Y: 0, Z: 0}})))
	// contains vector.Box
	equals(t, 2, len(o.ElementsIn(vector.Box{erutan.NetVector3{X: -1, Y: -1, Z: -1}, erutan.NetVector3{X: 2, Y: 2, Z: 2}})))

	// non-coincident point
	o.Add(2, erutan.NetVector3{X: 1, Y: 1, Z: 1})
	equals(t, true, o.root.hasChildren)
	equals(t, false, o.root.children == nil)
	equals(t, 8, len(o.root.children))
	equals(t, true, o.root.point == nil)

	// contains point
	equals(t, 2, len(o.ElementsIn(vector.Box{erutan.NetVector3{X: -1, Y: -1, Z: -1}, erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1}})))
	// 0 size at point
	equals(t, 2, len(o.ElementsIn(vector.Box{erutan.NetVector3{X: 0, Y: 0, Z: 0}, erutan.NetVector3{X: 0, Y: 0, Z: 0}})))
	// contains vector.Box
	equals(t, 3, len(o.ElementsIn(vector.Box{erutan.NetVector3{X: -1, Y: -1, Z: -1}, erutan.NetVector3{X: 2, Y: 2, Z: 2}})))

	// fresh octree
	o = NewOctree(vector.Box{erutan.NetVector3{X: 0, Y: 0, Z: 0}, erutan.NetVector3{X: 1, Y: 1, Z: 1}})
	equals(t, false, o.root.hasChildren)

	o.Add(11, erutan.NetVector3{X: 0.4, Y: 0.4, Z: 0.4})
	equals(t, 1, len(o.ElementsIn(vector.Box{erutan.NetVector3{X: -1, Y: -1, Z: -1}, erutan.NetVector3{X: 2, Y: 2, Z: 2}})))
	equals(t, 1, len(o.ElementsIn(vector.Box{erutan.NetVector3{X: 0.4, Y: 0.4, Z: 0.4}, erutan.NetVector3{X: 0.6, Y: 0.6, Z: 0.6}})))
	equals(t, 1, len(o.ElementsIn(vector.Box{erutan.NetVector3{X: -1, Y: 0.4, Z: 0.4}, erutan.NetVector3{X: 1, Y: 0.6, Z: 0.6}})))

	o.Add(12, erutan.NetVector3{X: 0.68, Y: 0.69, Z: 0.7})
	equals(t, 1, len(o.ElementsIn(vector.Box{erutan.NetVector3{X: -1, Y: 0.4, Z: 0.4}, erutan.NetVector3{X: 1, Y: 0.6, Z: 0.6}})))
	equals(t, 2, len(o.ElementsIn(vector.Box{erutan.NetVector3{X: -1, Y: -1, Z: -1}, erutan.NetVector3{X: 1, Y: 1, Z: 1}})))

	// add coincident point in octree
	o.Add(13, erutan.NetVector3{X: 0.68, Y: 0.69, Z: 0.7})
	equals(t, 3, len(o.ElementsIn(vector.Box{erutan.NetVector3{X: -1, Y: -1, Z: -1}, erutan.NetVector3{X: 1, Y: 1, Z: 1}})))
	equals(t, 2, len(o.ElementsIn(vector.Box{erutan.NetVector3{X: 0.68, Y: 0.69, Z: 0.7}, erutan.NetVector3{X: 0.68, Y: 0.69, Z: 0.7}})))
	equals(t, 1, len(o.ElementsIn(vector.Box{erutan.NetVector3{X: 0.35, Y: 0.35, Z: 0.35}, erutan.NetVector3{X: 0.45, Y: 0.45, Z: 0.45}})))

	o.Add(14, erutan.NetVector3{X: 0.1, Y: 0.9, Z: 0.1})

	// values
	equals(t, 11, o.ElementsIn(vector.Box{erutan.NetVector3{X: 0.35, Y: 0.35, Z: 0.35}, erutan.NetVector3{X: 0.45, Y: 0.45, Z: 0.45}})[0])
	equals(t, 12, o.ElementsIn(vector.Box{erutan.NetVector3{X: 0.65, Y: 0.65, Z: 0.65}, erutan.NetVector3{X: 0.75, Y: 0.75, Z: 0.75}})[0])
	equals(t, 13, o.ElementsIn(vector.Box{erutan.NetVector3{X: 0.65, Y: 0.65, Z: 0.65}, erutan.NetVector3{X: 0.75, Y: 0.75, Z: 0.75}})[1])

	// fresh octree
	o = NewOctree(vector.Box{erutan.NetVector3{X: 0, Y: 0, Z: 0}, erutan.NetVector3{X: 1, Y: 1, Z: 1}})
	equals(t, false, o.root.hasChildren)

	a := erutan.NetVector3{X: 0.4, Y: 0.4, Z: 0.4}
	b := a
	b.Y += 0.4
	o.Add(1, a)
	o.Add(2, b)

	// From a to a.Y+0.3
	equals(t, 1, len(o.ElementsIn(vector.Box{a, erutan.NetVector3{X: 0.4, Y: 0.7, Z: 0.4}})))
	equals(t, 2, len(o.ElementsIn(vector.Box{a, b})))
}

func TestRetrievesFirstElementIn(t *testing.T) {
	o := NewOctree(vector.Box{erutan.NetVector3{X: 0, Y: 0, Z: 0}, erutan.NetVector3{X: 1, Y: 1, Z: 1}})
	b := vector.Box{erutan.NetVector3{X: 0, Y: 0, Z: 0}, erutan.NetVector3{X: 1, Y: 1, Z: 1}}
	o.Add(1, erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})

	e := o.FirstElementIn(b)
	equals(t, 1, e)

	o.Add(2, erutan.NetVector3{X: 0.2, Y: 0.2, Z: 0.2})

	e = o.FirstElementIn(b)
	equals(t, 1, e)

	o.Add(3, erutan.NetVector3{X: 0.01, Y: 0.01, Z: 0.01})

	e = o.FirstElementIn(b)
	equals(t, 3, e)
}

func TestRetrievesElementsAt(t *testing.T) {
	o := NewOctree(vector.Box{erutan.NetVector3{X: 0, Y: 0, Z: 0}, erutan.NetVector3{X: 1, Y: 1, Z: 1}})

	o.Add(11, erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})
	// finds element at point
	equals(t, 1, len(o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})))
	equals(t, 11, (o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})[0]))

	// coincident point with different value
	o.Add(12, erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})

	// finds elements at point
	equals(t, 2, len(o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})))
	equals(t, 11, (o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})[0]))
	equals(t, 12, (o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})[1]))

	// finds elements at point after subdivision
	o.Add(13, erutan.NetVector3{X: 0.7, Y: 0.7, Z: 0.7})
	equals(t, 2, len(o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})))
	equals(t, 11, (o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})[0]))
	equals(t, 12, (o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})[1]))
	equals(t, 1, len(o.ElementsAt(erutan.NetVector3{X: 0.7, Y: 0.7, Z: 0.7})))
	equals(t, 13, (o.ElementsAt(erutan.NetVector3{X: 0.7, Y: 0.7, Z: 0.7})[0]))

	// finds elements at point after multiple subdivisions
	o.Add(14, erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.2})
	equals(t, 2, len(o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})))
	equals(t, 11, (o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})[0]))
	equals(t, 12, (o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})[1]))
	equals(t, 1, len(o.ElementsAt(erutan.NetVector3{X: 0.7, Y: 0.7, Z: 0.7})))
	equals(t, 13, (o.ElementsAt(erutan.NetVector3{X: 0.7, Y: 0.7, Z: 0.7})[0]))
	equals(t, 1, len(o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.2})))
	equals(t, 14, (o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.2})[0]))
}

func TestRemovesElements(t *testing.T) {
	o := NewOctree(vector.Box{erutan.NetVector3{X: 0, Y: 0, Z: 0}, erutan.NetVector3{X: 1, Y: 1, Z: 1}})

	// removes element
	o.Add(11, erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})
	equals(t, 1, len(o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})))
	equals(t, 11, (o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})[0]))
	equals(t, true, o.Remove(11))
	equals(t, 0, len(o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})))

	// remove correct element
	o.Add(11, erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})
	equals(t, 1, len(o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})))
	equals(t, 11, (o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})[0]))
	o.Add(12, erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})
	equals(t, 2, len(o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})))
	equals(t, 11, (o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})[0]))
	equals(t, 12, (o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})[1]))
	equals(t, true, o.Remove(11))
	equals(t, false, o.Remove(11))
	equals(t, 1, len(o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})))
	equals(t, 12, (o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})[0]))
	equals(t, true, o.Remove(12))
	equals(t, false, o.Remove(12))
	equals(t, 0, len(o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})))

	o.Add(11, erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})
	o.Add(12, erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})
	o.Add(13, erutan.NetVector3{X: 0.7, Y: 0.7, Z: 0.7})
	equals(t, 2, len(o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})))
	equals(t, 11, (o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})[0]))
	equals(t, 12, (o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})[1]))
	equals(t, 1, len(o.ElementsAt(erutan.NetVector3{X: 0.7, Y: 0.7, Z: 0.7})))
	equals(t, 13, (o.ElementsAt(erutan.NetVector3{X: 0.7, Y: 0.7, Z: 0.7})[0]))
	equals(t, true, o.Remove(11))
	equals(t, false, o.Remove(11))
	equals(t, 1, len(o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})))
	equals(t, 12, (o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})[0]))
	equals(t, 1, len(o.ElementsAt(erutan.NetVector3{X: 0.7, Y: 0.7, Z: 0.7})))
	equals(t, 13, (o.ElementsAt(erutan.NetVector3{X: 0.7, Y: 0.7, Z: 0.7})[0]))
	equals(t, true, o.Remove(12))
	equals(t, false, o.Remove(12))
	equals(t, 0, len(o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})))
	equals(t, 1, len(o.ElementsAt(erutan.NetVector3{X: 0.7, Y: 0.7, Z: 0.7})))
	equals(t, 13, (o.ElementsAt(erutan.NetVector3{X: 0.7, Y: 0.7, Z: 0.7})[0]))
	equals(t, true, o.Remove(13))
	equals(t, false, o.Remove(13))
	equals(t, 0, len(o.ElementsAt(erutan.NetVector3{X: 0.7, Y: 0.7, Z: 0.7})))
}

func TestRemovesElementsUsing(t *testing.T) {
	o := NewOctree(vector.Box{erutan.NetVector3{X: 0, Y: 0, Z: 0}, erutan.NetVector3{X: 1, Y: 1, Z: 1}})

	// removes element using node ref
	node11 := o.Add(11, erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})
	equals(t, 1, len(o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})))
	equals(t, 11, (o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})[0]))
	equals(t, true, o.RemoveUsing(11, node11))
	equals(t, 0, len(o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})))

	// removes element after subdivision using node ref
	node11 = o.Add(11, erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})
	node12 := o.Add(12, erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})
	node13 := o.Add(13, erutan.NetVector3{X: 0.7, Y: 0.7, Z: 0.7})
	node13b := o.Add(13, erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.2})
	equals(t, 2, len(o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})))
	equals(t, 11, (o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})[0]))
	equals(t, 12, (o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})[1]))
	equals(t, 1, len(o.ElementsAt(erutan.NetVector3{X: 0.7, Y: 0.7, Z: 0.7})))
	equals(t, true, o.RemoveUsing(13, node13))
	equals(t, false, o.RemoveUsing(13, node13))
	equals(t, 1, len(o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.2})))
	equals(t, 13, (o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.2})[0]))
	equals(t, true, o.RemoveUsing(13, node13b))
	equals(t, false, o.RemoveUsing(13, node13b))
	equals(t, 0, len(o.ElementsAt(erutan.NetVector3{X: 0.7, Y: 0.7, Z: 0.7})))
	equals(t, 2, len(o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})))
	equals(t, 11, (o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})[0]))
	equals(t, 12, (o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})[1]))
	equals(t, true, o.RemoveUsing(11, node11))
	equals(t, false, o.RemoveUsing(11, node11))
	equals(t, 1, len(o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})))
	equals(t, 12, (o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})[0]))
	equals(t, true, o.RemoveUsing(12, node12))
	equals(t, false, o.RemoveUsing(12, node12))
	equals(t, 0, len(o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})))
}

func TestClearTree(t *testing.T) {
	o := NewOctree(vector.Box{erutan.NetVector3{X: 0, Y: 0, Z: 0}, erutan.NetVector3{X: 1, Y: 1, Z: 1}})
	equals(t, 0, len(o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})))
	o.Add(11, erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})
	equals(t, 1, len(o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})))

	o.Clear()
	equals(t, 0, len(o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1})))
}

func TestRaycast(t *testing.T) {
	o := NewOctree(vector.Box{erutan.NetVector3{X: 0, Y: 0, Z: 0}, erutan.NetVector3{X: 1, Y: 1, Z: 1}})
	a := erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.1}
	b := a
	b.Y += 0.4 // Just above a
	o.Add(1, a)
	o.Add(2, b)
	origin := a
	origin.Y += 0.1 // We start the raycast just from above our position, to skip ourself ...
	hit := o.Raycast(origin, erutan.NetVector3{X: 0, Y: 1, Z: 0}, 1)
	/*

				b 0.1,0.5,0.1
		d=0.4	^ Raycast of length 1
				|
				a 0.1,0.1,0.1

		We're supposed to hit b

	*/
	equals(t, 2, hit)

	hit = o.Raycast(origin, erutan.NetVector3{X: 0, Y: 1, Z: 0}, 0.1)
	/*

				b 0.1,0.5,0.1
		d=0.4	^ Raycast of length 0.1
				|
				a 0.1,0.1,0.1

		We're not supposed to hit anything  there

	*/
	equals(t, nil, hit)

	hit = o.Raycast(origin, erutan.NetVector3{X: 0, Y: 1, Z: 0}, 0.4)
	/*

				b 0.1,0.5,0.1
		d=0.4	^ Raycast of length 0.4
				|
				a 0.1,0.1,0.1

		We're supposed to hit b (on the edge of the raycast)

	*/
	equals(t, 2, hit)
}
