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
	b := vector.NewBox(erutan.NewNetVector3(0, 0, 0), erutan.NewNetVector3(1, 1, 1))

	equals(t, true, b.ContainsPoint(erutan.NewNetVector3(0, 0, 0)))
	equals(t, true, b.ContainsPoint(erutan.NewNetVector3(1, 0, 0)))
	equals(t, true, b.ContainsPoint(erutan.NewNetVector3(0, 0, 1)))
	equals(t, true, b.ContainsPoint(erutan.NewNetVector3(0.5, 0.5, 0.5)))
	equals(t, true, b.ContainsPoint(erutan.NewNetVector3(0, 0, 0)))
	equals(t, false, b.ContainsPoint(erutan.NewNetVector3(-0.000001, 0.5, 0.5)))
	equals(t, false, b.ContainsPoint(erutan.NewNetVector3(0.5, -0.000001, 0.5)))
	equals(t, false, b.ContainsPoint(erutan.NewNetVector3(0.5, 0.5, -0.000001)))
	equals(t, false, b.ContainsPoint(erutan.NewNetVector3(1.000001, 0.5, 0.5)))
	equals(t, false, b.ContainsPoint(erutan.NewNetVector3(0.5, 1.000001, 0.5)))
	equals(t, false, b.ContainsPoint(erutan.NewNetVector3(0.5, 0.5, 1.000001)))
}

func Test_BoxContainsBox(t *testing.T) {
	b := vector.NewBox(erutan.NewNetVector3(0, 0, 0), erutan.NewNetVector3(1, 1, 1))
	b2 := vector.NewBox(erutan.NewNetVector3(0, 0, 0), erutan.NewNetVector3(1, 1, 1))

	// contains equal vector.Box, symmetrically
	equals(t, true, b2.Contains(b))
	equals(t, true, b.Contains(b2))
	equals(t, true, b2.IsContainedIn(b))
	equals(t, true, b.IsContainedIn(b2))

	// contained on edge
	b2 = vector.NewBox(erutan.NewNetVector3(0, 0, 0), erutan.NewNetVector3(0.5, 1, 1))

	equals(t, true, b.Contains(b2))
	equals(t, false, b2.Contains(b))
	equals(t, false, b.IsContainedIn(b2))
	equals(t, true, b2.IsContainedIn(b))

	// contained away from edges
	b2 = vector.NewBox(erutan.NewNetVector3(0.1, 0.1, 0.1), erutan.NewNetVector3(0.9, 0.9, 0.9))
	equals(t, true, b.Contains(b2))
	equals(t, false, b2.Contains(b))
	equals(t, false, b.IsContainedIn(b2))
	equals(t, true, b2.IsContainedIn(b))

	// 1 corner inside
	b2 = vector.NewBox(erutan.NewNetVector3(-0.1, -0.1, -0.1), erutan.NewNetVector3(0.9, 0.9, 0.9))
	equals(t, false, b.Contains(b2))
	equals(t, false, b2.Contains(b))
	equals(t, false, b.IsContainedIn(b2))
	equals(t, false, b2.IsContainedIn(b))

	b2 = vector.NewBox(erutan.NewNetVector3(0.9, 0.9, 0.9), erutan.NewNetVector3(1.1, 1.1, 1.1))
	equals(t, false, b.Contains(b2))
	equals(t, false, b2.Contains(b))
	equals(t, false, b.IsContainedIn(b2))
	equals(t, false, b2.IsContainedIn(b))
}

