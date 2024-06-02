// futils.go contains utility functions for the filters package.
// futils contains the following functions:
// 	- computeFourierTransform
// 	- computeInverseFourierTransform
// 	- diagonalizeMatrix

package filters

import (
	"example/gogsp/signals"
	"fmt"
	"math"

	"github.com/mjibson/go-dsp/fft"
	"gonum.org/v1/gonum/mat"
)

func computeFourierTransform(signal signals.Signal) signals.Signal {
	// Create a complex slice for the input to the Fourier transform
	complexSignal := make([]complex128, len(signal))
	for i, val := range signal {
		complexSignal[i] = complex(val, 0) // Assume real-valued signal
	}

	// Perform the Fourier transform
	fftResult := fft.FFT(complexSignal)

	// Compute the magnitude of the transformed signal
	transformedSignal := make(signals.Signal, len(signal))
	for i, val := range fftResult {
		transformedSignal[i] = math.Abs(real(val))
	}

	return transformedSignal
}

func computeInverseFourierTransform(transformedSignal signals.Signal) signals.Signal {
	// Create a complex slice for the input to the inverse Fourier transform
	complexTransformedSignal := make([]complex128, len(transformedSignal))
	for i, val := range transformedSignal {
		complexTransformedSignal[i] = complex(val, 0)
	}

	// Perform the inverse Fourier transform
	ifftResult := fft.IFFT(complexTransformedSignal)

	// Extract the real part of the inverse transformed signal
	inverseTransformedSignal := make(signals.Signal, len(transformedSignal))
	for i, val := range ifftResult {
		inverseTransformedSignal[i] = real(val)
	}

	return inverseTransformedSignal
}

func DiagonalizeMatrix(matrix [][]float64) ([]float64, [][]float64, error) {
	// Convert input matrix to gonum matrix
	m := mat.NewSymDense(len(matrix), nil)
	for i, row := range matrix {
		for j, val := range row {
			m.SetSym(i, j, val)
		}
	}

	// Perform diagonalization
	var es mat.EigenSym
	ok := es.Factorize(m, true)
	if !ok {
		return nil, nil, fmt.Errorf("diagonalization failed")
	}

	// Retrieve eigenvalues
	eigenvalues := es.Values(nil)

	// Prepare a Dense to receive the eigenvectors
	var eigenVectors mat.Dense
	es.VectorsTo(&eigenVectors)

	// Convert eigenvectors to a 2D slice
	eigenvectorsSlice := make([][]float64, len(matrix))
	for i := 0; i < len(matrix); i++ {
		eigenvectorsSlice[i] = make([]float64, len(matrix))
		for j := 0; j < len(matrix); j++ {
			eigenvectorsSlice[i][j] = eigenVectors.At(j, i)
		}
	}

	return eigenvalues, eigenvectorsSlice, nil
}
