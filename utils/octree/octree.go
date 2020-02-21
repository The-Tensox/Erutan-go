package octree

import (
	"fmt"
	"math"

	erutan "github.com/user/erutan/protos/realtime"
	"github.com/user/erutan/utils"
	"github.com/user/erutan/utils/vector"
)

// Octree An octree is a data structure that allows fast retrieval of elements based
// values in three dimensions.
type Octree struct {
	root *Node
}

// NewOctree Makes a new octree with the given box.
func NewOctree(box vector.Box) *Octree {
	return &Octree{
		root: &Node{
			box: box,
		},
	}
}

// Clear Removes all the data from the Octree while
// retaining its bounding box. Returns true if octree is ready for use
// (because it has previously been initialized).
func (o *Octree) Clear() bool {
	if o.root != nil {
		// if octree has been initializes, use the same box,
		// but create a new root, freeing the other memory
		// (except where outside references have been retained).
		o.root = &Node{box: o.root.box}
		return true
	}

	return false
}

// Range iterates octree elements
func (o *Octree) Range(f func(elements []interface{})) {
	// TODO: remove or fix
	utils.DebugLogf("e %v", o.root.elements)
	return
	for _, c := range o.root.children {
		f(o.ElementsIn(c.box))
	}
}

// Add Inserts the element in the tree at the specified point.
// If you may need to remove the element later, retain the
// returned node for fast removal.
func (o *Octree) Add(element interface{}, point erutan.NetVector3) *Node {
	return o.root.tryAdd([]interface{}{element}, &point)
}

// ElementsAt Retrieves a slice of elements that exist at
// the specified point in the tree.
func (o *Octree) ElementsAt(point erutan.NetVector3) []interface{} {
	return o.root.elementsAt(&point)
}

// ElementsIn Retrieves a slice of element that exist
// within the specified box.
func (o *Octree) ElementsIn(box vector.Box) []interface{} {
	return o.root.elementsIn(&box)
}

// FirstElementIn Retrieves an element that exist
// within the specified box. Should be faster than ElementsIn if we just want
// the first element
func (o *Octree) FirstElementIn(box vector.Box) interface{} {
	return o.root.firstElementIn(&box)
}

// Remove Removes the specified element from the tree.
// Generally, RemoveUsing should used as it is faster under
// most circumstances.
func (o *Octree) Remove(element interface{}) bool {
	return o.root.remove(element)
}

// RemoveUsing Removes the specified element from the tree; node constrains the search
// for the element and should usually be the node returned when this element
// was placed in the tree using Add()
func (o *Octree) RemoveUsing(element interface{}, node *Node) bool {
	if node != nil {
		return node.remove(element)
	}
	return false
}

// Raycast will try to find interfaces in a direction from an origin with a maximum distance
func (o *Octree) Raycast(origin erutan.NetVector3, direction erutan.NetVector3, maxDistance float64) interface{} {
	// Default value in go ...
	if maxDistance == -1 {
		maxDistance = math.MaxFloat64
	}
	b := vector.Box{
		Min: origin,
		Max: vector.Add(origin, vector.Mul(direction, maxDistance)),
	}
	//utils.DebugLogf("box %v -> %v", b.Min.Y, b.Max.Y)
	// TODO: optimization: should just stop at first element met, we don't care for others
	// Implement sort of "firstelementin"
	hits := o.ElementsIn(b)

	if len(hits) == 0 {
		return nil
	}
	utils.DebugLogf("Hit something: %v", hits[0])
	return hits[0]
}

// Spherecast TODO ...
func (o *Octree) Spherecast() {
}

// ToString Get a human readable representation of the state of
// this octree and its contents.
func (o *Octree) ToString() string {
	str := "nil"
	if o.root != nil {
		str = o.root.recursiveToString("  ", "  ")
	}

	return fmt.Sprintf("Octree{\n  root: %v\n}", str)
}

// Node An element within the tree that can either act as a leaf, that can directly hold a point
// and its corresponding elements or act as a branch and hold references to child nodes.
type Node struct {
	box         vector.Box
	point       *erutan.NetVector3
	elements    []interface{}
	hasChildren bool
	children    []*Node
}

func (n *Node) tryAdd(elements []interface{}, point *erutan.NetVector3) *Node {
	// attempt to add the elements in this node (or a descendant)
	// at the specified point.

	if !n.box.ContainsPoint(point) {
		return nil
	}

	if n.hasChildren {
		return n.addToChildren(elements, point)
	}

	if n.point != nil {
		// leaf already has assigned point
		if vector.ApproxEqual(*n.point, *point) {
			// points are equal
			n.elements = append(n.elements, elements...)
			return n
		}

		// subdivide because points are different
		return n.subdivide(elements, point)
	}

	// set own elements and point
	n.elements = elements
	n.point = point

	return n
}

