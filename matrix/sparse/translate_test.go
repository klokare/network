package sparse

import (
	"testing"

	"github.com/klokare/network"
)

const name = "matrix-sparse"

func BenchmarkTranslate(b *testing.B) {
	cases := network.GenerateCases(network.GenerateLayered)
	tr := network.Translator{Name: name, Translator: Translator{}}
	b.Run(name, network.TranslateBench(cases, tr))
}
