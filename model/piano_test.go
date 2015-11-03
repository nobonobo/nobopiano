package model

import (
	"math"
	"testing"
)

func BenchmarkPianoModel(b *testing.B) {
	freqs := make([]float32, 14)
	for i := range freqs {
		freqs[i] = float32(440.0 * math.Pow(2.0, float64(i+60-69)/12.0))
	}
	piano := NewPiano(freqs)
	osc := piano.GetOscillator()
	for i := 0; i < b.N; i++ {
		osc()
	}
}