func Test_BoxIntersectsBox(t *testing.T) {
	b := vector.NewBox(erutan.NewNetVector3(0, 0, 0), erutan.NewNetVector3(1, 1, 1))

	// not intersecting vector.Box above or below in each dimension
	b2 := vector.NewBox(erutan.NewNetVector3(1.1, 0, 0), erutan.NewNetVector3(2, 1, 1))
	equals(t, false, b.Intersects(b2))
	b2 = vector.NewBox(erutan.NewNetVector3(-1, 0, 0), erutan.NewNetVector3(-0.1, 1, 1))
	equals(t, false, b.Intersects(b2))
	b2 = vector.NewBox(erutan.NewNetVector3(0, 1.1, 0), erutan.NewNetVector3(1, 2, 1))
	equals(t, false, b.Intersects(b2))
	b2 = vector.NewBox(erutan.NewNetVector3(0, -1, 0), erutan.NewNetVector3(1, -0.1, 1))
	equals(t, false, b.Intersects(b2))
	b2 = vector.NewBox(erutan.NewNetVector3(0, 0, 1.1), erutan.NewNetVector3(1, 1, 2))
	equals(t, false, b.Intersects(b2))
	b2 = vector.NewBox(erutan.NewNetVector3(0, 0, -1), erutan.NewNetVector3(1, 1, -0.1))
	equals(t, false, b.Intersects(b2))

	// intersects equal vector.Box, symmetrically
	b2 = vector.NewBox(erutan.NewNetVector3(0, 0, 0), erutan.NewNetVector3(1, 1, 1))
	equals(t, true, b.Intersects(b2))
	equals(t, true, b2.Intersects(b))

	// intersects containing and contained
	b2 = vector.NewBox(erutan.NewNetVector3(0.1, 0.1, 0.1), erutan.NewNetVector3(0.9, 0.9, 0.9))

	equals(t, true, b.Intersects(b2))
	equals(t, true, b2.Intersects(b))

	// intersects partial containment on each corner
	b2 = vector.NewBox(erutan.NewNetVector3(0.9, 0.9, 0.9), erutan.NewNetVector3(2, 2, 2))

	equals(t, true, b.Intersects(b2))
	b2 = vector.NewBox(erutan.NewNetVector3(-1, 0.9, 0.9), erutan.NewNetVector3(1, 2, 2))

	equals(t, true, b.Intersects(b2))
	b2 = vector.NewBox(erutan.NewNetVector3(0.9, -1, 0.9), erutan.NewNetVector3(2, 0.1, 2))

	equals(t, true, b.Intersects(b2))
	b2 = vector.NewBox(erutan.NewNetVector3(-1, -1, 0.9), erutan.NewNetVector3(0.1, 0.1, 2))

	equals(t, true, b.Intersects(b2))
	b2 = vector.NewBox(erutan.NewNetVector3(0.9, 0.9, -1), erutan.NewNetVector3(2, 2, 0.1))

	equals(t, true, b.Intersects(b2))
	b2 = vector.NewBox(erutan.NewNetVector3(-1, 0.9, -1), erutan.NewNetVector3(0.1, 2, 0.1))

	equals(t, true, b.Intersects(b2))
	b2 = vector.NewBox(erutan.NewNetVector3(0.9, -1, -1), erutan.NewNetVector3(2, 0.1, 0.1))

	equals(t, true, b.Intersects(b2))
	b2 = vector.NewBox(erutan.NewNetVector3(-1, -1, -1), erutan.NewNetVector3(0.1, 0.1, 0.1))

	equals(t, true, b.Intersects(b2))

	// intersects 'beam'; where no corners inside
	// other but some contained
	b2 = vector.NewBox(erutan.NewNetVector3(-1, 0.1, 0.1), erutan.NewNetVector3(2, 0.9, 0.9))

	equals(t, true, b.Intersects(b2))
	b2 = vector.NewBox(erutan.NewNetVector3(0.1, -1, 0.1), erutan.NewNetVector3(0.9, 2, 0.9))

	equals(t, true, b.Intersects(b2))
	b2 = vector.NewBox(erutan.NewNetVector3(0.1, 0.1, -1), erutan.NewNetVector3(0.9, 0.9, 2))

	equals(t, true, b.Intersects(b2))
}

func TestInitializesRoot(t *testing.T) {
	o := NewOctree(vector.NewBox(erutan.NewNetVector3(0, 0, 0), erutan.NewNetVector3(1, 1, 1)))

	equals(t, true, o.root.point == nil)
	equals(t, *erutan.NewNetVector3(0, 0, 0), o.root.box.Min)
	equals(t, *erutan.NewNetVector3(1, 1, 1), o.root.box.Max)
	equals(t, false, o.root.hasChildren)
	equals(t, 0, len(o.root.children))
}

func TestInsertsContainedElements(t *testing.T) {
	o := NewOctree(vector.NewBox(erutan.NewNetVector3(0, 0, 0), erutan.NewNetVector3(1, 1, 1)))

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
	o := NewOctree(vector.NewBox(erutan.NewNetVector3(0, 0, 0), erutan.NewNetVector3(1, 1, 1)))

	o.Add(1, *erutan.NewNetVector3(0, 0, 0))
	o.Add(1, *erutan.NewNetVector3(0, 0, 0))
	equals(t, false, o.root.hasChildren)
	equals(t, true, o.root.children == nil)
	equals(t, true, vector.ApproxEqual(*o.root.point, *erutan.NewNetVector3(0, 0, 0)))
	o.Add(1, *erutan.NewNetVector3(1, 1, 1))
	equals(t, true, o.root.hasChildren)
	equals(t, false, o.root.children == nil)
	equals(t, true, o.root.point == nil)
}

