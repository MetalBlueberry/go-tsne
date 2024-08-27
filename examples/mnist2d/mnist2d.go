// Copyright (c) 2018 Daniel Augusto Rizzi Salvadori. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"strconv"

	grob "github.com/MetalBlueberry/go-plotly/generated/v2.31.1/graph_objects"
	"github.com/MetalBlueberry/go-plotly/pkg/offline"
	"github.com/MetalBlueberry/go-plotly/pkg/types"
	"github.com/danaugrs/go-tsne/tsne"
	"github.com/sjwhitworth/golearn/pca"
	"gonum.org/v1/gonum/mat"
)

// https://plotly.com/javascript/gapminder-example/

func main() {

	// Parameters
	pcaComponents := 50
	perplexity := float64(300)
	learningRate := float64(300)
	fmt.Printf("PCA Components = %v\nPerplexity = %v\nLearning Rate = %v\n\n", pcaComponents, perplexity, learningRate)

	// Load a subset of MNIST with 2500 records
	X, Y := LoadMNIST()

	// Pre-process the data with PCA (Principal Component Analysis)
	// reducing the number of dimensions from 784 (28x28) to the top pcaComponents principal components
	Xdense := mat.DenseCopyOf(X)
	pcaTransform := pca.NewPCA(pcaComponents)
	Xt := pcaTransform.FitTransform(Xdense)

	// Create output directory if not exists
	os.Mkdir("output", 0770)

	// Create the t-SNE dimensionality reductor and embed the MNIST data in 2D
	frames := []grob.Frame{}
	t := tsne.NewTSNE(2, perplexity, learningRate, 300, true)
	t.EmbedData(Xt, func(iter int, divergence float64, embedding mat.Matrix) bool {
		if iter%10 == 0 {
			fmt.Printf("Iteration %d: divergence is %v\n", iter, divergence)
			frames = append(frames, buildFrame(t.Y, Y, strconv.Itoa(iter)))
		}
		return false
	})

	fig := plotly2D(t.Y, Y, frames)
	offline.Serve(fig)
}

type Serie struct {
	X, Y []float64
}

func getSeries(Y, labels mat.Matrix) []*Serie {
	classes := []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	data := make([]*Serie, len(classes))
	for i := range classes {
		data[i] = &Serie{
			X: []float64{},
			Y: []float64{},
		}
	}

	n, _ := Y.Dims()
	for i := 0; i < n; i++ {
		label := int(labels.At(i, 0))
		data[label].X = append(data[label].X, Y.At(i, 0))
		data[label].Y = append(data[label].Y, Y.At(i, 1))
	}
	return data
}

func buildFrame(Y, labels mat.Matrix, name string) grob.Frame {
	series := getSeries(Y, labels)
	data := make([]types.Trace, len(series))
	for i, serie := range series {
		data[i] = &grob.Scatter{
			Name: types.S(strconv.Itoa(i)),
			X:    types.DataArray(serie.X),
			Y:    types.DataArray(serie.Y),
		}
	}

	f := grob.Frame{
		Data: data,
		Name: types.S(name),
	}
	return f

}

func plotly2D(Y, labels mat.Matrix, frames []grob.Frame) *grob.Fig {
	series := getSeries(Y, labels)
	data := make([]types.Trace, len(series))
	for i, serie := range series {
		data[i] = &grob.Scatter{
			Name: types.S(strconv.Itoa(i)),
			X:    types.DataArray(serie.X),
			Y:    types.DataArray(serie.Y),
			Mode: grob.ScatterModeMarkers,
		}
	}

	fig := &grob.Fig{
		Data: data,
		Layout: &grob.Layout{
			Height: types.N(500),
			Title:  &grob.LayoutTitle{Text: "t-SNE MNIST"},
			Xaxis: &grob.LayoutXaxis{
				Title: &grob.LayoutXaxisTitle{
					Text: "X",
				},
			},
			Yaxis: &grob.LayoutYaxis{
				Title: &grob.LayoutYaxisTitle{
					Text: "Y",
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
									Mode:        "immediate",
									FromCurrent: true,
									Frame: map[string]interface{}{
										"duration": 200,
										"redraw":   false,
									},
									Transition: map[string]interface{}{
										"duration": 200,
									},
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
				Duration: types.N(100),
				Redraw:   types.True,
			},
			Transition: &grob.AnimationTransition{
				Duration: types.N(50),
				Easing:   grob.AnimationTransitionEasingLinear,
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
