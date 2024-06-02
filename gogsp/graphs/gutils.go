// gutils.go contains the utility functions for the graphs package.
// gutils.go contains the following:
// 	- PrintGraph
// 	- PrintAdjacencyMatrix
// 	- UpdateWeightedGraph
//  - UpdateAdjacencyMatrix
// 	type UnionFind and function NewUnionFind

package graphs

import (
	"example/gogsp/signals"
	"fmt"

	"gonum.org/v1/gonum/mat"
)

// UpdateWeightedGraph updates the weighted graph matrix based on the current graph's adjacency list.
func (g *Graph) UpdateWeightedGraph() {
	size := len(g.AdjacencyList)
	g.WeightedGraph = make([][]Weight, size)
	for i := 0; i < size; i++ {
		g.WeightedGraph[i] = make([]Weight, size)
		for j := 0; j < size; j++ {
			weight := Weight(0)
			for _, edge := range g.AdjacencyList[Node(i)] {
				if edge.Node == Node(j) {
					weight = edge.Weight
					break
				}
			}
			g.WeightedGraph[i][j] = weight
		}
	}
}

// UpdateLaplacianMatrix updates the Laplacian matrix based on the current graph's weighted graph.
func (g *Graph) UpdateLaplacianMatrix() {
	size := len(g.WeightedGraph)
	g.LaplacianMatrix = make([][]Weight, size)
	for i := 0; i < size; i++ {
		g.LaplacianMatrix[i] = make([]Weight, size)
		degree := Weight(0)
		for j := 0; j < size; j++ {
			if i != j {
				degree += g.WeightedGraph[i][j]
				g.LaplacianMatrix[i][j] = -g.WeightedGraph[i][j]
			}
		}
		g.LaplacianMatrix[i][i] = degree
	}
}

// UpdateAdjacencyMatrix updates the adjacency matrix based on the current graph's adjacency list.
func (g *Graph) UpdateAdjacencyMatrix() {
	size := len(g.AdjacencyList)
	g.AdjacencyMatrix = make([][]Weight, size)
	for i := 0; i < size; i++ {
		g.AdjacencyMatrix[i] = make([]Weight, size)
		for j := 0; j < size; j++ {
			if i != j {
				for _, edge := range g.AdjacencyList[Node(i)] {
					if edge.Node == Node(j) {
						g.AdjacencyMatrix[i][j] = edge.Weight
						break
					}
				}
			}
		}
	}
}

func (g *Graph) IsFullyConnected() bool {
	// Get the number of nodes in the graph
	numNodes := len(g.AdjacencyList)
	if numNodes == 0 {
		// Empty graph, consider it fully connected
		return true
	}

	// Set to keep track of visited nodes
	visited := make(map[Node]bool)

	// Perform DFS traversal starting from an arbitrary node
	startNode := Node(0)
	g.dfs(startNode, visited)

	// Check if all nodes have been visited
	return len(visited) == numNodes
}

// Depth-First Search traversal
func (g *Graph) dfs(node Node, visited map[Node]bool) {
	visited[node] = true

	// Visit all adjacent nodes recursively
	for _, edge := range g.AdjacencyList[node] {
		if !visited[edge.Node] {
			g.dfs(edge.Node, visited)
		}
	}
}

func (g *Graph) LaplacianToMatDense() *mat.Dense {
	r, c := len(g.LaplacianMatrix), len(g.LaplacianMatrix[0])
	data := make([]float64, r*c)
	for i := 0; i < r; i++ {
		float64Row := convertWeightToFloat64(g.LaplacianMatrix[i])
		copy(data[i*c:(i+1)*c], float64Row)
	}
	return mat.NewDense(r, c, data)
}

func convertWeightToFloat64(weights []Weight) []float64 {
	float64s := make([]float64, len(weights))
	for i, weight := range weights {
		float64s[i] = float64(weight)
	}
	return float64s
}

