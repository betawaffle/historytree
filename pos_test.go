package historytree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	maxFrozenCount = MaxNodes - 1
)

func TestInvalidPos(t *testing.T) {
	assert.Panics(t, func() { Node(1, 1) })
	assert.Panics(t, func() { Node(1, 0).Child(1) })
	assert.Panics(t, func() { Node(0, 2).Child(2) })
	assert.Panics(t, func() { (Pos{1, MaxLayer}).FrozenCount() })
	assert.Panics(t, func() { Node(0, MaxLayer).Parent() })
	assert.Panics(t, func() { Node(0, 0).Walk(0) })
}

func TestPosAdjacent(t *testing.T) {
	type row struct {
		parent Pos
		left   Pos
		right  Pos
	}
	tests := [...]row{
		{Node(0, 1), Leaf(0), Leaf(1)},
		{Node(2, 1), Leaf(2), Leaf(3)},
		{Node(4, 1), Leaf(4), Leaf(5)},
		{Node(6, 1), Leaf(6), Leaf(7)},
		{Node(0, 2), Node(0, 1), Node(2, 1)},
		{Node(4, 2), Node(4, 1), Node(6, 1)},
		{Node(0, 3), Node(0, 2), Node(4, 2)},
	}
	for _, test := range tests {
		if actual := test.parent.Child(0); actual != test.left {
			t.Errorf("expected %s.Child(0) to return %s, got %s", test.parent, test.left, actual)
		}
		if actual := test.parent.Child(1); actual != test.right {
			t.Errorf("expected %s.Child(1) to return %s, got %s", test.parent, test.right, actual)
		}
		if actual := test.left.Parent(); actual != test.parent {
			t.Errorf("expected %s.Parent() to return %s, got %s", test.left, test.parent, actual)
		}
		if actual := test.left.Sibling(); actual != test.right {
			t.Errorf("expected %s.Sibling() to return %s, got %s", test.left, test.right, actual)
		}
		if actual := test.right.Parent(); actual != test.parent {
			t.Errorf("expected %s.Parent() to return %s, got %s", test.right, test.parent, actual)
		}
		if actual := test.right.Sibling(); actual != test.left {
			t.Errorf("expected %s.Sibling() to return %s, got %s", test.right, test.left, actual)
		}
	}
}

func TestPosFreezing(t *testing.T) {
	type row struct {
		p Pos
		i Index
	}
	tests := map[uint64]row{
		0:  {Leaf(0), 0},
		1:  {Leaf(1), 1},
		2:  {Node(0, 1), 1},
		3:  {Leaf(2), 2},
		4:  {Leaf(3), 3},
		5:  {Node(2, 1), 3},
		6:  {Node(0, 2), 3},
		7:  {Leaf(4), 4},
		8:  {Leaf(5), 5},
		9:  {Node(4, 1), 5},
		10: {Leaf(6), 6},
		11: {Leaf(7), 7},
		12: {Node(6, 1), 7},
		13: {Node(4, 2), 7},
		14: {Node(0, 3), 7},
	}
	for v, test := range tests {
		if actual := test.p.FrozenBy(); actual != test.i {
			t.Errorf("expected %s.FrozenBy() to return %d, got %d", test.p, test.i, actual)
		}
		if actual := test.p.FrozenCount(); actual != v {
			t.Errorf("expected %s.FrozenCount() to return %d, got %d", test.p, v, actual)
		}
		if actual := test.p.Offset(); actual != int64(v*LabelSize) {
			t.Errorf("expected %s.FrozenCount() to return %d, got %d", test.p, v*LabelSize, actual)
		}
	}
	if n := Node(0, MaxLayer).FrozenCount(); n != maxFrozenCount {
		t.Errorf("expected (0,%d).FrozenCount() to return %d, got %d", MaxLayer, maxFrozenCount, n)
	}
	if n := Node(MaxIndex, 0).FrozenCount(); n != maxFrozenCount-maxLayer {
		t.Errorf("expected (%d,0).FrozenCount() to return %d, got %d", MaxIndex, maxFrozenCount-maxLayer, n)
	}
}

func TestPosString(t *testing.T) {
	assert.Equal(t, "(1,0)", Node(1, 0).String())
	assert.Equal(t, "(0,1)", Node(0, 1).String())
}

func TestPosWalk(t *testing.T) {
	type row struct {
		parent Pos
		index  Index
		child  Pos
	}
	tests := [...]row{
		{Node(0, 1), 0, Leaf(0)},
		{Node(0, 1), 1, Leaf(1)},
		{Node(0, 2), 2, Node(2, 1)},
		{Node(0, 2), 3, Node(2, 1)},
		{Node(0, 3), 3, Node(0, 2)},
		{Node(0, 3), 5, Node(4, 2)},
		{Node(0, 3), 5, Node(4, 2)},
		{Node(4, 2), 5, Node(4, 1)},

		// now for the hard ones...
		{Node(4, 2), 2, Node(0, 3)},
	}
	for _, test := range tests {
		if actual := test.parent.Walk(test.index); actual != test.child {
			t.Errorf("expected %s.Descend(%d) to return %s, got %s", test.parent, test.index, test.child, actual)
		}
	}
}

func BenchmarkPosFrozenCount(b *testing.B) {
	p := Node(0, MaxIndex.SubtreeDepth()) // root of the largest tree?
	b.ResetTimer()

	var n uint64
	for i := 0; i < b.N; i++ {
		n = p.FrozenCount()
	}
	if n != maxFrozenCount {
		b.Fatalf("expected %s.FrozenCount() to return %d, got %d", p, maxFrozenCount, n)
	}
}
