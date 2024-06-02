// This code is part of the filters package which handles one of the three main components of the GSP framework: filters.
// Filters are functions that can be applied to a signal on a graph.
// The filters package contains a FilterFunc type which is a function that can process a signal, given a graph and filter coefficients.
// Specific filters built on top of this FilterFunc type are also defined in this package.

package filters

import (
	"errors"
	"example/gogsp/graphs"
	"example/gogsp/signals"
)

type FilterFunc func(graph *graphs.Graph, coefficients []float64, signal signals.Signal) (signals.Signal, error)

func ApplyFilter(filter FilterFunc, graph *graphs.Graph, signal signals.Signal) (signals.Signal, error) {
	graph.UpdateAdjacencyMatrix()
	coefficients := make([]float64, len(signal))

	for node := range signal {
		// Use degree as coefficient - nodes with higher degree are considered more important
		coefficients[node] = float64(len(graph.AdjacencyList[graphs.Node(node)]))
	}

	filteredSignal, err := filter(graph, coefficients, signal)
	if err != nil {
		return nil, err
	}
	return filteredSignal, nil
}

var LaplacianFilter FilterFunc = func(graph *graphs.Graph, coefficients []float64, signal signals.Signal) (signals.Signal, error) {
	if len(coefficients) != len(signal) {
		return nil, errors.New("mismatch in size between coefficients and signal")
	}

	output := make(signals.Signal, len(signal))
	for node := range signal {
		total := 0.0
		edges := graph.AdjacencyList[graphs.Node(node)]
		for _, edge := range edges {
			total += float64(edge.Weight) * (signal[node] - signal[edge.Node])
		}
		// multiply with calculated coefficients
		output[node] = coefficients[node] * total
	}

	return output, nil
}

var HighPassFilter FilterFunc = func(graph *graphs.Graph, coefficients []float64, signal signals.Signal) (signals.Signal, error) {
	// same code as before, but coefficient calculation is removed 
	if len(coefficients) != len(signal) {
		return nil, errors.New("mismatch in size between coefficients and signal")
	}

	output := make(signals.Signal, len(signal))
	for node := range signal {
		output[node] = coefficients[node] * signal[node]
	}
	return output, nil
}


var FourierFilter FilterFunc = func(graph *graphs.Graph, coefficients []float64, signal signals.Signal) (signals.Signal, error) {
	if len(coefficients) != len(signal) {
		return nil, errors.New("mismatch in size between coefficients and signal")
	}

	// Compute graph Fourier transform of the signal
	ftSignal := graph.FourierTransform(signal)

	output := make(signals.Signal, len(signal))

	for node := range signal {
		// The Fourier transform already provides the representation in the frequency domain.
		// We just need to multiply with the coefficients to complete the filtering operation.
		output[node] = coefficients[node] * ftSignal[node]
	}

	// Apply inverse Fourier transform to get back to the spatial domain
	filteredSignal := graph.InverseFourierTransform(output)

	return filteredSignal, nil
}
