package sitter

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/sitty/pkg/winter"
	"github.com/kettek/sitty/sitters"
)

type Instance struct {
	sitter                    *sitters.Sitter
	sitterImageTicker         int
	winter                    *winter.Winter
	rootX, rootY              int
	currentX, currentY        int
	targetX, targetY          int
	targetWidth, targetHeight int
	moveTimer                 int
}

func (i *Instance) Init() error {

	// Set up our ebiten window.
	ebiten.SetInitFocused(false)
	ebiten.SetScreenTransparent(true)
	ebiten.SetWindowDecorated(false)
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowSize(64, 64)
	ebiten.SetWindowPosition(32, 32)
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetWindowTitle("Go Sitter")

	// Get our initial sitter.
	s, err := sitters.LoadSitter("gopher")
	if err != nil {
		panic(err)
	}
	i.sitter = s
	if err := i.sitter.Init(); err != nil {
		panic(err)
	}

	// Set up our window manager interactor.
	w, err := winter.NewWinter()
	if err != nil {
		return err
	}
	i.winter = w

	return nil
}

func (i *Instance) Update() error {
	x, y, w, h, err := i.winter.GetActiveWindowDimensions()
	if err == nil {
		i.rootX = x
		i.rootY = y
		i.targetWidth = w
		i.targetHeight = h
	}

	sw, sh := i.sitter.Size()

	if i.rootX != i.targetX || i.targetY != i.rootY-sh {
		i.moveTimer = 0
	}

	i.targetX = i.rootX
	i.targetY = i.rootY - sh

	x, y = ebiten.WindowPosition()
	if x != i.targetX || y != i.targetY {
		i.moveTimer += 4

		r := math.Atan2(float64(i.targetY-y), float64(i.targetX-x))
		x1 := math.Cos(r) * float64(i.moveTimer)
		y1 := math.Sin(r) * float64(i.moveTimer)
		if math.Abs(x1) > 60 || math.Abs(float64(i.targetX)-(float64(x)+x1)) < 20 {
			x = i.targetX
		} else {
			x += int(x1)
		}
		if math.Abs(y1) > 60 || math.Abs(float64(i.targetY)-(float64(y)+y1)) < 20 {
			y = i.targetY
		} else {
			y += int(y1)
		}

		ebiten.SetWindowPosition(x, y)
	}

	ww, wh := ebiten.WindowSize()
	if ww != sw || wh != sh {
		ebiten.SetWindowSize(sw, sh)
	}

	i.sitter.Tick()

	return nil
}

func (i *Instance) Draw(screen *ebiten.Image) {
	i.sitter.Draw(screen)
}

func (i *Instance) Layout(ow, oh int) (int, int) {
	return ow, oh
}
