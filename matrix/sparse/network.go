package sparse

import (
	"errors"
	"fmt"

	"github.com/james-bowman/sparse"
	"github.com/klokare/evo"
	"gonum.org/v1/gonum/mat"
)

// Known errors
var (
	ErrIncorrectNumberInputColumns = errors.New("number of columns in the inputs matrix does not match number of input neurons")
)

type Network struct {
	Layers  [][]evo.Activation
	Weights []*sparse.CSR
}

func (z Network) Activate(inputs []float64) (outputs []float64, err error) {

	// Add bias and make the matrix
	imat := mat.NewDense(1, len(inputs)+1, nil)
	for i, x := range inputs {
		imat.Set(0, i, x)
	}

	// Incorrect matrix size
	n, m := imat.Dims()
	if m != len(z.Layers[0]) {
		err = ErrIncorrectNumberInputColumns
		return
	}

	// Execute the network
	var src, tgt *mat.Dense
	src = imat

	//var tgt *sparse.CSR
	for w := 0; w < len(z.Weights); w++ {

		// Set the bias values
		for i := 0; i < n; i++ {
			src.Set(i, m-1, 1.0)
		}

		// Multiply the matrices
		// TODO: should s be a sparse array?
		s := new(mat.Dense)
		s.Mul(src, z.Weights[w])

		// Apply the activations
		tgt = new(mat.Dense)
		tgt.Apply(func(i int, j int, v float64) float64 {
			return z.Layers[w+1][j].Activate(v)
		}, s)

		// Move to the next layer
		src = tgt
	}

	// copy the output layer
	outputs = tgt.RawRowView(0)
	return
}

func show(name string, m mat.Matrix) {
	fa := mat.Formatted(m, mat.Prefix("    "), mat.Squeeze())
	fmt.Printf("%s:\na = %v\n\n", name, fa)
}
