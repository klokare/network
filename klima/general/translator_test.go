package general

import (
	"context"
	"testing"

	"github.com/klokare/evo"
	"github.com/klokare/network"
)

const name = "klima-general"

var expectedLayered = &Network{
	NP: []int{0, 0, 0, 3, 6, 9},
	NC: []int{0, 0, 3, 3, 3},
	W: []float64{
		-7.900839, 7.417331, 5.509324,
		-7.949483, 7.118832, -3.439597,
		-3.771610, 2.606777, 3.187178,
	},
	M:  []int{0, 1, 2, 0, 1, 3, 2, 3, 4},
	F:  []evo.Activation{evo.Direct, evo.Direct, evo.SteepenedSigmoid, evo.SteepenedSigmoid, evo.SteepenedSigmoid},
	N:  5,
	NI: 2,
	NH: 2,
	NO: 1,
}

/*
{Position: evo.Position{Layer: 0.5, X: 0.5}, Neuron: evo.Hidden, Activation: evo.SteepenedSigmoid, Bias: -1.695151},
{Position: evo.Position{Layer: 1.0, X: 0.5}, Neuron: evo.Output, Activation: evo.SteepenedSigmoid, Bias: -1.967445},
*/
var expectedForward = &Network{
	NP: []int{0, 0, 0, 3, 7},
	NC: []int{0, 0, 3, 4},
	W: []float64{
		3.650676, -4.790058, -1.695151,
		-4.028692, 3.972927, 7.995010, -1.967445,
	},
	M:  []int{0, 1, 2, 0, 1, 2, 3},
	F:  []evo.Activation{evo.Direct, evo.Direct, evo.SteepenedSigmoid, evo.SteepenedSigmoid},
	N:  4,
	NI: 2,
	NH: 1,
	NO: 1,
}

func TestTranslateForward(t *testing.T) {
	testTranslate(t, network.ForwardXor, expectedForward)
}

func TestTranslateLayered(t *testing.T) {
	testTranslate(t, network.LayeredXor, expectedLayered)
}

func testTranslate(t *testing.T, sub evo.Substrate, exp *Network) {

	// Create the network
	net, err := new(Translator).Translate(context.Background(), sub)
	if err != nil {
		t.Errorf("there should be no error: expected nil, actual %v", err)
	}
	if net == nil {
		t.Errorf("there should be a network. there wasn't")
	}
	act := net.(*Network)

	// There should be the correct activation functions
	if len(exp.F) != len(act.F) {
		t.Errorf("incorrect number of activation functions: expected %d, actual %d", len(exp.F), len(act.F))
	} else {
		for i := 0; i < len(exp.F); i++ {
			if exp.F[i] != act.F[i] {
				t.Errorf("incorrect activation function at %d: expected %v, actual %v", i, exp.F[i], act.F[i])
			}
		}
	}

	// The number of incoming connections + bias should be correct
	if len(exp.NC) != len(act.NC) {
		t.Errorf("incorrect number of incoming connections: expected %d, actual %d", len(exp.NC), len(act.NC))
	} else {
		for i := 0; i < len(exp.NC); i++ {
			if exp.NC[i] != act.NC[i] {
				t.Errorf("incorrect number of incoming connectinos for neuron %d: expected %v, actual %v", i, exp.NC[i], act.NC[i])
			}
		}
	}

	// The neuron pointer for each layer should be correct
	if len(exp.NP) != len(act.NP) {
		t.Errorf("incorrect number of neuron pointer records: expected %d, actual %d", len(exp.NP), len(act.NP))
	} else {
		for i := 0; i < len(exp.NP); i++ {
			if exp.NP[i] != act.NP[i] {
				t.Errorf("incorrect neuron pointer at neuron %d: expected %v, actual %v", i, exp.NP[i], act.NP[i])
			}
		}
	}

	// The node for weight index for each weight should be correct
	if len(exp.M) != len(act.M) {
		t.Errorf("incorrect number of node index for weight records: expected %d, actual %d", len(exp.M), len(act.M))
	} else {
		for i := 0; i < len(exp.M); i++ {
			if exp.M[i] != act.M[i] {
				t.Errorf("incorrect node index for weight%d: expected %v, actual %v", i, exp.M[i], act.M[i])
			}
		}
	}

	// The weight array should be correct
	if len(exp.W) != len(act.W) {
		t.Errorf("incorrect number of weights: expected %d, actual %d", len(exp.W), len(act.W))
	} else {
		for i := 0; i < len(exp.W); i++ {
			if exp.W[i] != act.W[i] {
				t.Errorf("incorrect weight at %d: expected %v, actual %v", i, exp.W[i], act.W[i])
			}
		}
	}
}

func BenchmarkTranslateLayered(b *testing.B) {
	cases := network.GenerateCases(network.GenerateLayered)
	tr := network.Translator{Name: name, Translator: Translator{}}
	b.Run(name, network.TranslateBench(cases, tr))
}

func BenchmarkTranslateForward(b *testing.B) {
	cases := network.GenerateCases(network.GenerateForward)
	tr := network.Translator{Name: name, Translator: Translator{}}
	b.Run(name, network.TranslateBench(cases, tr))
}
