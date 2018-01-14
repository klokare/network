package general

import (
	"github.com/klokare/evo"
)

// Network ...
type Network struct {
	NP            []int
	NC            []int
	W             []float64
	M             []int
	F             []evo.Activation // activation functions
	N, NI, NH, NO int
}

// Activate ...
func (net *Network) Activate(inputs []float64) (outputs []float64, err error) {

	// Copy input
	n := make([]float64, net.N)
	copy(n, inputs)

	// Set remaining neuron states to 1
	for ni := net.NI; ni < net.N; ni++ {
		n[ni] = 1.0
	}
	//log.Println("initial state", n)

	// main loop
	var idx int
	for ni := net.NI; ni < net.N; ni++ {
		signal := 0.0
		for wi := net.NP[ni]; wi < net.NP[ni+1]; wi++ {
			idx = net.M[wi]
			signal += net.W[wi] * n[idx]
		}
		//log.Println(ni, "pre-activate ", n, "weights", net.W[net.NP[ni]:net.NP[ni+1]])
		n[ni] = net.F[ni].Activate(signal)
		//log.Println(ni, "post-activate ", n)
	}

	// Return the outputs
	outputs = make([]float64, net.NO)
	copy(outputs, n[len(n)-net.NO:])
	return
}
