package network

import (
	"context"
	"fmt"
	"testing"

	"github.com/klokare/evo"
)

// Translator wraps the helper with a name
type Translator struct {
	Name string
	evo.Translator
}

// TranslateTest verifies that the translator returns errors for malformed substrates and forms
// a network from the XOR substrate that properly solves the problem.
func TranslateTest(cases []Case, tr Translator) func(*testing.T) {
	return func(t *testing.T) {
		// Iterate the translators
		ctx := context.Background()
		for _, c := range cases {
			t.Run(fmt.Sprintf("%s %s", tr.Name, c.Name), func(t *testing.T) {

				// Translate the network
				var net evo.Network
				var err error
				net, err = tr.Translator.Translate(ctx, c.Substrate)
				if t.Run("error", ErrorTest(c.HasError, err)); t.Failed() || c.HasError {
					return
				}

				// The network was the valid XOR network structure. Ensure it delivers the correct
				// result.
				t.Run("xor", XorTest(net))
			})
		}
	}
}

// TranslateBench provides a common function to benchmarking the translation of networks
func TranslateBench(cases []Case, tr Translator) func(*testing.B) {
	return func(b *testing.B) {
		for _, c := range cases {
			b.Run(fmt.Sprintf("%s-%s", tr.Name, c.Name), func(b *testing.B) {
				ctx := context.Background()
				for i := 0; i < b.N; i++ {
					tr.Translator.Translate(ctx, c.Substrate)
				}
			})
		}
	}
}
