package sparse

import (
	"context"
	"sort"

	"github.com/james-bowman/sparse"
	"github.com/klokare/evo"
)

type Translator struct{}

func (z Translator) Translate(ctx context.Context, sub evo.Substrate) (evo.Network, error) {

	// Sort substrate nodes by source then target so that all layers line up
	sort.Slice(sub.Nodes, func(i, j int) bool { return sub.Nodes[i].Compare(sub.Nodes[j]) < 0 })

	// Create the layers while mapping the nodes to the indices of the layer
	type rec struct {
		Layer int
		Index int
	}
	m := make(map[evo.Position]rec, len(sub.Nodes))
	layers := make([][]evo.Activation, 0, 10)
	last := sub.Nodes[0].Layer - 1.0 // last layer read
	var l, idx int
	for _, n := range sub.Nodes {
		if n.Layer > last {
			layers = append(layers, make([]evo.Activation, 0, 100))
			idx = 0
			l = len(layers) - 1
		}
		layers[l] = append(layers[l], n.Activation)
		m[n.Position] = rec{Layer: l, Index: idx}
		idx++
		last = n.Layer
	}

	// Append the dummy bias node to each layer
	for i := 0; i < len(layers)-1; i++ {
		layers[i] = append(layers[i], evo.Direct)
	}

	// Create the matrices for the weights beetween layers
	weights := make([]*sparse.CSR, 0, len(layers)-1)
	for i := 1; i < len(layers); i++ {

		// Create the data matrix for this
		weights = append(weights, sparse.NewDOK(len(layers[i-1]), len(layers[i])).ToCSR())
	}

	// Set the bias values
	for _, n := range sub.Nodes {
		rec := m[n.Position]
		if rec.Layer > 0 {
			w := rec.Layer - 1
			weights[w].Set(len(layers[rec.Layer-1])-1, rec.Index, n.Bias)
		}
	}

	// Set weights
	for _, c := range sub.Conns {
		if c.Enabled {
			src := m[c.Source]
			tgt := m[c.Target]

			// Special case, bias
			if src.Index == -1 {
				src.Index = len(layers[src.Layer]) - 1
			}

			// Set the value
			w := tgt.Layer - 1
			weights[w].Set(src.Index, tgt.Index, c.Weight)
		}
	}

	// Return the new network
	return Network{
		Layers:  layers[:len(layers)],
		Weights: weights[:len(weights)],
	}, nil
}
