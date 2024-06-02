package main

import (
	"example/gogsp/filters"
	"example/gogsp/graphs"
	"example/gogsp/signals"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// Setting seed for randomness
	rand.Seed(time.Now().UnixNano())

	// Create a random 7 node graph
	g := graphs.RandomWeightedGraph(7)

	// Print the graph
	g.PrintGraph()
	fmt.Println()

	// Print the weighted graph
	g.UpdateWeightedGraph()
	g.PrintWeightedGraph()
	fmt.Println()

	// Print Laplacian matrix
	g.UpdateLaplacianMatrix()
	g.PrintLaplacianMatrix()
	fmt.Println()

	// Create a signal on the graph
	signal := signals.CreateSignal(7)
	signal.SetSignal(g, []float64{3, 4, 5, 7, 11, 13, 17})

	// Print original signal
	fmt.Println("Original signal:")
	signal.PrintSignal()

	// Compute Graph Fourier Transform
	transformedSignal, err := signal.GraphFourierTransform(g)
	if err != nil {
		fmt.Println(err)
	}

	// Print Fourier Transformed signal
	fmt.Println("\nSignal after Graph Fourier Transform:")
	transformedSignal.PrintSignal()

	// Apply Laplacian filter
	filteredSignal, err := filters.ApplyFilter(filters.LaplacianFilter, g, signal)
	if err != nil {
		fmt.Println(err)
	}

	// Print signal after Laplacian filtering
	fmt.Println("\nSignal after Laplacian filtering:")
	filteredSignal.PrintSignal()

	// Apply HighPass filter
	filteredSignal, err = filters.ApplyFilter(filters.HighPassFilter, g, signal)
	if err != nil {
		fmt.Println(err)
	}

	// Print signal after HighPass filtering
	fmt.Println("\nSignal after HighPass filtering:")
	filteredSignal.PrintSignal()

	// Apply LowPass filter
	filteredSignal, err = filters.ApplyFilter(filters.LowPassFilter, g, signal)
	if err != nil {
		fmt.Println(err)
	}

	// Print signal after LowPass filtering
	fmt.Println("\nSignal after LowPass filtering:")
	filteredSignal.PrintSignal()

	// Apply FourierFilter
	filteredSignal, err = filters.ApplyFilter(filters.FourierFilter, g, signal)
	if err != nil {
		fmt.Println(err)
	}

	// Print signal after Fourier filtering
	fmt.Println("\nSignal after Fourier filtering:")
	filteredSignal.PrintSignal()

}
