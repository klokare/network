package network

import (
	"math/rand"
	"sort"

	"github.com/klokare/evo"
)

// Sizes of generated networks
var Sizes = [][]int{
	{10, 5, 2},
	{100, 40, 10, 4},
	{1000, 250, 250, 50},
	{5000, 1000, 500, 250, 100},
}

func GenerateForward(density float64, sizes ...int) evo.Substrate {

	// Density is the proportion of connections that remain, so 1 - density is the probability
	// that it is removed
	p := 1.0 - density

	// Create the layers
	layers := make([][]evo.Node, 0, len(sizes))
	l := 0.0
	dl := 1.0 / float64(len(sizes)-1)
	tot := 0
	for i, cnt := range sizes {

		nodes := make([]evo.Node, 0, cnt)

		// Determine the starting x value and its increment
		x := 0.5
		dx := 0.0
		if cnt > 1 {
			x = 0.0
			dx = 1.0 / float64(cnt-1)
		}
		for j := 0; j < cnt; j++ {
			node := evo.Node{Position: evo.Position{Layer: l, X: x}}

			// Create the right kind of node
			switch i {
			case 0: // Input layer
				node.Activation = evo.Direct
				node.Neuron = evo.Input
			case len(sizes) - 1: // Output layer
				node.Activation = evo.Sigmoid
				node.Neuron = evo.Output
			default: // Hidden layer
				node.Activation = evo.Sigmoid
				node.Neuron = evo.Hidden
			}

			// Append the node and increment x
			nodes = append(nodes, node)
			x += dx
		}

		// Append the layer and increment layer
		layers = append(layers, nodes)
		l += dl
		tot += len(nodes)
	}

	// Fully connect the layers to every prior layer
	conns := make([]evo.Conn, 0, sizes[0]*sizes[1]*len(sizes)-1)
	if density > 0.0 {

		// Iterate the layers, connecting to every prior layer
		for t := len(layers) - 1; t > 0; t-- {
			tlayer := layers[t]
			for _, tgt := range tlayer {
				for s := 0; s < t; s++ {
					slayer := layers[s]
					for _, src := range slayer {
						conns = append(conns, evo.Conn{
							Source:  src.Position,
							Target:  tgt.Position,
							Weight:  rand.Float64(),
							Enabled: true,
						})
					}
				}
			}
		}
	}

	// Prune the substrate based on the probability.
	if density < 1.0 {
		tmp := conns
		conns = make([]evo.Conn, 0, len(conns))
		for _, c := range tmp {
			if rand.Float64() < p {
				continue // removes the connection
			}
			conns = append(conns, c)
		}
	}

	// Sort
	sort.Slice(conns, func(i, j int) bool { return conns[i].Compare(conns[j]) < 0 })

	// Return the substrate
	sub := evo.Substrate{
		Nodes: make([]evo.Node, 0, tot),
		Conns: conns,
	}
	for _, layer := range layers {
		sub.Nodes = append(sub.Nodes, layer...)
	}
	return sub
}

func GenerateLayered(density float64, sizes ...int) evo.Substrate {

	// Density is the proportion of connections that remain, so 1 - density is the probability
	// that it is removed
	p := 1.0 - density

	// Create the layers
	nodes := make([]evo.Node, 0, sizes[0]*len(sizes))
	l := 0.0
	dl := 1.0 / float64(len(sizes)-1)
	for i, cnt := range sizes {

		// Determine the starting x value and its increment
		x := 0.5
		dx := 0.0
		if cnt > 1 {
			x = 0.0
			dx = 1.0 / float64(cnt-1)
		}
		for j := 0; j < cnt; j++ {
			node := evo.Node{Position: evo.Position{Layer: l, X: x}}

			// Create the right kind of node
			switch i {
			case 0: // Input layer
				node.Activation = evo.Direct
				node.Neuron = evo.Input
			case len(sizes) - 1: // Output layer
				node.Activation = evo.Sigmoid
				node.Neuron = evo.Output
			default: // Hidden layer
				node.Activation = evo.Sigmoid
				node.Neuron = evo.Hidden
			}

			// Append the node and increment x
			nodes = append(nodes, node)
			x += dx
		}

		// Increment layer
		l += dl
	}

	// Fully connect the layers
	conns := make([]evo.Conn, 0, sizes[0]*sizes[1]*len(sizes)-1)
	if density > 0.0 {

		// Iterate the layers
		offset := 0
		for l := 0; l < len(sizes)-1; l++ {
			scnt := sizes[l+0]
			tcnt := sizes[l+1]
			for s := 0; s < scnt; s++ {
				src := nodes[offset+s]
				for t := 0; t < tcnt; t++ {
					tgt := nodes[offset+scnt+t]
					conns = append(conns, evo.Conn{
						Source:  src.Position,
						Target:  tgt.Position,
						Weight:  rand.NormFloat64(),
						Enabled: true,
					})
				}
			}

			// Increment the offset
			offset += scnt
		}
	}

	// Prune the substrate based on the probability.
	if density < 1.0 {
		tmp := conns
		conns = make([]evo.Conn, 0, len(conns))
		for _, c := range tmp {
			if rand.Float64() < p {
				continue // removes the connection
			}
			conns = append(conns, c)
		}
	}

	// Return the substrate
	return evo.Substrate{
		Nodes: nodes,
		Conns: conns,
	}
}
