package dense

import (
	"context"
	"math/rand"
	"testing"

	"github.com/klokare/network"
	"gonum.org/v1/gonum/mat"
)

func TestActivate(t *testing.T) {
	t.Run(name, network.XorTest(expectedForward))
}

func BenchmarkActivate(b *testing.B) {
	cases := network.GenerateCases(network.GenerateLayered)
	tr := network.Translator{Name: name, Translator: Translator{}}
	b.Run(name, network.ActivateBench(cases, tr))
}

func TestSliceVsMatrix(t *testing.T) {

	// Create the input data
	n := 500                          // 500 rows of data
	m := 1000                         // 1000 columns
	islices := make([][]float64, n*m) // 1000 columns
	imatrix := mat.NewDense(n, m+1, nil)

	for i := 0; i < n; i++ {
		islices[i] = make([]float64, m)
		for j := 0; j < m; j++ {
			x := rand.Float64()
			islices[i][j] = x
			imatrix.Set(i, j, x)
		}
	}

	// Create the dense network
	sub := network.GenerateLayered(1.0, m, 1000, 100, 2)
	tmp, _ := new(Translator).Translate(context.Background(), sub)
	net := tmp.(Network)

	// Process as slice. Store output in matrix for easy compare
	sout := mat.NewDense(n, 2, nil)
	for i := 0; i < n; i++ {
		out, _ := net.Activate(islices[i])
		for j, v := range out {
			sout.Set(i, j, v)
		}
	}

	// Process as matrix
	var mout mat.Matrix
	mout, _ = net.ActivateMatrix(imatrix)

	// Compare dimensions
	sr, sc := sout.Dims()
	mr, mc := mout.Dims()
	if sr != mr {
		t.Errorf("incorrect number of rows: slices %d, matrix %d", sr, mr)
	}
	if sc != mc {
		t.Errorf("incorrect number of cols: slices %d, matrix %d", sc, mc)
	}

	// Compare the values
	for i := 0; i < sr; i++ {
		for j := 0; j < sc; j++ {
			sx := sout.At(i, j)
			mx := mout.At(i, j)
			if sx != mx {
				t.Errorf("incorrect value at [%d,%d]: expected %f, actual %f", i, j, sx, mx)
			}
		}
	}
}

func BenchmarkSliceVsMatrix(b *testing.B) {

	// Create the input data
	n := 500                          // 500 rows of data
	m := 1000                         // 1000 columns
	islices := make([][]float64, n*m) // 1000 columns
	imatrix := mat.NewDense(n, m+1, nil)

	for i := 0; i < n; i++ {
		islices[i] = make([]float64, m)
		for j := 0; j < m; j++ {
			x := rand.Float64()
			islices[i][j] = x
			imatrix.Set(i, j, x)
		}
	}
	b.ResetTimer()

	// Create the dense network
	sub := network.GenerateLayered(1.0, m, 1000, 1000, 100, 2)
	tmp, _ := new(Translator).Translate(context.Background(), sub)
	net := tmp.(Network)

	// Compare the two options
	b.Run("slices", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				for _, inputs := range islices {
					net.Activate(inputs)
				}
			}
		})
	})
	b.Run("matrix", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				net.ActivateMatrix(imatrix)
			}
		})
	})
}
