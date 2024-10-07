package assetstore

import (
	"fmt"
	"slices"
	"strings"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type Texture struct {
	texture *sdl.Texture
	id      string
}

type AssetStore struct {
	textures []Texture
}

func (as *AssetStore) AddTexture(renderer *sdl.Renderer, assetID, filename string) error {
	surface, err := img.Load(filename)
	if err != nil {
		return fmt.Errorf("error adding texture to asset store", err)
	}

	tex, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return fmt.Errorf("error adding texture to asset store", err)
	}
	surface.Free()

	as.textures = append(as.textures, Texture{
		texture: tex,
		id:      assetID,
	})

	return nil
}

func (as *AssetStore) GetTexture(assetID string) *Texture {
	i, found := slices.BinarySearchFunc(as.textures, Texture{id: assetID}, func(a, b Texture) int {
		return strings.Compare(a.id, b.id)
	})
	if found {
		return &as.textures[i]
	}
	return nil
}