func TestRetrievesElementsIn(t *testing.T) {
	o := NewOctree(vector.NewBox(erutan.NewNetVector3(0, 0, 0), erutan.NewNetVector3(1, 1, 1)))

	o.Add(11, *erutan.NewNetVector3(0, 0, 0))
	// contains point
	equals(t, 1, len(o.ElementsIn(*vector.NewBox(erutan.NewNetVector3(-1, -1, -1), erutan.NewNetVector3(0.1, 0.1, 0.1)))))
	// 0 size at point
	equals(t, 1, len(o.ElementsIn(*vector.NewBox(erutan.NewNetVector3(0, 0, 0), erutan.NewNetVector3(0, 0, 0)))))
	// contains vector.Box
	equals(t, 1, len(o.ElementsIn(*vector.NewBox(erutan.NewNetVector3(-1, -1, -1), erutan.NewNetVector3(2, 2, 2)))))

	// coincident point
	o.Add(12, *erutan.NewNetVector3(0, 0, 0))
	// contains point
	equals(t, 2, len(o.ElementsIn(*vector.NewBox(erutan.NewNetVector3(-1, -1, -1), erutan.NewNetVector3(0.1, 0.1, 0.1)))))
	// 0 size at point
	equals(t, 2, len(o.ElementsIn(*vector.NewBox(erutan.NewNetVector3(0, 0, 0), erutan.NewNetVector3(0, 0, 0)))))
	// contains vector.Box
	equals(t, 2, len(o.ElementsIn(*vector.NewBox(erutan.NewNetVector3(-1, -1, -1), erutan.NewNetVector3(2, 2, 2)))))

	// non-coincident point
	o.Add(2, *erutan.NewNetVector3(1, 1, 1))
	equals(t, true, o.root.hasChildren)
	equals(t, false, o.root.children == nil)
	equals(t, 8, len(o.root.children))
	equals(t, true, o.root.point == nil)

	// contains point
	equals(t, 2, len(o.ElementsIn(*vector.NewBox(erutan.NewNetVector3(-1, -1, -1), erutan.NewNetVector3(0.1, 0.1, 0.1)))))
	// 0 size at point
	equals(t, 2, len(o.ElementsIn(*vector.NewBox(erutan.NewNetVector3(0, 0, 0), erutan.NewNetVector3(0, 0, 0)))))
	// contains vector.Box
	equals(t, 3, len(o.ElementsIn(*vector.NewBox(erutan.NewNetVector3(-1, -1, -1), erutan.NewNetVector3(2, 2, 2)))))

	// fresh octree
	o = NewOctree(vector.NewBox(erutan.NewNetVector3(0, 0, 0), erutan.NewNetVector3(1, 1, 1)))
	equals(t, false, o.root.hasChildren)

	o.Add(11, *erutan.NewNetVector3(0.4, 0.4, 0.4))
	equals(t, 1, len(o.ElementsIn(*vector.NewBox(erutan.NewNetVector3(-1, -1, -1), erutan.NewNetVector3(2, 2, 2)))))
	equals(t, 1, len(o.ElementsIn(*vector.NewBox(erutan.NewNetVector3(0.4, 0.4, 0.4), erutan.NewNetVector3(0.6, 0.6, 0.6)))))
	equals(t, 1, len(o.ElementsIn(*vector.NewBox(erutan.NewNetVector3(-1, 0.4, 0.4), erutan.NewNetVector3(1, 0.6, 0.6)))))

	o.Add(12, *erutan.NewNetVector3(0.68, 0.69, 0.7))
	equals(t, 1, len(o.ElementsIn(*vector.NewBox(erutan.NewNetVector3(-1, 0.4, 0.4), erutan.NewNetVector3(1, 0.6, 0.6)))))
	equals(t, 2, len(o.ElementsIn(*vector.NewBox(erutan.NewNetVector3(-1, -1, -1), erutan.NewNetVector3(1, 1, 1)))))

	// add coincident point in octree
	o.Add(13, *erutan.NewNetVector3(0.68, 0.69, 0.7))
	equals(t, 3, len(o.ElementsIn(*vector.NewBox(erutan.NewNetVector3(-1, -1, -1), erutan.NewNetVector3(1, 1, 1)))))
	equals(t, 2, len(o.ElementsIn(*vector.NewBox(erutan.NewNetVector3(0.68, 0.69, 0.7), erutan.NewNetVector3(0.68, 0.69, 0.7)))))
	equals(t, 1, len(o.ElementsIn(*vector.NewBox(erutan.NewNetVector3(0.35, 0.35, 0.35), erutan.NewNetVector3(0.45, 0.45, 0.45)))))

	o.Add(14, *erutan.NewNetVector3(0.1, 0.9, 0.1))

	// values
	equals(t, 11, o.ElementsIn(*vector.NewBox(erutan.NewNetVector3(0.35, 0.35, 0.35), erutan.NewNetVector3(0.45, 0.45, 0.45)))[0])
	equals(t, 12, o.ElementsIn(*vector.NewBox(erutan.NewNetVector3(0.65, 0.65, 0.65), erutan.NewNetVector3(0.75, 0.75, 0.75)))[0])
	equals(t, 13, o.ElementsIn(*vector.NewBox(erutan.NewNetVector3(0.65, 0.65, 0.65), erutan.NewNetVector3(0.75, 0.75, 0.75)))[1])

	// fresh octree
	o = NewOctree(vector.NewBox(erutan.NewNetVector3(0, 0, 0), erutan.NewNetVector3(1, 1, 1)))
	equals(t, false, o.root.hasChildren)

	o.Add(1, *erutan.NewNetVector3(0.4, 0.4, 0.4))
	o.Add(2, *erutan.NewNetVector3(0.4, 0.8, 0.4))

	// From a to a.Y+0.3
	equals(t, 1, len(o.ElementsIn(*vector.NewBox(erutan.NewNetVector3(0.4, 0.4, 0.4), erutan.NewNetVector3(0.4, 0.7, 0.4)))))

	// From a to b
	equals(t, 2, len(o.ElementsIn(*vector.NewBox(erutan.NewNetVector3(0.4, 0.4, 0.4), erutan.NewNetVector3(0.4, 0.8, 0.4)))))

	o = NewOctree(vector.NewBox(erutan.NewNetVector3(-1, -1, -1), erutan.NewNetVector3(1, 1, 1)))

	o.Add(3, *erutan.NewNetVector3(0, 0.1, 0))
	o.Add(4, *erutan.NewNetVector3(0, -0.5, 0))

	// From 0 to -1
	equals(t, 4, o.ElementsIn(*vector.NewBox(erutan.NewNetVector3(0, 0, 0), erutan.NewNetVector3(0, -1, 0)))[0])
	equals(t, 2, len(o.ElementsIn(*vector.NewBox(erutan.NewNetVector3(0, 0.1, 0), erutan.NewNetVector3(0, -0.5, 0)))))
}