func (g *Graph) LaplacianToMatSymDense() *mat.SymDense {
	// assuming Laplacian matrix is symmetric
	r := len(g.LaplacianMatrix)
	data := make([]float64, r*r)
	for i := 0; i < r; i++ {
		float64Row := convertWeightToFloat64(g.LaplacianMatrix[i])
		copy(data[i*r:(i+1)*r], float64Row)
	}
	return mat.NewSymDense(r, data)
}

// PrintWeightedGraph prints the weighted graph matrix.
func (g *Graph) PrintWeightedGraph() {
	fmt.Println("Weighted Graph:")
	for i, row := range g.WeightedGraph {
		fmt.Printf("Node %v :", Node(i))
		for _, weight := range row {
			// Format the weights to two decimal places
			fmt.Printf(" %.2f", weight)
		}
		fmt.Println()
	}
}

type UnionFind struct {
	parent []Node
	rank   []int
}

func NewUnionFind(size int) *UnionFind {
	uf := &UnionFind{
		parent: make([]Node, size),
		rank:   make([]int, size),
	}
	for i := range uf.parent {
		uf.parent[i] = Node(i)
		uf.rank[i] = 0
	}
	return uf
}

func (uf *UnionFind) Find(node Node) Node {
	if uf.parent[node] != node {
		uf.parent[node] = uf.Find(uf.parent[node])
	}
	return uf.parent[node]
}

func (uf *UnionFind) Union(node1, node2 Node) {
	root1 := uf.Find(node1)
	root2 := uf.Find(node2)
	if root1 != root2 {
		if uf.rank[root1] < uf.rank[root2] {
			uf.parent[root1] = root2
		} else if uf.rank[root1] > uf.rank[root2] {
			uf.parent[root2] = root1
		} else {
			uf.parent[root2] = root1
			uf.rank[root1]++
		}
	}
}
func (g *Graph) GraphFourierTransform(s signals.Signal) (signals.Signal, error) {
	// Convert g.LaplacianMatrix to *mat.SymDense
	l := g.LaplacianToMatSymDense()

	// Prepare EigenSym and compute the eigendecomposition of the Laplacian
	var es mat.EigenSym
	ok := es.Factorize(l, true)
	if !ok {
		return nil, fmt.Errorf("failed to factorize Laplacian matrix")
	}

	// Get eigen vectors
	var eigenVectors mat.Dense
	es.VectorsTo(&eigenVectors)

	// Convert signal to mat.VecDense
	sVec := mat.NewVecDense(len(s), s)

	// Prepare the transformed signal
	ts := make(signals.Signal, len(s))

	// Transformed signal as a vector
	var tsVec mat.VecDense

	tsVec.MulVec(eigenVectors.T(), sVec)

	// Copy the results back to a Signal
	for i := range ts {
		ts[i] = tsVec.AtVec(i)
	}

	return ts, nil
}

func (g *Graph) InverseGraphFourierTransform(s signals.Signal) (signals.Signal, error) {
	// Convert g.LaplacianMatrix to *mat.SymDense
	l := g.LaplacianToMatSymDense()

	// Use EigenSym to compute the eigendecomposition of the Laplacian
	var es mat.EigenSym
	ok := es.Factorize(l, true)
	if !ok {
		return nil, fmt.Errorf("failed to factorize Laplacian matrix")
	}

	// Get eigen vectors
	var eigenVectors mat.Dense
	es.VectorsTo(&eigenVectors)

	// Convert signal to mat.VecDense
	sVec := mat.NewVecDense(len(s), s)

	// Prepare the inverse transformed signal
	its := make(signals.Signal, len(s))

	// Inverse transformed signal as a vector
	var itsVec mat.VecDense

	// Now, you use eigenVectors (not Transposed).
	itsVec.MulVec(&eigenVectors, sVec)

	// Copy the results back to a Signal
	for i := range its {
		its[i] = itsVec.AtVec(i)
	}

	return its, nil
}
