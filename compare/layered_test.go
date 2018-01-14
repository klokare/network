package compare

import (
	"context"
	"testing"

	"github.com/klokare/network"
	"github.com/klokare/network/klima/general"
	"github.com/klokare/network/klima/layered"
	"github.com/klokare/network/klima/layeredblas"
	"github.com/klokare/network/klokare/forward"
	"github.com/klokare/network/matrix/dense"
	"github.com/klokare/network/matrix/sparse"
)

var layeredTranslators = []network.Translator{
	{Name: "matrix-dense", Translator: dense.Translator{}},
	{Name: "matrix-sparse", Translator: sparse.Translator{}},
	{Name: "klima-forward", Translator: layered.Translator{}},
	{Name: "klima-blas", Translator: layeredblas.Translator{}},
	{Name: "klima-general", Translator: general.Translator{}},
	{Name: "forward-push", Translator: forward.Translator{}},
}

func TestLayeredTranslate(t *testing.T) {
	cases := []network.Case{
		{Name: "xor", Substrate: network.LayeredXor},
	}
	for _, tr := range layeredTranslators {
		t.Run(tr.Name, network.TranslateTest(cases, tr))
	}
}

func TestLayeredActivate(t *testing.T) {
	for _, tr := range layeredTranslators {
		net, _ := tr.Translate(context.Background(), network.LayeredXor)
		t.Run(tr.Name, network.XorTest(net))
	}
}

func BenchmarkLayeredTranslate(b *testing.B) {
	cases := network.GenerateCases(network.GenerateLayered)
	for _, tr := range layeredTranslators {
		b.Run(tr.Name, network.TranslateBench(cases, tr))
	}
}

func BenchmarkLayeredActivate(b *testing.B) {
	cases := network.GenerateCases(network.GenerateLayered)
	for _, tr := range forwardTranslators {
		b.Run(tr.Name, network.ActivateBench(cases, tr))
	}
}