func TestRetrievesFirstElementIn(t *testing.T) {
	o := NewOctree(vector.NewBox(erutan.NewNetVector3(0, 0, 0), erutan.NewNetVector3(1, 1, 1)))
	b := *vector.NewBox(erutan.NewNetVector3(0, 0, 0), erutan.NewNetVector3(1, 1, 1))
	o.Add(1, *erutan.NewNetVector3(0.1, 0.1, 0.1))

	e := o.FirstElementIn(b)
	equals(t, 1, e)

	o.Add(2, *erutan.NewNetVector3(0.2, 0.2, 0.2))

	e = o.FirstElementIn(b)
	equals(t, 1, e)

	o.Add(3, *erutan.NewNetVector3(0.01, 0.01, 0.01))

	e = o.FirstElementIn(b)
	equals(t, 3, e)
}

func TestRetrievesElementsAt(t *testing.T) {
	o := NewOctree(vector.NewBox(erutan.NewNetVector3(0, 0, 0), erutan.NewNetVector3(1, 1, 1)))

	o.Add(11, *erutan.NewNetVector3(0.1, 0.1, 0.1))
	// finds element at point
	equals(t, 1, len(o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))))
	equals(t, 11, (o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))[0]))

	// coincident point with different value
	o.Add(12, *erutan.NewNetVector3(0.1, 0.1, 0.1))

	// finds elements at point
	equals(t, 2, len(o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))))
	equals(t, 11, (o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))[0]))
	equals(t, 12, (o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))[1]))

	// finds elements at point after subdivision
	o.Add(13, *erutan.NewNetVector3(0.7, 0.7, 0.7))
	equals(t, 2, len(o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))))
	equals(t, 11, (o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))[0]))
	equals(t, 12, (o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))[1]))
	equals(t, 1, len(o.ElementsAt(*erutan.NewNetVector3(0.7, 0.7, 0.7))))
	equals(t, 13, (o.ElementsAt(*erutan.NewNetVector3(0.7, 0.7, 0.7))[0]))

	// finds elements at point after multiple subdivisions
	o.Add(14, erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.2})
	equals(t, 2, len(o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))))
	equals(t, 11, (o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))[0]))
	equals(t, 12, (o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))[1]))
	equals(t, 1, len(o.ElementsAt(*erutan.NewNetVector3(0.7, 0.7, 0.7))))
	equals(t, 13, (o.ElementsAt(*erutan.NewNetVector3(0.7, 0.7, 0.7))[0]))
	equals(t, 1, len(o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.2})))
	equals(t, 14, (o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.2})[0]))
}

