package main

import (
	"encoding/binary"
	"log"
	"sync"

	"golang.org/x/mobile/exp/audio/al"

	"github.com/nobonobo/nobopiano/model"
)

const (
	SampleRate = 22050
	NS         = 256 // number of samples
	CZ         = 2   // bytes/1-sample for al.FormatMono16
	Fmt        = al.FormatMono16
	QUEUE      = 8
)

type Audio struct {
	sync.RWMutex
	source     al.Source
	queue      []al.Buffer
	oscillator model.Oscillator
}

func NewAudio(oscillator model.Oscillator) *Audio {
	if err := al.OpenDevice(); err != nil {
		log.Fatal(err)
	}
	s := al.GenSources(1)
	if code := al.Error(); code != 0 {
		log.Fatalln("openal error:", code)
	}
	return &Audio{
		source:     s[0],
		queue:      []al.Buffer{},
		oscillator: oscillator,
	}
}

func (c *Audio) Play() {
	c.Lock()
	defer c.Unlock()
	n := c.source.BuffersProcessed()
	if n > 0 {
		rm, split := c.queue[:n], c.queue[n:]
		c.queue = split
		c.source.UnqueueBuffers(rm...)
		al.DeleteBuffers(rm...)
	}
	for len(c.queue) < QUEUE {
		b := al.GenBuffers(1)
		buf := make([]byte, NS*CZ)
		for n := 0; n < NS*CZ; n += CZ {
			f := c.oscillator()
			if f < -1.0 {
				f = -1.0
			}
			if f > 1.0 {
				f = 1.0
			}
			v := int16(float32(32767) * f)
			binary.LittleEndian.PutUint16(buf[n:n+2], uint16(v))
		}
		b[0].BufferData(Fmt, buf, SampleRate)
		c.source.QueueBuffers(b...)
		c.queue = append(c.queue, b...)
	}
	if c.source.State() != al.Playing {
		al.PlaySources(c.source)
	}
}

func (c *Audio) Close() {
	c.Lock()
	defer c.Unlock()
	al.StopSources(c.source)
	c.source.UnqueueBuffers(c.queue...)
	al.DeleteBuffers(c.queue...)
	c.queue = nil
	al.DeleteSources(c.source)
}
