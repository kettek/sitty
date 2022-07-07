package sitters

import (
	"fmt"
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

// Sitter is our sitter.
type Sitter struct {
	States            map[string]SitterState `yaml:"states"`
	InterpreterSource string                 `yaml:"interpreter"`
	state             string
	frame             int
	lifetime          int
	interp            *interp.Interpreter
	Click             func(s *Sitter) bool
}

func (s *Sitter) Init() (err error) {
	exports := make(interp.Exports)
	exports["sitters/sitters"] = map[string]reflect.Value{
		"Sitter":      reflect.ValueOf((*Sitter)(nil)),
		"SitterState": reflect.ValueOf((*SitterState)(nil)),
	}

	// Build our state interpreters.
	for k, state := range s.States {
		if state.InterpreterSource != "" {
			state.interp = interp.New(interp.Options{})
			state.interp.Use(stdlib.Symbols)

			if err := state.interp.Use(exports); err != nil {
				fmt.Println(err)
			}

			_, err = state.interp.Eval(fmt.Sprintf(`
		import (
			"fmt"
			"sitters"
		)

		%s
	`, state.InterpreterSource))

			if val, err := state.interp.Eval("Click"); err == nil {
				state.Click = val.Interface().(func(s *Sitter) bool)
			} else {
				fmt.Println(err)
			}
			s.States[k] = state
		}
	}

	err = s.SetState("idle")

	s.interp = interp.New(interp.Options{})
	s.interp.Use(stdlib.Symbols)

	if err := s.interp.Use(exports); err != nil {
		fmt.Println(err)
	}

	_, err = s.interp.Eval(fmt.Sprintf(`
		import (
			"fmt"
			"sitters"
		)

		%s
	`, s.InterpreterSource))

	if val, err := s.interp.Eval("Click"); err == nil {
		s.Click = val.Interface().(func(s *Sitter) bool)
	} else {
		fmt.Println(err)
	}

	return
}

func (s *Sitter) SetState(state string) error {
	_, ok := s.States[state]
	if !ok {
		return fmt.Errorf("no such state %s", state)
	}
	s.state = state
	s.frame = 0

	return nil
}

func (s *Sitter) Tick() {
	s.lifetime++

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		if s.States[s.state].Click != nil {
			if !s.States[s.state].Click(s) {
				s.Click(s)
			}
		} else {
			s.Click(s)
		}
	}

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
	Rate              int          `yaml: "rate"`
	Images            SitterImages `yaml: "images"`
	InterpreterSource string       `yaml:"interpreter"`
	interp            *interp.Interpreter
	Click             func(s *Sitter) bool
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