func TestRemovesElements(t *testing.T) {
	o := NewOctree(vector.NewBox(erutan.NewNetVector3(0, 0, 0), erutan.NewNetVector3(1, 1, 1)))

	// removes element
	o.Add(11, *erutan.NewNetVector3(0.1, 0.1, 0.1))
	equals(t, 1, len(o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))))
	equals(t, 11, (o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))[0]))
	equals(t, true, o.Remove(11))
	equals(t, 0, len(o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))))

	// remove correct element
	o.Add(11, *erutan.NewNetVector3(0.1, 0.1, 0.1))
	equals(t, 1, len(o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))))
	equals(t, 11, (o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))[0]))
	o.Add(12, *erutan.NewNetVector3(0.1, 0.1, 0.1))
	equals(t, 2, len(o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))))
	equals(t, 11, (o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))[0]))
	equals(t, 12, (o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))[1]))
	equals(t, true, o.Remove(11))
	equals(t, false, o.Remove(11))
	equals(t, 1, len(o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))))
	equals(t, 12, (o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))[0]))
	equals(t, true, o.Remove(12))
	equals(t, false, o.Remove(12))
	equals(t, 0, len(o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))))

	o.Add(11, *erutan.NewNetVector3(0.1, 0.1, 0.1))
	o.Add(12, *erutan.NewNetVector3(0.1, 0.1, 0.1))
	o.Add(13, *erutan.NewNetVector3(0.7, 0.7, 0.7))
	equals(t, 2, len(o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))))
	equals(t, 11, (o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))[0]))
	equals(t, 12, (o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))[1]))
	equals(t, 1, len(o.ElementsAt(*erutan.NewNetVector3(0.7, 0.7, 0.7))))
	equals(t, 13, (o.ElementsAt(*erutan.NewNetVector3(0.7, 0.7, 0.7))[0]))
	equals(t, true, o.Remove(11))
	equals(t, false, o.Remove(11))
	equals(t, 1, len(o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))))
	equals(t, 12, (o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))[0]))
	equals(t, 1, len(o.ElementsAt(*erutan.NewNetVector3(0.7, 0.7, 0.7))))
	equals(t, 13, (o.ElementsAt(*erutan.NewNetVector3(0.7, 0.7, 0.7))[0]))
	equals(t, true, o.Remove(12))
	equals(t, false, o.Remove(12))
	equals(t, 0, len(o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))))
	equals(t, 1, len(o.ElementsAt(*erutan.NewNetVector3(0.7, 0.7, 0.7))))
	equals(t, 13, (o.ElementsAt(*erutan.NewNetVector3(0.7, 0.7, 0.7))[0]))
	equals(t, true, o.Remove(13))
	equals(t, false, o.Remove(13))
	equals(t, 0, len(o.ElementsAt(*erutan.NewNetVector3(0.7, 0.7, 0.7))))
}

