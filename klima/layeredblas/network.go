package layeredblas

import (
	"github.com/klokare/evo"
	"gonum.org/v1/gonum/blas/blas64"
)

type Network struct {
	N  []int            // number of neurons in each layer
	np []int            // neuron pointer for layer
	wp []int            // weight pointer for layer
	w  []float64        // weights
	f  []evo.Activation // activation functions
}

func (net *Network) Activate(inputs []float64) (outputs []float64, err error) {

	// Create the neuron states array
	n := make([]float64, len(net.f))
	copy(n, inputs)

	// Main loop
	ni := net.np[1]
	wi := 0
	for l := 1; l < len(net.N); l++ {
		for i := 0; i < net.N[l]; i++ {
			signal := net.w[wi] // Bias
			wi++
			npi := net.np[l-1] // Index of neuron in the previous layer
			cnt := net.N[l-1]
			signal += blas64.Dot(cnt, blas64.Vector{Inc: 1, Data: net.w[wi : wi+cnt]}, blas64.Vector{Inc: 1, Data: n[npi : npi+cnt]})
			// for j := 0; j < net.N[l-1]; j++ {
			// 	signal += net.w[wi] * n[npi]
			// 	wi++
			// 	npi++
			// }
			wi += cnt
			npi += cnt
			n[ni] = net.f[ni].Activate(signal)
			ni++
		}
	}

	// Copy the outputs
	outputs = make([]float64, net.N[len(net.N)-1])
	copy(outputs, n[len(n)-len(outputs):])
	return
}
