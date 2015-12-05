package main

import (
	"fmt"
	"image"
	"log"
	"time"

	_ "image/png"

	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/app/debug"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/exp/sprite/clock"
	"golang.org/x/mobile/exp/sprite/glsprite"
	"golang.org/x/mobile/gl"

	"github.com/nobonobo/nobopiano/model"
)

var (
	whites  = []int{0, 1, 3, 5, 6, 8, 10, 12, 13}
	blacks  = []int{-1, -1, 2, 4, -1, 7, 9, 11, -1, -1}
	whiltsN = float32(len(whites))
)

type Context struct {
	piano     *model.Piano
	audio     *Audio
	lastSeq   map[touch.Sequence]int
	startTime time.Time
	images    *glutil.Images
	eng       sprite.Engine
	scene     *sprite.Node
	sz        size.Event
	fps       *debug.FPS
}

func NewContext() *Context {
	model.SampleRate = SampleRate
	piano := model.NewPiano([]float32{
		246.941650628,
		261.625565301,
		277.182630977,
		293.664767917,
		311.126983722,
		329.627556913,
		349.228231433,
		369.994422712,
		391.995435982,
		415.30469758,
		440.0,
		466.163761518,
		493.883301256,
		523.251130601,
	})
	lastSeq := map[touch.Sequence]int{}

	return &Context{
		piano:     piano,
		lastSeq:   lastSeq,
		startTime: time.Now(),
	}
}

func (c *Context) Start(glctx gl.Context) {
	c.images = glutil.NewImages(glctx)
	c.eng = glsprite.Engine(c.images)
	c.loadScene()
	log.Println("Start")
	c.audio = NewAudio(c.piano.GetOscillator())
	c.fps = debug.NewFPS(c.images)
}

func (c *Context) Stop() {
	if c.audio != nil {
		c.audio.Close()
		c.audio = nil
	}
	log.Println("Stop")
	c.fps.Release()
	c.eng.Release()
	c.images.Release()
}

func (c *Context) Paint(glctx gl.Context) {
	glctx.ClearColor(1, 1, 1, 1)
	glctx.Clear(gl.COLOR_BUFFER_BIT)
	now := clock.Time(time.Since(c.startTime) * 60 / time.Second)
	c.eng.Render(c.scene, now, c.sz)
	c.fps.Draw(c.sz)
}

func (c *Context) Play() {
	if c.audio != nil {
		c.audio.Play()
	}
}

func (c *Context) Size(sz size.Event) {
	c.sz = sz
}

func (c *Context) hit2key(x, y float32) int {
	w := float32(c.sz.WidthPx)
	h := float32(c.sz.HeightPx)
	offset := w / whiltsN / 2
	if y < h*0.7 {
		index := int(whiltsN * (x + offset) / w)
		if index >= len(blacks) {
			index = len(blacks) - 1
		}
		if key := blacks[index]; key >= 0 {
			return key
		}
	}
	index := int(whiltsN * x / w)
	if index >= len(whites) {
		index = len(whites) - 1
	}
	return whites[index]
}

func (c *Context) Touch(e touch.Event) {
	key := c.hit2key(e.X, e.Y)
	old, ok := c.lastSeq[e.Sequence]
	keyChanged := !ok || (old != key)
	if ok && keyChanged {
		c.piano.NoteOff(old)
		fmt.Println("notechoff:", old, e.X, e.Y, e.Type, e.Sequence)
	}
	c.lastSeq[e.Sequence] = key
	if e.Type == touch.TypeBegin || (keyChanged && e.Type == touch.TypeMove) {
		c.piano.NoteOn(key)
		fmt.Println("noteon:", key, e.X, e.Y, e.Type, e.Sequence)
	}
	if e.Type == touch.TypeEnd {
		c.piano.NoteOff(key)
		fmt.Println("noteoff:", key, e.X, e.Y, e.Type, e.Sequence)
	}
}

func (c *Context) loadScene() {
	keyboard := loadKeyboard(c.eng)
	c.scene = &sprite.Node{}
	c.eng.Register(c.scene)
	c.eng.SetTransform(c.scene, f32.Affine{
		{1, 0, 0},
		{0, 1, 0},
	})

	n := &sprite.Node{}
	c.eng.Register(n)
	c.scene.AppendChild(n)
	// TODO: Shouldn't arranger pass the size.Event?
	n.Arranger = arrangerFunc(func(eng sprite.Engine, n *sprite.Node, t clock.Time) {
		eng.SetSubTex(n, keyboard)
		eng.SetTransform(n, f32.Affine{
			{float32(c.sz.WidthPt), 0, 0},
			{0, float32(c.sz.HeightPt), 0},
		})
	})
}

func loadKeyboard(eng sprite.Engine) sprite.SubTex {
	a, err := asset.Open("piano-octave.png")
	if err != nil {
		log.Fatal(err)
	}
	defer a.Close()

	img, _, err := image.Decode(a)
	if err != nil {
		log.Fatal(err)
	}
	t, err := eng.LoadTexture(img)
	if err != nil {
		log.Fatal(err)
	}
	return sprite.SubTex{t, image.Rect(0, 0, 500, 249)}
}

type arrangerFunc func(e sprite.Engine, n *sprite.Node, t clock.Time)

func (a arrangerFunc) Arrange(e sprite.Engine, n *sprite.Node, t clock.Time) { a(e, n, t) }
