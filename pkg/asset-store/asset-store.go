package assetstore

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type AssetStore struct {
	textures map[string]*sdl.Texture
	fonts    map[string]*ttf.Font
}

func New() *AssetStore {
	return &AssetStore{
		textures: make(map[string]*sdl.Texture),
		fonts:    make(map[string]*ttf.Font),
	}
}

func (as *AssetStore) AddTexture(renderer *sdl.Renderer, assetID, filename string) error {
	surface, err := img.Load(filename)
	if err != nil {
		return fmt.Errorf("error adding texture to asset store: %w", err)
	}

	tex, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return fmt.Errorf("error adding texture to asset store: %w", err)
	}
	surface.Free()

	as.textures[assetID] = tex

	return nil
}

func (as *AssetStore) ClearAssets() {
	for _, t := range as.textures {
		if err := t.Destroy(); err != nil {
			// TODO: return a bunch of possible errors to caller?
			slog.Error("error destroying texture")
		}
	}
	as.textures = make(map[string]*sdl.Texture)

	for _, f := range as.fonts {
		f.Close()
	}
	as.fonts = make(map[string]*ttf.Font)

}

func (as *AssetStore) GetTexture(assetID string) (*sdl.Texture, error) {
	tex, ok := as.textures[assetID]
	if !ok {
		return nil, errors.New("texture not found")
	}
	return tex, nil
}

func (as *AssetStore) AddFont(assetID, path string, size int) error {
	font, err := ttf.OpenFont(path, size)
	if err != nil {
		return fmt.Errorf("error opening font: %w", err)
	}

	as.fonts[assetID] = font

	return nil
}

func (as *AssetStore) GetFont(assetID string) (*ttf.Font, error) {
	font, ok := as.fonts[assetID]
	if !ok {
		return nil, errors.New("font not found")
	}
	return font, nil
}
