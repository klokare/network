package layered

import (
	"testing"

	"github.com/klokare/network"
)

func TestActivate(t *testing.T) {
	t.Run(name, network.XorTest(expectedForward))
}

func BenchmarkActivate(b *testing.B) {
	cases := network.GenerateCases(network.GenerateLayered)
	tr := network.Translator{Name: name, Translator: Translator{}}
	b.Run(name, network.ActivateBench(cases, tr))
}
