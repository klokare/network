package layeredblas

import (
	"context"
	"sort"

	"github.com/klokare/evo"
)

type Translator struct{}

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

	// Begin the network
	n := make([]int, 0, 5)
	np := make([]int, 0, 5)
	wp := make([]int, 0, 5)
	f := make([]evo.Activation, 0, len(sub.Nodes))

	type lrec struct {
		Layer int
		Index int
	}
	m := make(map[evo.Position]lrec, len(sub.Nodes))
	last := sub.Nodes[0].Layer - 1.0
	l := 0
	i := 0
	for _, node := range sub.Nodes {

		// Begin the new layer
		if node.Layer > last {
			n = append(n, 0)
			l = len(n) - 1
			np = append(np, 0)
			wp = append(wp, 0)
			if l > 0 {
				np[l] = np[l-1] + n[l-1]
				if l > 1 {
					wp[l] = wp[l-1] + n[l-1]*(n[l-2]+1)
				}
			}
			i = 0
		}

		// Add the node
		n[l]++
		f = append(f, node.Activation)
		m[node.Position] = lrec{Layer: l, Index: i}
		i++
		last = node.Layer
	}

	// Write the bias values
	w := make([]float64, wp[l]+n[l]*(n[l-1]+1))
	for _, node := range sub.Nodes {
		rec := m[node.Position]
		l := rec.Layer
		i := rec.Index
		if l > 0 {
			wi := wp[l] + (n[l-1]+1)*i
			w[wi] = node.Bias
		}
	}

	// Write the connection weights
	for _, conn := range sub.Conns {
		if conn.Enabled {
			src := m[conn.Source]
			tgt := m[conn.Target]
			i := tgt.Index
			j := src.Index
			l := tgt.Layer
			wi := wp[l] + (n[l-1]+1)*i + j + 1
			w[wi] = conn.Weight
		}
	}

	net = &Network{
		N:  n,
		f:  f,
		w:  w,
		np: np,
		wp: wp,
	}
	return
}
