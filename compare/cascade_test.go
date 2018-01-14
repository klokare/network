package compare

import (
	"context"
	"testing"

	"github.com/klokare/network"
	"github.com/klokare/network/klima/general"
	"github.com/klokare/network/klokare/forward"
)

var forwardTranslators = []network.Translator{
	{Name: "klima-general", Translator: general.Translator{}},
	{Name: "klokare-forward", Translator: forward.Translator{}},
}

func TestForwardTranslate(t *testing.T) {
	cases := []network.Case{
		{Name: "xor", Substrate: network.ForwardXor},
	}
	for _, tr := range forwardTranslators {
		t.Run(tr.Name, network.TranslateTest(cases, tr))
	}
}

func TestForwardActivate(t *testing.T) {
	for _, tr := range forwardTranslators {
		net, _ := tr.Translate(context.Background(), network.ForwardXor)
		t.Run(tr.Name, network.XorTest(net))
	}
}

func BenchmarkForwardTranslate(b *testing.B) {
	cases := network.GenerateCases(network.GenerateForward)
	for _, tr := range forwardTranslators {
		b.Run(tr.Name, network.TranslateBench(cases, tr))
	}
}

func BenchmarkForwardActivate(b *testing.B) {
	cases := network.GenerateCases(network.GenerateForward)
	for _, tr := range forwardTranslators {
		b.Run(tr.Name, network.ActivateBench(cases, tr))
	}
}
