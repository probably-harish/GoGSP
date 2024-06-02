// signals.go contains the Signal type and utility functions for working with signals.

package signals

import (
	"fmt"
	"math"

	"example/gogsp/graphs"

	"gonum.org/v1/gonum/mat"
)

// Signal represents a signal over the nodes of a graph
type Signal []float64

// CreateSignal generates a new signal with the given size
func CreateSignal(size int) Signal {
	return make(Signal, size)
}

// Set sets the value of the signal at a specific node
func (s Signal) Set(n graphs.Node, value float64) {
	s[n] = value
}

func (s Signal) SetSignal(g *graphs.Graph, arr []float64) {
	for i := 0; i < len(arr); i++ {
		s[graphs.Node(i)] = arr[i]
	}
}

// Get returns the value of the signal at a specific node
func (s Signal) Get(n graphs.Node) float64 {
	return s[n]
}

// Mean calculates the mean value of the signal
func (s Signal) Mean() float64 {
	total := 0.0
	for _, value := range s {
		total += value
	}
	return total / float64(len(s))
}

// Normalize normalizes the signal to have a mean of 0 and standard deviation of 1
func (s Signal) Normalize() {
	mean := s.Mean()
	variance := 0.0
	for _, value := range s {
		diff := value - mean
		variance += diff * diff
	}
	variance = variance / float64(len(s))
	stdDev := math.Sqrt(variance)

	for i := range s {
		s[i] = (s[i] - mean) / stdDev
	}
}

// PrintSignal prints the signal in vector format
func (s Signal) PrintSignal() {
	fmt.Printf("%.4f", s)
	fmt.Print("\n")
}

// getMatDense creates mat.Dense from signal. Why? Because gonum/mat doesn't have a constructor for mat.Dense from []float64
func getMatDense(s Signal) *mat.Dense {
	data := make([]float64, len(s))
	copy(data, s)
	v := mat.NewDense(len(data), 1, data)
	return v
}

func (g *graphs.Graph) GraphFourierTransform(s Signal) (Signal, error) {
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
	ts := make(Signal, len(s))

	// Transformed signal as a vector
	var tsVec mat.VecDense

	tsVec.MulVec(eigenVectors.T(), sVec)

	// Copy the results back to a Signal
	for i := range ts {
		ts[i] = tsVec.AtVec(i)
	}

	return ts, nil
}

func (g *graphs.Graph) InverseGraphFourierTransform(s Signal) (Signal, error) {
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
	its := make(Signal, len(s))

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
