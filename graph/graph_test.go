package graph

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGraph_Resolve(t *testing.T) {
	nodeA := NewNode("a", "b", "c")
	nodeB := NewNode("b", "b_a", "b_b")
	nodeB_A := NewNode("b_a")
	nodeB_B := NewNode("b_b", "c")
	nodeC := NewNode("c", "c_a")
	nodeC_A := NewNode("c_a", "c_b")
	nodeC_B := NewNode("c_b")

	nodes := []*Node{
		nodeA,
		nodeB,
		nodeB_A,
		nodeB_B,
		nodeC,
		nodeC_A,
		nodeC_B,
	}

	_ = []*Node{
		nodeC_B,
		nodeB_A,
		nodeC_A,
		nodeC,
		nodeB_B,
		nodeB,
		nodeA,
	}

	sut := New(nodes...)

	resolved, err := sut.Resolve()
	assert.NoError(t, err)
	for i, node := range resolved {
		fmt.Printf("[%d] %s\n", i, node.Name)
	}
	// assert.Equal(t, expected, resolved.Nodes())
}
