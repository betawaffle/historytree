package historytree

const (
	// MaxLayers is the maximum number of layers supported. It is 1-indexed.
	MaxLayers = 58 // 1-indexed

	// MaxLayer is the highest supported layer. It is 0-indexed.
	MaxLayer = Layer(maxLayer)
	maxLayer = MaxLayers - 1 // untyped version
)

// Layer represents the layer component of a node, or an offset from the base
// of the tree.
type Layer uint8

// MaxIndex returns the highest index in a full subtree rooted at (0,r).
func (r Layer) MaxIndex() Index {
	return 1<<r - 1
}

// SubtreeCount returns the number of nodes in a full subtree rooted at (0,r).
func (r Layer) SubtreeCount() uint64 {
	return 1<<r - 1 | 1<<r
}

func (r Layer) assertValid() {
	if r > MaxLayer {
		panic("historytree: layer out of range")
	}
}
