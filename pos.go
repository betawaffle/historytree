package historytree

import "fmt"

const (
	// MaxNodes is the maximum number of nodes that can be in a single tree.
	MaxNodes = 1<<MaxLayers - 1

	// confirm that MaxNodes fits in uint64
	_ = uint64(MaxNodes)
)

// Pos represents the location of a node in the tree.
type Pos struct {
	i Index
	r Layer
}

// Node returns a Pos at the given index and layer.
// If i or r are invalid, or i is invalid at r, Node panics.
func Node(i Index, r Layer) (p Pos) {
	p.i = i
	p.r = r
	p.assertValid()
	return
}

// Leaf returns the Pos of a leaf node at i.
func Leaf(i Index) Pos {
	return Node(i, 0)
}

// Child returns a Pos representing the left (0) or right (1) child of p.
// Child may panic if p is a leaf (layer 0), or if i > 1.
func (p Pos) Child(i Index) Pos {
	// if debugMode {
	p.assertValid()

	if p.r == 0 {
		panic("historytree: layer 0 has no children")
	}
	if i > 1 {
		panic("historytree: child index out of range")
	}
	// }
	r := p.r - 1
	return Pos{
		i: p.i + i<<r,
		r: r,
	}
}

// FrozenBy returns the version where p would be frozen.
func (p Pos) FrozenBy() Index {
	p.assertValid()

	return p.i + p.r.MaxIndex()
}

// FrozenCount returns the number nodes frozen before p.
func (p Pos) FrozenCount() uint64 {
	n := uint64(p.r)
	i := p.i + (1<<p.r - 1) // = p.FrozenBy()

	if i > MaxIndex || i < p.i {
		panic("historytree: index out of range")
	}

	// for each bit set in the index, add the subtree
	for r, bit := Layer(0), Index(1); i != 0; bit = bit << 1 {
		r++
		if i&bit != 0 {
			i, n = i&^bit, n+(1<<r-1) // = n + r.MaxIndex()
		}
	}
	return n
}

// Offset returns the byte offset of p in an append-only log.
func (p Pos) Offset() int64 {
	n := p.FrozenCount()
	return int64(n) * LabelSize
}

// Parent returns a Pos representing the parent of p.
// Parent may panic if the layer would overflow.
func (p Pos) Parent() Pos {
	p.assertValid()

	if p.r == MaxLayer {
		panic("historytree: parent layer out of range")
	}

	r := p.r + 1
	return Pos{
		i: p.i &^ r.MaxIndex(), // clear all the descendant bits
		r: r,
	}
}

// Sibling returns the opposite child under the same parent.
func (p Pos) Sibling() Pos {
	return Pos{
		i: p.i ^ (1 << p.r), // flip the bit at this layer
		r: p.r,
	}
}

func (p Pos) String() string {
	return fmt.Sprintf("(%d,%d)", p.i, p.r)
}

// Walk returns the adjacent node closest to (i,0).
// If (i,0) is not a descendent of p, it will return the parent.
func (p Pos) Walk(i Index) Pos {
	if i&p.i != p.i {
		return p.Parent()
	}

	p.assertValid()

	if p.r == 0 {
		panic("historytree: layer 0 has no children")
	}

	r := p.r - 1
	return Pos{
		i: p.i | 1<<r&i,
		r: r,
	}
}

func (p Pos) assertValid() {
	p.i.assertValid()
	p.r.assertValid()

	if p.i&p.r.MaxIndex() != 0 {
		panic("historytree: invalid index for layer")
	}
}
