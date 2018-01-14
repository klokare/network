package network

import (
	"fmt"

	"github.com/klokare/evo"
)

type Case struct {
	Name      string
	Substrate evo.Substrate
	HasError  bool
}

func GenerateCases(gen func(float64, ...int) evo.Substrate) []Case {
	names := []string{"small", "medium", "large", "xlarge"}
	dname := []string{"empty", "sparse", "dense", "full"}
	dense := []float64{0.0, 0.25, 0.75, 1.0}
	cases := make([]Case, 0, len(names)*len(dname))
	for i, name := range names {
		for j, dn := range dname {
			cases = append(cases, Case{
				Name:      fmt.Sprintf("%s-%s", name, dn),
				Substrate: gen(dense[j], Sizes[i]...),
			})
		}
	}
	return cases
}
