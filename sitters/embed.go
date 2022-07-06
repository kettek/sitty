package sitters

import (
	"embed"
	"image"
	_ "image/png"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/go-multipath/v2"
	"gopkg.in/yaml.v3"
)

//go:embed *
var embedFS embed.FS

var FS multipath.FS

func init() {
	FS.AddFS(embedFS)
	FS.AddFS(os.DirFS("sitters"))
}

func LoadSitter(s string) (sitter *Sitter, err error) {
	bytes, err := FS.ReadFile(filepath.Join(s, "sitter.yml"))
	if err != nil {
		return nil, err
	}
	sitter = &Sitter{}

	if err := yaml.Unmarshal(bytes, sitter); err != nil {
		return nil, err
	}

	for k, state := range sitter.States {
		for _, imgFile := range state.Images.imageFiles {
			imgFiles, _ := FS.Glob(filepath.Join(s, imgFile))
			for _, imgFile := range imgFiles {
				f, err := FS.Open(imgFile)
				if err != nil {
					continue
				}
				img, _, err := image.Decode(f)
				f.Close()
				if err != nil {
					continue
				}
				eimg := ebiten.NewImageFromImage(img)
				state.Images.Images = append(state.Images.Images, eimg)
			}
		}
		sitter.States[k] = state
	}

	return sitter, nil
}
