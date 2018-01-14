package network

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
)

func ActivateBench(cases []Case, tr Translator) func(*testing.B) {
	return func(b *testing.B) {

		// Create the input data
		inputs := make([][]float64, len(Sizes))
		for i, sizes := range Sizes {
			inputs[i] = make([]float64, sizes[0]-1)
			for j := 0; j < len(inputs[i]); j++ {
				inputs[i][j] = rand.Float64()
			}
		}
		b.ResetTimer()

		// Execute the benchmarks
		for j, c := range cases {
			ds := j / 4
			b.Run(fmt.Sprintf("%s-%s", tr.Name, c.Name), func(b *testing.B) {

				// Translate the network
				ctx := context.Background()
				net, err := tr.Translator.Translate(ctx, c.Substrate)
				if err != nil {
					b.Log("error translating network:", err)
					return
				}
				b.ReportAllocs()
				b.ResetTimer()

				// Activate the network
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						net.Activate(inputs[ds])
					}
				})
			})
		}
	}
}
