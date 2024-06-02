# GoGSP - Graph Signal Processing with Golang

This library aims to bring the power of signal processing to the context of graphs. It allows the creation and manipulation of Graphs, Signals and Filters and applies the principles of signal processing in the graph domain.

## Features

- Defines the Graph structure with adjacency list, adjacency matrix, and Laplacian matrix.
- Supports addition and deletion of node, additions of weighted edges.
- Defines a standard `Signal` type to operate on Graphs.
- Provides out-of-the-box filters such as Laplacian filter, high-pass filter, low-pass filter.
- Ability to perform a Graph Fourier Transform operation.

## Usage 

First declare a new graph:

g := graphs.NewGraph()
```

Adding nodes and edge:

``` go
g.AddNode(1)
g.AddNode(2)
g.AddEdge(1, 2, 3.5)
```

Creating a signal:

``` go
signal := signals.CreateSignal(2)
signal.SetSignal(g, []float64{3, 7})
```

Applying filters:

``` go
filteredSignal, err := filters.ApplyFilter(filters.LaplacianFilter, g, signal)
```

A detailed example can be found in `main.go`.


## Disclaimer 

Please note that this library is a simplified implementation of graph signal processing, and as such, it does not support some features you may find in more comprehensive libraries.
