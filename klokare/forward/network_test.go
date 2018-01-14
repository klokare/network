package forward

import (
	"testing"

	"github.com/klokare/network"
)

func TestActivateLayered(t *testing.T) {
	t.Run(name, network.XorTest(expectedLayered))
}

func TestActivateForward(t *testing.T) {
	t.Run(name, network.XorTest(expectedForward))
}

func BenchmarkActivateLayered(b *testing.B) {
	cases := network.GenerateCases(network.GenerateLayered)
	tr := network.Translator{Name: name, Translator: Translator{}}
	b.Run(name, network.ActivateBench(cases, tr))
}

func BenchmarkActivateForward(b *testing.B) {
	cases := network.GenerateCases(network.GenerateForward)
	tr := network.Translator{Name: name, Translator: Translator{}}
	b.Run(name, network.ActivateBench(cases, tr))
}

func TestActivateWrongInputSize(t *testing.T) {
	net := expectedForward
	_, err := net.Activate([]float64{1.0, 1.0, 1.0})
	if err == nil {
		t.Errorf("error expected")
	}
}

func TestActivateInvalidActivation(t *testing.T) {
	net := &Network{
		Inputs:   expectedForward.Inputs,
		Hidden:   expectedForward.Hidden,
		Outputs:  expectedForward.Outputs,
		Neurons:  make([]Neuron, len(expectedForward.Neurons)),
		Synapses: make([]Synapse, len(expectedForward.Synapses)),
	}
	copy(net.Neurons, expectedForward.Neurons)
	copy(net.Synapses, expectedForward.Synapses)
	net.Neurons[3].Activation = 0
	_, err := net.Activate([]float64{1.0, 1.0})
	if err == nil {
		t.Errorf("error expected")
	}
}
