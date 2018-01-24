package network

import (
	"testing"

	"github.com/klokare/evo"
)

var (
	ForwardXor = evo.Substrate{
		Nodes: []evo.Node{
			{Position: evo.Position{Layer: 0.0, X: 0.0}, Neuron: evo.Input, Activation: evo.Direct},
			{Position: evo.Position{Layer: 0.0, X: 1.0}, Neuron: evo.Input, Activation: evo.Direct},
			{Position: evo.Position{Layer: 0.5, X: 0.5}, Neuron: evo.Hidden, Activation: evo.SteepenedSigmoid, Bias: -1.695151},
			{Position: evo.Position{Layer: 1.0, X: 0.5}, Neuron: evo.Output, Activation: evo.SteepenedSigmoid, Bias: -1.967445},
		},
		Conns: []evo.Conn{
			{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 0.5, X: 0.5}, Weight: 3.650676, Enabled: true},
			{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: -4.028692, Enabled: true},
			{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 0.5, X: 0.5}, Weight: -4.790058, Enabled: true},
			{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 3.972927, Enabled: true},
			{Source: evo.Position{Layer: 0.5, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 7.995010, Enabled: true},
		},
	}
	LayeredXor = evo.Substrate{
		Nodes: []evo.Node{
			{Position: evo.Position{Layer: 0.0, X: 0.0}, Neuron: evo.Input, Activation: evo.Direct},
			{Position: evo.Position{Layer: 0.0, X: 1.0}, Neuron: evo.Input, Activation: evo.Direct},
			{Position: evo.Position{Layer: 0.5, X: 0.0}, Neuron: evo.Hidden, Activation: evo.SteepenedSigmoid, Bias: 5.509324},
			{Position: evo.Position{Layer: 0.5, X: 1.0}, Neuron: evo.Hidden, Activation: evo.SteepenedSigmoid, Bias: -3.439597},
			{Position: evo.Position{Layer: 1.0, X: 0.5}, Neuron: evo.Output, Activation: evo.SteepenedSigmoid, Bias: 3.187178},
		},
		Conns: []evo.Conn{
			{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 0.5, X: 0.0}, Weight: -7.900839, Enabled: true},
			{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 0.5, X: 1.0}, Weight: -7.949483, Enabled: true},
			{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 0.5, X: 0.0}, Weight: 7.417331, Enabled: true},
			{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 0.5, X: 1.0}, Weight: 7.118832, Enabled: true},
			{Source: evo.Position{Layer: 0.5, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: -3.771610, Enabled: true},
			{Source: evo.Position{Layer: 0.5, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 2.606777, Enabled: true},
		},
	}
)

func XorTest(net evo.Network) func(*testing.T) {
	return func(t *testing.T) {
		// Evaluate XOR wit hthe network
		in, out, solved := XorEvaluate(net)
		if !solved {
			for i := 0; i < len(out); i++ {
				if i == 0 || i == 3 {
					if out[i] >= 0.5 {
						t.Errorf("incorrect outputs for %v: expected < 0.5, actual %f", in, out[i])
					}
				} else {
					if out[i] <= 0.5 {
						t.Errorf("incorrect outputs for %v: expected > 0.5, actual %f", in, out[i])
					}
				}
			}
		}
	}
}

func XorEvaluate(net evo.Network) (in [][]float64, out []float64, solved bool) {
	in = [][]float64{{0, 0}, {1, 0}, {0, 1}, {1, 1}}
	out = make([]float64, len(in))
	solved = true // be hopeful :)
	for i, inputs := range in {
		outputs, _ := net.Activate(inputs)
		out[i] = outputs[0]
		if i == 0 || i == 3 {
			solved = solved && outputs[0] <= 0.5
		} else {
			solved = solved && outputs[0] >= 0.5
		}
	}
	return
}
