package sitters

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sitter struct {
	States   map[string]SitterState `yaml: "states"`
	state    string
	frame    int
	lifetime int
}

func (s *Sitter) Init() (err error) {
	err = s.SetState("idle")

	return
}

func (s *Sitter) SetState(state string) error {
	_, ok := s.States[state]
	if !ok {
		return fmt.Errorf("no such state %s", state)
	}
	s.state = state

	return nil
}

func (s *Sitter) Tick() {
	s.lifetime++

	if s.lifetime >= s.States[s.state].Rate {
		s.lifetime = 0
		s.IterateFrame()
	}
}

func (s *Sitter) IterateFrame() {
	s.frame++
	if s.frame >= len(s.States[s.state].Images.Images) {
		s.frame = 0
	}
}

func (s *Sitter) Draw(screen *ebiten.Image) {
	op := ebiten.DrawImageOptions{}
	screen.DrawImage(s.States[s.state].Images.Images[s.frame], &op)
}

func (s *Sitter) Size() (width, height int) {
	b := s.States[s.state].Images.Images[s.frame].Bounds()
	return b.Dx(), b.Dy()
}

type SitterState struct {
	Rate   int          `yaml: "rate"`
	Images SitterImages `yaml: "images"`
}

type SitterImages struct {
	Images     []*ebiten.Image
	imageFiles []string
}

func (s *SitterImages) UnmarshalYAML(unmarshal func(interface{}) error) error {
	err := unmarshal(&s.imageFiles)
	if err != nil {
		return err
	}
	return nil
}
