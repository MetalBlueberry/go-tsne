module github.com/danaugrs/go-tsne/examples/mnist3d

go 1.22

toolchain go1.22.3

// To run with the local go-tsne uncomment the following line
// replace github.com/danaugrs/go-tsne/tsne => ../../tsne

require (
	github.com/danaugrs/go-tsne/tsne v0.0.0-20220306153449-0ee45704632c
	github.com/sjwhitworth/golearn v0.0.0-20211014193759-a8b69c276cd8
	gonum.org/v1/gonum v0.9.3
)

require (
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	golang.org/x/sys v0.24.0 // indirect
)

require github.com/MetalBlueberry/go-plotly v0.6.2-0.20240827072729-c72e96e3021a
