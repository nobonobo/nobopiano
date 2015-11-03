package model

import (
	"math"

	"golang.org/x/mobile/exp/f32"
)

var SampleRate = 22050

const (
	Pi = float32(math.Pi)
)

type Oscillator func() float32

func G(gain float32, f Oscillator) Oscillator {
	return func() float32 {
		return gain * f()
	}
}

func Multiplex(fs ...Oscillator) Oscillator {
	return func() float32 {
		res := float32(0)
		for _, osc := range fs {
			res += osc()
		}
		return res
	}
}

func GenOscillator(freq float32) Oscillator {
	dt := 1.0 / float32(SampleRate)
	k := 2.0 * Pi * freq
	T := 1.0 / freq
	t := float32(0.0)
	return func() float32 {
		res := f32.Sin(k * t)
		t += dt
		if t > T {
			t -= T
		}
		return res
	}
}

func GenEnvelope(press *bool, f Oscillator) Oscillator {
	dt := 1.0 / float32(SampleRate)
	top := false
	gain := float32(0.0)
	attackd := dt / 0.01
	dekeyd := dt / 0.03
	sustainlevel := float32(0.3)
	sustaind := dt / 7.0
	released := dt / 0.8
	return func() float32 {
		if *press {
			if !top {
				gain += attackd
				if gain > 1.0 {
					top = true
					gain = 1.0
				}
			} else {
				if gain > sustainlevel {
					gain -= dekeyd
				} else {
					gain -= sustaind
				}
				if gain < 0.0 {
					gain = 0.0
				}
			}
		} else {
			top = false
			gain -= released
			if gain < 0.0 {
				gain = 0.0
			}
		}
		return gain * f()
	}
}
