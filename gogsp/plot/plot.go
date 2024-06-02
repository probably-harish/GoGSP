package plot

import (
	"example/gogsp/graphs"
	"example/gogsp/signals"
	"fmt"
	"math"
	"os"

	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

func PlotSignal(signal signals.Signal, name string) {
	// Prepare the data for plotting
	xValues := make([]float64, len(signal))
	yOriginal := make([]float64, len(signal))

	for i := range signal {
		xValues[i] = float64(i)
		yOriginal[i] = signal[i]
	}

	// Create a new chart
	graphChart := chart.Chart{
		XAxis: chart.XAxis{
			Name: "Node",
		},
		YAxis: chart.YAxis{
			Name: "Signal Value",
		},
		Series: []chart.Series{
			&chart.ContinuousSeries{
				Name:    "Original Signal",
				XValues: xValues,
				YValues: yOriginal,
				Style: chart.Style{
					Show:        true,
					StrokeColor: drawing.ColorBlue,
					StrokeWidth: 3,
				},
			},
			&chart.ContinuousSeries{
				Name:    "Filtered Signal",
				XValues: xValues,
				Style: chart.Style{
					Show:        true,
					StrokeColor: drawing.ColorRed,
					StrokeWidth: 3,
				},
			},
		},
	}

	// Save the chart to a file
	file_name := fmt.Sprintf("%s.png", name)
	file, err := os.Create(file_name)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	err = graphChart.Render(chart.PNG, file)
	if err != nil {
		fmt.Println("Error rendering chart:", err)
		return
	}

	fmt.Println("Chart saved to graph_signal.png")
}

func applyFunctionToFloat64Array(arr []float64, fn func(float64) float64) []float64 {
	result := make([]float64, len(arr))
	for i, v := range arr {
		result[i] = fn(v)
	}
	return result
}
func sqrt(x float64) float64 {
	return math.Sqrt(Abs(1 - x*x))
}

func Abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func PlotGraph(graph *graphs.Graph, name string) {
	// Prepare the data for plotting
	nodes := make([]float64, 0)
	edgesX := make([]float64, 0)
	edgesY := make([]float64, 0)

	for node, edges := range graph.AdjacencyList {
		nodes = append(nodes, float64(node))
		for _, edge := range edges {
			edgesX = append(edgesX, float64(node))
			edgesX = append(edgesX, float64(edge.Node))
			edgesY = append(edgesY, float64(node))
			edgesY = append(edgesY, float64(edge.Node))
		}
	}
	fmt.Print(nodes, "\n")

	viridisByY := func(xr, yr chart.Range, index int, x, y float64) drawing.Color {
		return chart.Viridis(y, yr.GetMin(), yr.GetMax())
	}

	// Duplicate the first node at the end to form a ring/oval shape
	nodes = []float64{(1 / math.Sqrt(2)), (1 / math.Sqrt(2)), (-1 / math.Sqrt(2)), (-1 / math.Sqrt(2))}
	resultArray := []float64{(1 / math.Sqrt(2)), (-1 / math.Sqrt(2)), (-1 / math.Sqrt(2)), (1 / math.Sqrt(2))}

	// Create a new chart
	graphChart := chart.Chart{
		Series: []chart.Series{
			&chart.ContinuousSeries{
				Style: chart.Style{
					Show:             true,
					StrokeWidth:      chart.Disabled,
					DotWidth:         5,
					DotColorProvider: viridisByY,
				},
				XValues: nodes,
				YValues: resultArray,
			},
		},
	}

	// Save the chart to a file
	fileName := fmt.Sprintf("%s.png", name)
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	err = graphChart.Render(chart.PNG, file)
	if err != nil {
		fmt.Println("Error rendering chart:", err)
		return
	}

	fmt.Println("Chart saved to", fileName)
}
