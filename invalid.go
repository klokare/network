package network

import "github.com/klokare/evo"

var (
	InvalidInputs = evo.Substrate{
		Nodes: []evo.Node{
			{Position: evo.Position{Layer: 0.5, X: 0.5}, Neuron: evo.Hidden, Activation: evo.SteepenedSigmoid},
			{Position: evo.Position{Layer: 1.0, X: 0.5}, Neuron: evo.Output, Activation: evo.SteepenedSigmoid},
		},
	}
	InvalidOutputs = evo.Substrate{
		Nodes: []evo.Node{
			{Position: evo.Position{Layer: 0.0, X: 1.0}, Neuron: evo.Input, Activation: evo.Direct},
			{Position: evo.Position{Layer: 0.0, X: 1.0}, Neuron: evo.Input, Activation: evo.Direct},
			{Position: evo.Position{Layer: 0.5, X: 0.5}, Neuron: evo.Hidden, Activation: evo.SteepenedSigmoid},
		},
	}
)
