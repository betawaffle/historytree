package historytree

const (
	// MaxIndex is the highest supported leaf index.
	MaxIndex = Index(maxIndex)
	maxIndex = 1<<maxLayer - 1 // untyped version
)

// Index represents the a node index or a version of the tree.
type Index uint64

// SubtreeDepth returns the 0-indexed depth of the subtree containing (i,0).
func (i Index) SubtreeDepth() (d Layer) {
	// FIXME: There's a faster way to count trailing ones, right?
	for i&1 != 0 {
		d++
		i >>= 1
	}
	return
}

// TreeDepth returns the minimum depth for a tree containing (i,0).
func (i Index) TreeDepth() (d Layer) {
	// FIXME: This is a poor man's log2.
	for i != 0 {
		d++
		i >>= 1
	}
	return
}

func (i Index) assertValid() {
	if i > MaxIndex {
		panic("historytree: index out of range")
	}
}
