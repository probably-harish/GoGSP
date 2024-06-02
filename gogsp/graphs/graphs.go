// Author: Harish Raju
// github: github.com/probably-harish

// Package graphs implements graph data structures and algorithms to be used in the context of graph signal processing
// The package is inspired by the Python package PyGSP

package graphs

import (
	"fmt"
	"math/rand"
	"strconv"
)

type Node int

type Weight float64

type Edge struct {
	Node   Node
	Weight Weight
}

type Graph struct {
	AdjacencyList   map[Node][]Edge
	WeightedGraph   [][]Weight
	LaplacianMatrix [][]Weight
	AdjacencyMatrix [][]Weight
}

func NewGraph() *Graph {
	return &Graph{
		AdjacencyList:   make(map[Node][]Edge),
		WeightedGraph:   nil,
		LaplacianMatrix: nil,
		AdjacencyMatrix: nil,
	}
}

func (g *Graph) AddNode(n Node) {
	if _, present := g.AdjacencyList[n]; !present {
		g.AdjacencyList[n] = []Edge{}
	}
}

func (g *Graph) AddEdge(n1, n2 Node, weight Weight) {
	g.AddNode(n1)
	g.AddNode(n2)
	g.AdjacencyList[n1] = append(g.AdjacencyList[n1], Edge{n2, weight})
}

func RandomWeightedGraph(size int) *Graph {
	g := NewGraph()

	type Edge struct {
		Node1  Node
		Node2  Node
		Weight Weight
	}

	// Generate all possible edges with random weights
	edges := make([]Edge, 0, size*(size-1)/2)
	for i := 0; i < size; i++ {
		for j := i + 1; j < size; j++ {
			weight := Weight(rand.Float64())
			edges = append(edges, Edge{
				Node1:  Node(i),
				Node2:  Node(j),
				Weight: weight,
			})
		}
	}

	// Shuffle the edges randomly
	rand.Shuffle(len(edges), func(i, j int) {
		edges[i], edges[j] = edges[j], edges[i]
	})

	uf := NewUnionFind(size)
	for _, edge := range edges {
		node1, node2 := edge.Node1, edge.Node2
		if uf.Find(node1) != uf.Find(node2) {
			uf.Union(node1, node2)
			g.AddEdge(node1, node2, edge.Weight)
			g.AddEdge(node2, node1, edge.Weight)
		}
	}

	return g
}

func (g *Graph) PrintGraph() {
	for node, edges := range g.AdjacencyList {
		fmt.Printf("\nNode %v:", node)
		for _, edge := range edges {
			weight := strconv.FormatFloat(float64(edge.Weight), 'f', 2, 64)
			fmt.Printf("\n%v   Weight: %v", edge.Node, weight)
		}
		fmt.Printf("\n")
	}
}

// PrintAdjacencyMatrix prints the adjacency matrix.
func (g *Graph) PrintAdjacencyMatrix() {
	fmt.Println("Adjacency Matrix:")
	for i, row := range g.AdjacencyMatrix {
		fmt.Printf("Node %v :", Node(i))
		for _, weight := range row {
			// Format the weights to two decimal places
			fmt.Printf(" %.2f", weight)
		}
		fmt.Println()
	}
}

// PrintLaplacianMatrix prints the Laplacian matrix.
func (g *Graph) PrintLaplacianMatrix() {
	fmt.Println("Laplacian Matrix:")
	for i, row := range g.LaplacianMatrix {
		fmt.Printf("Node %v :", Node(i))
		for _, weight := range row {
			// Format the weights to two decimal places
			fmt.Printf(" %.2f", weight)
		}
		fmt.Println()
	}
}
