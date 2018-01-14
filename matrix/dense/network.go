package dense

import (
	"errors"
	"fmt"

	"github.com/klokare/evo"
	"gonum.org/v1/gonum/mat"
)

// Known errors
var (
	ErrIncorrectNumberInputColumns = errors.New("number of columns in the inputs matrix does not match number of input neurons")
)

type Network struct {
	Layers  [][]evo.Activation
	Weights []*mat.Dense
}

func (net Network) Activate(inputs []float64) (outputs []float64, err error) {

	// Add bias and make the matrix
	inputs = append(inputs, 1.0) // Append the bias value
	imat := mat.NewDense(1, len(inputs), inputs)

	// Active the matrix
	var omat *mat.Dense
	if omat, err = net.ActivateMatrix(imat); err != nil {
		return
	}

	// copy the output layer
	outputs = omat.RawRowView(0)
	return
}

func (net Network) ActivateMatrix(inputs *mat.Dense) (outputs *mat.Dense, err error) {

	// Incorrect matrix size
	var n, m int
	_, m = inputs.Dims()
	if m != len(net.Layers[0]) {
		err = ErrIncorrectNumberInputColumns
		return
	}

	// Execute the network
	var src, tgt *mat.Dense
	src = inputs

	for w := 0; w < len(net.Weights); w++ {

		// Set the bias values
		n, m = src.Dims()
		for i := 0; i < n; i++ {
			src.Set(i, m-1, 1.0)
		}

		// Multiply the matrices
		// TODO: should s be a sparse array?
		s := new(mat.Dense)
		s.Mul(src, net.Weights[w])

		// Apply the activations
		tgt = new(mat.Dense)
		tgt.Apply(func(i int, j int, v float64) float64 {
			return net.Layers[w+1][j].Activate(v)
		}, s)

		// Move to the next layer
		src = tgt
	}

	// Return the output matrix
	outputs = tgt
	return
}

func show(name string, m mat.Matrix) {
	fa := mat.Formatted(m, mat.Prefix("    "), mat.Squeeze())
	fmt.Printf("%s:\na = %v\n\n", name, fa)
}
