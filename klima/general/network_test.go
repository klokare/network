package general

import (
	"testing"

	"github.com/klokare/network"
)

func TestActivateLayered(t *testing.T) {
	t.Run(name, network.XorTest(expectedLayered))
}

func TestActivateForward(t *testing.T) {
	t.Run(name, network.XorTest(expectedForward))
}

func BenchmarkActivateLayered(b *testing.B) {
	cases := network.GenerateCases(network.GenerateLayered)
	tr := network.Translator{Name: name, Translator: Translator{}}
	b.Run(name, network.ActivateBench(cases, tr))
}

func BenchmarkActivateForward(b *testing.B) {
	cases := network.GenerateCases(network.GenerateForward)
	tr := network.Translator{Name: name, Translator: Translator{}}
	b.Run(name, network.ActivateBench(cases, tr))
}
