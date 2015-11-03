// +build darwin
// +build arm arm64

package main

import (
	"log"

	"golang.org/x/mobile/app"
)

func init() {
	log.Println("ios mode")
}

func repaint(a app.App) {}
