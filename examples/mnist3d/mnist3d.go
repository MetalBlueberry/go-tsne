// Copyright (c) 2018 Daniel Augusto Rizzi Salvadori. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strconv"

	grob "github.com/MetalBlueberry/go-plotly/generated/v2.31.1/graph_objects"
	"github.com/MetalBlueberry/go-plotly/pkg/offline"
	"github.com/MetalBlueberry/go-plotly/pkg/types"
	"github.com/danaugrs/go-tsne/tsne"
	"github.com/sjwhitworth/golearn/pca"
	"gonum.org/v1/gonum/mat"
)

// ExampleMNIST3D is this example application displaying t-SNE being performed in real-time, in 3D,
// using a subset of the MNIST dataset of handwritten digits.

func main() {

	pcaComponents := 100
	perplexity := float64(500)
	learningRate := float64(500)
	fmt.Printf("PCA Components = %v\nPerplexity = %v\nLearning Rate = %v\n\n", pcaComponents, perplexity, learningRate)

	// Load a subset of MNIST with 2500 records
	X, Y := LoadMNIST()

	Xdense := mat.DenseCopyOf(X)
	pcaTransform := pca.NewPCA(pcaComponents)
	Xt := pcaTransform.FitTransform(Xdense)

	frames := []grob.Frame{}

	// Create the t-SNE dimensionality reductor and embed the MNIST data in 3D
	t := tsne.NewTSNE(3, perplexity, learningRate, 300, true)
	t.EmbedData(Xt, func(iter int, divergence float64, embedding mat.Matrix) bool {
		if iter%10 == 0 {
			fmt.Printf("Iteration %d: divergence is %v\n", iter, divergence)
			frames = append(frames, buildFrame(t.Y, Y, strconv.Itoa(iter)))
		}
		return false
	})
	fig := plotly3D(t.Y, Y, frames)
	offline.Serve(fig)

}

type Serie struct {
	X, Y, Z []float64
}

func getSeries(Y, labels mat.Matrix) []*Serie {
	classes := []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	data := make([]*Serie, len(classes))
	for i := range classes {
		data[i] = &Serie{
			X: []float64{},
			Y: []float64{},
			Z: []float64{},
		}
	}

	n, _ := Y.Dims()
	for i := 0; i < n; i++ {
		label := int(labels.At(i, 0))
		data[label].X = append(data[label].X, Y.At(i, 0))
		data[label].Y = append(data[label].Y, Y.At(i, 1))
		data[label].Z = append(data[label].Z, Y.At(i, 2))
	}
	return data
}

var colorPalette = []types.Color{
	"rgba(255, 0, 0, 0.8)",
	"rgba(0, 255, 0, 0.8)",
	"rgba(0, 0, 255, 0.8)",
	"rgba(255, 255, 0, 0.8)",
	"rgba(255, 0, 255, 0.8)",
	"rgba(0, 255, 255, 0.8)",
	"rgba(128, 0, 128, 0.8)",
	"rgba(255, 165, 0, 0.8)",
	"rgba(75, 0, 130, 0.8)",
	"rgba(0, 0, 0, 0.8)",
}
var lineColorPalette = []types.Color{
	"rgba(255, 0, 0, 1)",
	"rgba(0, 255, 0, 1)",
	"rgba(0, 0, 255, 1)",
	"rgba(255, 255, 0, 1)",
	"rgba(255, 0, 255, 1)",
	"rgba(0, 255, 255, 1)",
	"rgba(128, 0, 128, 1)",
	"rgba(255, 165, 0, 1)",
	"rgba(75, 0, 130, 1)",
	"rgba(0, 0, 0, 1)",
}

func buildSeries(Y, labels mat.Matrix) []types.Trace {
	series := getSeries(Y, labels)
	data := make([]types.Trace, len(series))
	for i, serie := range series {
		data[i] = &grob.Scatter3d{
			Name: types.S(strconv.Itoa(i)),
			X:    types.DataArray(serie.X),
			Y:    types.DataArray(serie.Y),
			Z:    types.DataArray(serie.Z),
			Mode: grob.Scatter3dModeMarkers,
			Marker: &grob.Scatter3dMarker{
				Color: types.ArrayOKValue(types.UseColor(colorPalette[i])),
				Line: &grob.Scatter3dMarkerLine{
					Color: types.ArrayOKValue(types.UseColor(lineColorPalette[i])),
					Width: types.N(3),
				},
			},
		}
	}
	return data

}

func buildFrame(Y, labels mat.Matrix, name string) grob.Frame {
	data := buildSeries(Y, labels)
	f := grob.Frame{
		Data: data,
		Name: types.S(name),
	}
	return f

}

func plotly3D(Y, labels mat.Matrix, frames []grob.Frame) *grob.Fig {
	data := buildSeries(Y, labels)

	fig := &grob.Fig{
		Data: data,
		Layout: &grob.Layout{
			Height: types.N(800),
			Title:  &grob.LayoutTitle{Text: "t-SNE MNIST3D"},
			Scene: &grob.LayoutScene{
				Xaxis: &grob.LayoutSceneXaxis{
					Title: &grob.LayoutSceneXaxisTitle{
						Text: "X",
					},
					Range: []int{-500, 500},
				},
				Yaxis: &grob.LayoutSceneYaxis{
					Title: &grob.LayoutSceneYaxisTitle{
						Text: "Y",
					},
					Range: []int{-500, 500},
				},
				Zaxis: &grob.LayoutSceneZaxis{
					Title: &grob.LayoutSceneZaxisTitle{
						Text: "Z",
					},
					Range: []int{-500, 500},
				},
			},
			Updatemenus: []grob.LayoutUpdatemenu{
				{
					Type:       "buttons",
					Showactive: types.False,
					Buttons: []grob.LayoutUpdatemenuButton{
						{
							Label:  "Play",
							Method: "animate",
							Args: []*ButtonArgs{
								nil,
								{
									Mode: "immediate",
								},
							},
						},
					},
				},
			},
		},
		Config: &grob.Config{
			Responsive: types.True,
		},
		Frames: frames,
		Animation: &grob.Animation{
			Frame: &grob.AnimationFrame{
				Duration: types.N(0),
				Redraw:   types.False,
			},
			Transition: &grob.AnimationTransition{
				Duration: types.N(0),
			},
		},
	}

	return fig
}

// Define the ButtonArgs type
type ButtonArgs struct {
	Frame       map[string]interface{} `json:"frame,omitempty"`
	Transition  map[string]interface{} `json:"transition,omitempty"`
	FromCurrent bool                   `json:"fromcurrent,omitempty"`
	Mode        string                 `json:"mode,omitempty"`
}
