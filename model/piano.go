package model

const N = 4

type Piano struct {
	notes      []bool
	oscillator Oscillator
}

func NewPiano(freqs []float32) *Piano {
	p := new(Piano)
	p.notes = make([]bool, len(freqs))
	envelopes := []Oscillator{}
	for i, f := range freqs {
		base := []Oscillator{}
		for j := float32(1.0); j <= N; j++ {
			base = append(base, G(0.5/j, GenOscillator(f*j)))
		}
		base = append(base, G(0.3, GenOscillator(f+2)))
		osc := Multiplex(base...)
		envelopes = append(envelopes, G(0.4, GenEnvelope(&p.notes[i], osc)))
	}
	p.oscillator = Multiplex(envelopes...) // all note oscilator multiplex
	return p
}

func (p *Piano) NoteOn(key int) {
	p.notes[key] = true
}

func (p *Piano) NoteOff(key int) {
	p.notes[key] = false
}

func (p *Piano) GetOscillator() Oscillator { return p.oscillator }
