package main

import (
	"log"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/gl"
)

func main() {
	ctx := NewContext()
	app.Main(func(a app.App) {
		var glctx gl.Context
		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					log.Println("CrossOn")
					glctx, _ = e.DrawContext.(gl.Context)
					ctx.Start(glctx)
					repaint(a) //a.Send(paint.Event{})
				case lifecycle.CrossOff:
					log.Println("CrossOff")
					ctx.Stop()
					glctx = nil
				}
			case size.Event:
				log.Println(e)
				ctx.Size(e)
			case paint.Event:
				if glctx == nil || e.External {
					continue
				}
				ctx.Paint(glctx)
				a.Publish()
				ctx.Play()
				repaint(a) // keep animating
			case touch.Event:
				ctx.Touch(e)
			}
		}
	})
}
