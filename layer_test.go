package historytree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvalidLayer(t *testing.T) {
	assert.Panics(t, func() { (MaxLayer + 1).assertValid() })
}

func TestLayerMaxIndex(t *testing.T) {
	type row struct {
		r Layer
		i Index
	}
	tests := [...]row{
		{0, 0},
		{1, 1},
		{2, 3},
		{3, 7},
	}
	for _, test := range tests {
		if actual := test.r.MaxIndex(); actual != test.i {
			t.Errorf("expected %d.MaxIndex() to return %d, got %d", test.r, test.i, actual)
		}
	}
	for r := Layer(4); r < MaxLayer; r++ {
		if actual := (r + 1).MaxIndex(); actual != Index(r.SubtreeCount()) {
			t.Errorf("expected %d.MaxIndex() to return %d, got %d", r+1, r.SubtreeCount(), actual)
		}
	}
	if i := MaxLayer.MaxIndex(); i != MaxIndex {
		t.Errorf("expected MaxLayer.MaxIndex() to return MaxIndex (%d), got %d", MaxIndex, i)
	}
}

func TestLayerSubtreeCount(t *testing.T) {
	type row struct {
		r Layer
		n uint64
	}
	tests := [...]row{
		{0, 1},
		{1, 3},
		{2, 7},
		{3, 15},
	}
	for _, test := range tests {
		if actual := test.r.SubtreeCount(); actual != test.n {
			t.Errorf("expected %d.SubtreeCount() to return %d, got %d", test.r, test.n, actual)
		}
	}
}
