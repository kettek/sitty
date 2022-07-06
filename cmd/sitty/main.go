package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/sitty/pkg/sitter"
)

func main() {
	sitter := sitter.Instance{}

	if err := sitter.Init(); err != nil {
		panic(err)
	}

	if err := ebiten.RunGame(&sitter); err != nil {
		panic(err)
	}
}