func (n *Node) addToChildren(elements []interface{}, point *erutan.NetVector3) *Node {
	for _, child := range n.children {
		// try adding to child
		leaf := child.tryAdd(elements, point)

		if leaf != nil {
			// succeeded adding
			return leaf
		}
	}

	// Error: box.contains evaluated to true, but none of the children added the point
	return nil
}

func (n *Node) subdivide(addElements []interface{}, atPoint *erutan.NetVector3) *Node {
	// create child nodes for what is currently a leaf,
	// moving its current contents to one of those leafs.

	// setup this node's children
	n.hasChildren = true
	subBoxes := n.box.MakeSubBoxes()

	for i := 0; i < 8; i++ {
		n.children = append(n.children, &Node{box: subBoxes[i]})
	}

	// add node's elements and point to a child
	n.addToChildren(n.elements, n.point)

	// clear elements and point from self
	n.elements = nil
	n.point = nil

	// add the new element to a child
	return n.addToChildren(addElements, atPoint)
}

func (n *Node) elementsAt(point *erutan.NetVector3) []interface{} {
	// get any alements in this node (or a descendant)
	// at the specified point

	if n.hasChildren {
		for _, child := range n.children {
			if child.box.ContainsPoint(point) {
				return child.elementsAt(point)
			}
		}
	} else {
		// when a leaf
		if n.point != nil && vector.ApproxEqual(*point, *n.point) {
			return n.elements
		}
	}

	return nil
}

func (n *Node) elementsIn(box *vector.Box) []interface{} {
	// get any elements in this node (or a descendant)
	// within the specified box

	if n.hasChildren {
		elements := []interface{}{}
		for _, child := range n.children {
			if child.box.IsContainedIn(box) {
				// fully contained
				elements = append(elements, child.elementsIn(&child.box)...)
			} else if child.box.Contains(box) || child.box.Intersects(box) {
				// partially contained
				elements = append(elements, child.elementsIn(box)...)
			}
		}

		return elements
	}

	// when a leaf
	if n.point != nil && box.ContainsPoint(n.point) {
		return n.elements
	}

	return nil
}

func (n *Node) firstElementIn(box *vector.Box) interface{} {
	// get first element in this node (or a descendant)
	// within the specified box

	if n.hasChildren {
		for _, child := range n.children {
			if child.box.IsContainedIn(box) {
				// fully contained
				return child.firstElementIn(&child.box)
			} else if child.box.Contains(box) || child.box.Intersects(box) {
				// partially contained
				return child.firstElementIn(box)
			}
		}
	}

	// when a leaf
	if n.point != nil && box.ContainsPoint(n.point) {
		return n.elements[0] // TODO: is it correct, think of a test to check
	}

	return nil
}

func (n *Node) remove(element interface{}) bool {
	// remove the first instance of the specified element
	// in this node (or in a descendant)

	if n.hasChildren {
		for _, child := range n.children {
			if child.remove(element) {
				return true
			}
		}
		return false
	}

	for idx, val := range n.elements {
		if val == element {
			// remove element from the slice
			n.elements = append(n.elements[:idx], n.elements[idx+1:]...)
			return true
		}
	}
	return false
}

// ToString Get a human readable representation of the state of
// this node and its contents.
func (n *Node) ToString() string {
	return n.recursiveToString("", "  ")
}

func (n *Node) recursiveToString(curIndent, stepIndent string) string {
	singleIndent := curIndent + stepIndent

	// default values
	childStr := "nil"
	pointStr := "nil"
	elementStr := "nil"

	if n.hasChildren {
		doubleIndent := singleIndent + stepIndent

		// accumulate child strings
		childStr = ""
		for i, child := range n.children {
			childStr = childStr + fmt.Sprintf("%v%d: %v,\n", doubleIndent, i, child.recursiveToString(doubleIndent, stepIndent))
		}

		childStr = fmt.Sprintf("[\n%v%v]", childStr, singleIndent)
	}

	if n.point != nil {
		pointStr = vector.ToString(n.point)
	}

	if n.elements != nil {
		// not stringifying elements since their type is unknown
		elementStr = fmt.Sprintf("[%d]", len(n.elements))
	}

	return fmt.Sprintf("Node{\n%vchildren: %v,\n%vbox: %v,\n%vpoint: %v\n%velements: %v,\n%v}", singleIndent, childStr, singleIndent, n.box.ToString(), singleIndent, pointStr, singleIndent, elementStr, curIndent)
}
