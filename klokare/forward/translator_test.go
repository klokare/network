package forward

import (
	"context"
	"testing"

	"github.com/klokare/evo"
	"github.com/klokare/network"
)

const name = "klokare-forward"

var expectedLayered = &Network{
	Inputs:  2,
	Hidden:  2,
	Outputs: 1,
	Neurons: []Neuron{
		{Neuron: evo.Input, Activation: evo.Direct, Bias: 0.0},
		{Neuron: evo.Input, Activation: evo.Direct, Bias: 0.0},
		{Neuron: evo.Hidden, Activation: evo.SteepenedSigmoid, Bias: 5.509324},
		{Neuron: evo.Hidden, Activation: evo.SteepenedSigmoid, Bias: -3.439597},
		{Neuron: evo.Output, Activation: evo.SteepenedSigmoid, Bias: 3.187178},
	},
	Synapses: []Synapse{
		{Source: 0, Target: 2, Weight: -7.900839},
		{Source: 0, Target: 3, Weight: -7.949483},
		{Source: 1, Target: 2, Weight: 7.417331},
		{Source: 1, Target: 3, Weight: 7.118832},
		{Source: 2, Target: 4, Weight: -3.771610},
		{Source: 3, Target: 4, Weight: 2.606777},
	},
}

var expectedForward = &Network{
	Inputs:  2,
	Hidden:  1,
	Outputs: 1,
	Neurons: []Neuron{
		{Neuron: evo.Input, Activation: evo.Direct, Bias: 0.0},
		{Neuron: evo.Input, Activation: evo.Direct, Bias: 0.0},
		{Neuron: evo.Hidden, Activation: evo.SteepenedSigmoid, Bias: -1.695151},
		{Neuron: evo.Output, Activation: evo.SteepenedSigmoid, Bias: -1.967445},
	},
	Synapses: []Synapse{
		{Source: 0, Target: 2, Weight: 3.650676},
		{Source: 0, Target: 3, Weight: -4.028692},
		{Source: 1, Target: 2, Weight: -4.790058},
		{Source: 1, Target: 3, Weight: 3.972927},
		{Source: 2, Target: 3, Weight: 7.995010},
	},
}

func TestInvalid(t *testing.T) {
	cases := []network.Case{
		{Name: "invalid inputs", Substrate: network.InvalidInputs, HasError: true},
		{Name: "invalid outputs", Substrate: network.InvalidOutputs, HasError: true},
	}
	t.Run("invalid", network.TranslateTest(cases, network.Translator{Name: name, Translator: new(Translator)}))
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

	// There should be the correct neurons
	if exp.Inputs != act.Inputs {
		t.Errorf("incorrect number of input neurons: expected %d, actual %d", exp.Inputs, act.Inputs)
	}
	if exp.Hidden != act.Hidden {
		t.Errorf("incorrect number of hidden neurons: expected %d, actual %d", exp.Hidden, act.Hidden)
	}
	if exp.Outputs != act.Outputs {
		t.Errorf("incorrect number of output neurons: expected %d, actual %d", exp.Outputs, act.Outputs)
	}
	if len(exp.Neurons) != len(act.Neurons) {
		t.Errorf("inccorect number of neurons: expected %d, actual %d", len(exp.Neurons), len(act.Neurons))
	}
	for i := 0; i < len(exp.Neurons); i++ {
		if exp.Neurons[i].Neuron != act.Neurons[i].Neuron {
			t.Errorf("incorrect neuron type for %d: expected %d, actual %d", i, exp.Neurons[i].Neuron, act.Neurons[i].Neuron)
		}
		if exp.Neurons[i].Activation != act.Neurons[i].Activation {
			t.Errorf("incorrect activation type for %d: expected %d, actual %d", i, exp.Neurons[i].Activation, act.Neurons[i].Activation)
		}
		if exp.Neurons[i].Bias != act.Neurons[i].Bias {
			t.Errorf("incorrect bias for %d: expected %f, actual %f", i, exp.Neurons[i].Bias, act.Neurons[i].Bias)
		}
	}

	// There should be the correct synapses
	if len(exp.Synapses) != len(act.Synapses) {
		t.Errorf("inccorect number of synapses: expected %d, actual %d", len(exp.Synapses), len(act.Synapses))
	}
	for i := 0; i < len(exp.Synapses); i++ {
		if exp.Synapses[i].Source != act.Synapses[i].Source {
			t.Errorf("incorrect source for %d: expected %d, actual %d", i, exp.Synapses[i].Source, act.Synapses[i].Source)
		}
		if exp.Synapses[i].Target != act.Synapses[i].Target {
			t.Errorf("incorrect target for %d: expected %d, actual %d", i, exp.Synapses[i].Target, act.Synapses[i].Target)
		}
		if exp.Synapses[i].Weight != act.Synapses[i].Weight {
			t.Errorf("incorrect weight for %d: expected %d, actual %d", i, exp.Synapses[i].Weight, act.Synapses[i].Weight)
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
