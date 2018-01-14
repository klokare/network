package general

import (
	"context"
	"sort"

	"github.com/klokare/evo"
)

// Translator ...
type Translator struct{}

// Translate ...
func (t Translator) Translate(ctx context.Context, sub evo.Substrate) (net evo.Network, err error) {

	// Ensure the substrate order
	sort.Slice(sub.Nodes, func(i, j int) bool { return sub.Nodes[i].Compare(sub.Nodes[j]) < 0 })
	sort.Slice(sub.Conns, func(i, j int) bool {
		x := sub.Conns[i].Target.Compare(sub.Conns[j].Target)
		if x == 0 {
			return sub.Conns[i].Source.Compare(sub.Conns[j].Source) < 0
		}
		return x < 0
	})

	// Map the nodes to their indexes and count the neuron types
	var n, ni, nh, no int
	n = len(sub.Nodes)
	f := make([]evo.Activation, n)
	n2c := make(map[int][]int, n)
	n2i := make(map[evo.Position]int, n)
	for i, node := range sub.Nodes {
		switch node.Neuron {
		case evo.Input:
			ni++
		case evo.Hidden:
			nh++
		case evo.Output:
			no++
		}
		n2i[node.Position] = i
		f[i] = node.Activation
	}

	// Map the incomming connections to their nodes
	var in []int
	var ok bool
	for i, conn := range sub.Conns {
		ti := n2i[conn.Target]
		if in, ok = n2c[ti]; !ok {
			in = make([]int, 0, 10)
		}
		in = append(in, i)
		n2c[ti] = in
	}

	// Set the neuron and weight pointers
	np := make([]int, n+1)
	nc := make([]int, n)
	for i := 1; i < n; i++ {
		np[i] = np[i-1] + nc[i-1]
		nc[i] = len(n2c[i])
		if sub.Nodes[i].Neuron != evo.Input {
			nc[i]++
		}
	}
	np[n] = np[n-1] + nc[n-1]

	// Set the weights and the index of source node
	w := make([]float64, np[n])
	m := make([]int, len(w))
	wi := 0
	for i, node := range sub.Nodes {

		// Add the connections
		for _, c := range n2c[i] {
			w[wi] = sub.Conns[c].Weight
			m[wi] = n2i[sub.Conns[c].Source]
			wi++
		}

		// Add the bias
		if node.Neuron != evo.Input {
			w[wi] = node.Bias
			m[wi] = i
			wi++
		}
	}

	// Create the network
	net = &Network{
		NP: np,
		NC: nc,
		W:  w,
		M:  m,
		F:  f,
		N:  n,
		NI: ni,
		NH: nh,
		NO: no,
	}
	return
}
