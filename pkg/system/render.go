package system

import (
	assetstore "bonecollector/pkg/asset-store"
	"bonecollector/pkg/entity"

	"github.com/veandco/go-sdl2/sdl"
)

type Render struct {
	componentSigature string // ??
	entities          []*entity.Entity
}

func (rs Render) Update(renderer *sdl.Renderer, camera sdl.Rect, assetStore *assetstore.AssetStore) {
	r := sdl.Rect{
		X: 20,
		Y: 20,
		W: 20,
		H: 20,
	}
	_ = renderer.SetDrawColor(255, 21, 21, 255)
	renderer.FillRect(&r)
}
