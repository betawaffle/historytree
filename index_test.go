package historytree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvalidIndex(t *testing.T) {
	assert.Panics(t, func() { (MaxIndex + 1).assertValid() })
}

func TestIndexSubtreeDepth(t *testing.T) {
	type row struct {
		i Index
		d Layer
	}
	tests := [...]row{
		{0, 0},
		{1, 1},
		{2, 0},
		{3, 2},
		{4, 0},
		{5, 1},
		{6, 0},
		{7, 3},
	}
	for _, test := range tests {
		if actual := test.i.SubtreeDepth(); actual != test.d {
			t.Errorf("expected %d.SubtreeDepth() to return %d, got %d", test.i, test.d, actual)
		}
	}
}

func TestIndexTreeDepth(t *testing.T) {
	type row struct {
		i Index
		d Layer
	}
	tests := [...]row{
		{0, 0},
		{1, 1},
		{2, 2},
		{3, 2},
		{4, 3},
		{5, 3},
		{6, 3},
		{7, 3},
	}
	for _, test := range tests {
		if actual := test.i.TreeDepth(); actual != test.d {
			t.Errorf("expected %d.TreeDepth() to return %d, got %d", test.i, test.d, actual)
		}
	}
	if d := MaxIndex.TreeDepth(); d != MaxLayer {
		t.Errorf("expected MaxIndex.TreeDepth() to return MaxLayer (%d), got %d", MaxLayer, d)
	}
}
