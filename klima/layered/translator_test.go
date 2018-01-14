package layered

import (
	"context"
	"testing"

	"github.com/klokare/evo"
	"github.com/klokare/network"
)

const name = "klima-layered"

var expectedForward = &Network{
	N:  []int{2, 2, 1},
	f:  []evo.Activation{evo.Direct, evo.Direct, evo.SteepenedSigmoid, evo.SteepenedSigmoid, evo.SteepenedSigmoid},
	np: []int{0, 2, 4},
	wp: []int{0, 0, 6},
	w: []float64{
		5.509324, -7.900839, 7.417331,
		-3.439597, -7.949483, 7.118832,
		3.187178, -3.771610, 2.606777,
	},
}

func TestTranslate(t *testing.T) {

	// Given a working XOR substrate
	sub := network.LayeredXor

	// We should get a network like this
	exp := expectedForward

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
	if len(exp.f) != len(act.f) {
		t.Errorf("incorrect number of activation functions: expected %d, actual %d", len(exp.f), len(act.f))
	} else {
		for i := 0; i < len(exp.f); i++ {
			if exp.f[i] != act.f[i] {
				t.Errorf("incorrect activation function at %d: expected %v, actual %v", i, exp.f[i], act.f[i])
			}
		}
	}

	// The layer sizes should be correct
	if len(exp.N) != len(act.N) {
		t.Errorf("incorrect number of layers: expected %d, actual %d", len(exp.N), len(act.N))
	} else {
		for i := 0; i < len(exp.N); i++ {
			if exp.N[i] != act.N[i] {
				t.Errorf("incorrect neuron count at layer %d: expected %v, actual %v", i, exp.N[i], act.N[i])
			}
		}
	}

	// The neuron pointer for each layer should be correct
	if len(exp.np) != len(act.np) {
		t.Errorf("incorrect number of neuron pointer records: expected %d, actual %d", len(exp.np), len(act.np))
	} else {
		for i := 0; i < len(exp.np); i++ {
			if exp.np[i] != act.np[i] {
				t.Errorf("incorrect neuron pointer at layer %d: expected %v, actual %v", i, exp.np[i], act.np[i])
			}
		}
	}

	// The weight pointer for each layer should be correct
	if len(exp.wp) != len(act.wp) {
		t.Errorf("incorrect number of weight pointer records: expected %d, actual %d", len(exp.wp), len(act.wp))
	} else {
		for i := 0; i < len(exp.wp); i++ {
			if exp.wp[i] != act.wp[i] {
				t.Errorf("incorrect weight pointer at layer %d: expected %v, actual %v", i, exp.wp[i], act.wp[i])
			}
		}
	}

	// The weight array should be correct
	if len(exp.w) != len(act.w) {
		t.Errorf("incorrect number of weights: expected %d, actual %d", len(exp.w), len(act.w))
	} else {
		for i := 0; i < len(exp.w); i++ {
			if exp.w[i] != act.w[i] {
				t.Errorf("incorrect weight at %d: expected %v, actual %v", i, exp.w[i], act.w[i])
			}
		}
	}
}

func BenchmarkTranslate(b *testing.B) {
	cases := network.GenerateCases(network.GenerateLayered)
	tr := network.Translator{Name: name, Translator: Translator{}}
	b.Run(name, network.TranslateBench(cases, tr))
}