func TestRemovesElementsUsing(t *testing.T) {
	o := NewOctree(vector.NewBox(erutan.NewNetVector3(0, 0, 0), erutan.NewNetVector3(1, 1, 1)))

	// removes element using node ref
	node11 := o.Add(11, *erutan.NewNetVector3(0.1, 0.1, 0.1))
	equals(t, 1, len(o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))))
	equals(t, 11, (o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))[0]))
	equals(t, true, o.RemoveUsing(11, node11))
	equals(t, 0, len(o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))))

	// removes element after subdivision using node ref
	node11 = o.Add(11, *erutan.NewNetVector3(0.1, 0.1, 0.1))
	node12 := o.Add(12, *erutan.NewNetVector3(0.1, 0.1, 0.1))
	node13 := o.Add(13, *erutan.NewNetVector3(0.7, 0.7, 0.7))
	node13b := o.Add(13, *erutan.NewNetVector3(0.1, 0.1, 0.2))
	equals(t, 2, len(o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))))
	equals(t, 11, (o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))[0]))
	equals(t, 12, (o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))[1]))
	equals(t, 1, len(o.ElementsAt(*erutan.NewNetVector3(0.7, 0.7, 0.7))))
	equals(t, true, o.RemoveUsing(13, node13))
	equals(t, false, o.RemoveUsing(13, node13))
	equals(t, 1, len(o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.2})))
	equals(t, 13, (o.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.1, Z: 0.2})[0]))
	equals(t, true, o.RemoveUsing(13, node13b))
	equals(t, false, o.RemoveUsing(13, node13b))
	equals(t, 0, len(o.ElementsAt(*erutan.NewNetVector3(0.7, 0.7, 0.7))))
	equals(t, 2, len(o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))))
	equals(t, 11, (o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))[0]))
	equals(t, 12, (o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))[1]))
	equals(t, true, o.RemoveUsing(11, node11))
	equals(t, false, o.RemoveUsing(11, node11))
	equals(t, 1, len(o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))))
	equals(t, 12, (o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))[0]))
	equals(t, true, o.RemoveUsing(12, node12))
	equals(t, false, o.RemoveUsing(12, node12))
	equals(t, 0, len(o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))))
}

func TestClearTree(t *testing.T) {
	o := NewOctree(vector.NewBox(erutan.NewNetVector3(0, 0, 0), erutan.NewNetVector3(1, 1, 1)))
	equals(t, 0, len(o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))))
	o.Add(11, *erutan.NewNetVector3(0.1, 0.1, 0.1))
	equals(t, 1, len(o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))))

	o.Clear()
	equals(t, 0, len(o.ElementsAt(*erutan.NewNetVector3(0.1, 0.1, 0.1))))
}

func TestRaycast(t *testing.T) {
	o := NewOctree(vector.NewBox(erutan.NewNetVector3(0, 0, 0), erutan.NewNetVector3(1, 1, 1)))
	o.Add(1, *erutan.NewNetVector3(0.1, 0.1, 0.1))
	o.Add(2, *erutan.NewNetVector3(0.1, 0.5, 0.1))
	hit := o.Raycast(*erutan.NewNetVector3(0.1, 0.2, 0.1), *erutan.NewNetVector3(0, 1, 0), 1)
	/*

				b 0.1,0.5,0.1
		d=0.4	^ Raycast of length 1
				|
				a 0.1,0.1,0.1

		We're supposed to hit b

	*/
	equals(t, 2, hit)

	hit = o.Raycast(*erutan.NewNetVector3(0.1, 0.2, 0.1), *erutan.NewNetVector3(0, 1, 0), 0.1)
	/*

				b 0.1,0.5,0.1
		d=0.4	^ Raycast of length 0.1
				|
				a 0.1,0.1,0.1

		We're not supposed to hit anything  there

	*/
	equals(t, nil, hit)

	hit = o.Raycast(*erutan.NewNetVector3(0.1, 0.2, 0.1), *erutan.NewNetVector3(0, 1, 0), 0.4)
	/*

				b 0.1,0.5,0.1
		d=0.4	^ Raycast of length 0.4
				|
				a 0.1,0.1,0.1

		We're supposed to hit b (on the edge of the raycast)

	*/
	equals(t, 2, hit)

	o = NewOctree(vector.NewBox(erutan.NewNetVector3(-1, -1, -1), erutan.NewNetVector3(1, 1, 1)))
	o.Add(3, *erutan.NewNetVector3(0, 0.1, 0))
	o.Add(4, *erutan.NewNetVector3(0, -0.4, 0))
	hit = o.Raycast(*erutan.NewNetVector3(0, 0, 0), *erutan.NewNetVector3(0, -1, 0), 0.5)
	/*

				c 0,0.1,0
				|Raycast of length 0.5
				-
				d 0,-0.3,0

		We're supposed to hit 4

	*/
	equals(t, 4, hit)
}
