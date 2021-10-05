package graph

import (
	"errors"
	"fmt"
	"strings"

	mapset "github.com/deckarep/golang-set"
)

type Graph []*Node

func (g Graph) String() string {
	return displayGraph(g)
}

// Displays the dependency graph
func displayGraph(graph Graph) string {
	sb := &strings.Builder{}
	for _, node := range graph {
		for i, dep := range node.Deps {
			sb.WriteString(fmt.Sprintf("[%d] %s -> %s\n", i, node.Name, dep))
		}
	}
	return sb.String()
}

// Resolve resolves the Graph returning a new Graph with Node in a resolvable order
func (g *Graph) Resolve() (res Graph, err error) {
	return resolveGraph(*g)
}

// Resolves the dependency graph
func resolveGraph(graph Graph) (Graph, error) {
	// A map containing the node names and the actual node object
	nodeNames := make(map[string]*Node)

	// A map containing the nodes and their dependencies
	nodeDependencies := make(map[string]mapset.Set)

	// Populate the maps
	for _, node := range graph {
		nodeNames[node.Name] = node

		dependencySet := mapset.NewSet()
		for _, dep := range node.Deps {
			dependencySet.Add(dep)
		}
		nodeDependencies[node.Name] = dependencySet
	}

	// Iteratively find and remove nodes from the graph which have no dependencies.
	// If at some point there are still nodes in the graph, and we cannot find
	// nodes without dependencies, that means we have a circular dependency
	var resolved Graph
	for len(nodeDependencies) != 0 {
		// Get all nodes from the graph which have no dependencies
		readySet := mapset.NewSet()
		for name, deps := range nodeDependencies {
			if deps.Cardinality() == 0 {
				readySet.Add(name)
			}
		}

		// If there aren't any ready nodes, then we have a circular dependency
		if readySet.Cardinality() == 0 {
			var g Graph
			for name := range nodeDependencies {
				g = append(g, nodeNames[name])
			}

			return g, errors.New("circular dependency found")
		}

		// Remove the ready nodes and add them to the resolved graph
		for name := range readySet.Iter() {
			delete(nodeDependencies, name.(string))
			resolved = append(resolved, nodeNames[name.(string)])
		}

		// Also make sure to remove the ready nodes from the
		// remaining node dependencies as well
		for name, deps := range nodeDependencies {
			diff := deps.Difference(readySet)
			nodeDependencies[name] = diff
		}
	}

	return resolved, nil
}

// Node represents a single node in the graph with it's dependencies
type Node struct {
	// Name of the node
	Name string

	// Value of the node
	value string

	// Dependencies of the node
	Deps []string
}

// New returns a new Graph of the provided nodes
func New(nodes ...*Node) Graph {
	g := Graph{}
	for _, node := range nodes {
		g = append(g, node)
	}
	return g
}

// Nodes returns the Node of the Graph
func (g Graph) Nodes() []*Node {
	return g
}

// NewNode creates a new node
func NewNode(name string, deps ...string) *Node {
	n := &Node{
		Name: name,
		Deps: deps,
	}

	return n
}
