package dense

import (
	"context"
	"testing"

	"github.com/klokare/evo"
	"github.com/klokare/network"
	"gonum.org/v1/gonum/mat"
)

const name = "matrix-dense"

var expectedForward = Network{
	Layers: [][]evo.Activation{
		[]evo.Activation{evo.Direct, evo.Direct, evo.Direct},
		[]evo.Activation{evo.SteepenedSigmoid, evo.SteepenedSigmoid, evo.Direct},
		[]evo.Activation{evo.SteepenedSigmoid},
	},
	Weights: []*mat.Dense{
		mat.NewDense(3, 3, []float64{
			-7.900839, -7.949483, 0,
			7.417331, 7.118832, 0,
			5.509324, -3.439597, 0,
		}),
		mat.NewDense(3, 1, []float64{
			-3.771610,
			2.606777,
			3.187178,
		}),
	},
}

func TestTranslate(t *testing.T) {

	// XOR substrate
	sub := network.LayeredXor

	// Expected network
	exp := expectedForward

	// Translate the network
	tmp, err := new(Translator).Translate(context.Background(), sub)
	if err != nil {
		t.Errorf("there should be no error. actual %v", err)
	}
	if tmp == nil {
		t.Errorf("network not returned")
	}
	act := tmp.(Network)

	// There should be the same activations
	if len(exp.Layers) != len(act.Layers) {
		t.Errorf("incorrect number of layers: expected %d, actual %d", len(exp.Layers), len(act.Layers))
	}
	for l := 0; l < len(exp.Layers); l++ {
		if len(exp.Layers[l]) != len(act.Layers[l]) {
			t.Errorf("incorrect number of activations for layer %d: expected %d, actual %d", l, len(exp.Layers[l]), len(act.Layers[l]))
		}
		for a := 0; a < len(exp.Layers[l]); a++ {
			if exp.Layers[l][a] != act.Layers[l][a] {
				t.Errorf("incorrect activation %d  in layer %d: expected %s, actual %s", a, l, exp.Layers[l][a], act.Layers[l][a])
			}
		}
	}

	// There should be the same weights
	if len(exp.Weights) != len(act.Weights) {
		t.Errorf("incorrect number of weight matrices: expected %d, actual %d", len(exp.Weights), len(act.Weights))
	}
	for w := 0; w < len(exp.Weights); w++ {
		em := exp.Weights[w]
		am := act.Weights[w]
		er, ec := em.Dims()
		ar, ac := am.Dims()
		if er != ar || ec != ac {
			t.Errorf("incorrected dimension of matrix %d: expected [%d,%d], actual [%d,%d]", w, er, ec, ar, ac)
		}
		for r := 0; r < er; r++ {
			for c := 0; c < ec; c++ {
				ex := em.At(r, c)
				ax := am.At(r, c)
				if ex != ax {
					t.Errorf("incorrect value at [%d,%d] of matrix %d: expected %f, actual %f", r, c, w, ex, ax)
				}
			}
		}
	}
}

func BenchmarkTranslate(b *testing.B) {
	cases := network.GenerateCases(network.GenerateLayered)
	tr := network.Translator{Name: name, Translator: Translator{}}
	b.Run(name, network.TranslateBench(cases, tr))
}
