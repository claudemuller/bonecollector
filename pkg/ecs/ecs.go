package ecs

import (
	assetstore "bonecollector/pkg/asset-store"
	"bonecollector/pkg/component"
	"bonecollector/pkg/entity"
	"bonecollector/pkg/system"
	"errors"
	"log/slog"

	"github.com/veandco/go-sdl2/sdl"
)

type EntityManager struct {
	entityByTag               map[string]*entity.Entity
	tagByEntity               map[int]string
	entityByGroup             map[string]*entity.Entity
	groupByEntity             map[int]string
	components                map[int]interface{}
	entityComponentSignatures map[int]int
	systems                   map[int]interface{}
}

func New() EntityManager {
	slog.Info("creating entity manager")
	return EntityManager{
		entityByTag:               make(map[string]*entity.Entity),
		tagByEntity:               make(map[int]string),
		entityByGroup:             make(map[string]*entity.Entity),
		groupByEntity:             make(map[int]string),
		components:                make(map[int]interface{}),
		entityComponentSignatures: make(map[int]int),
		systems:                   make(map[int]interface{}),
	}
}

func (em *EntityManager) AddSystem(sysType int) {
	switch sysType {
	case system.MovementSystem:
		em.systems[sysType] = &system.Movement{}

	case system.RenderSystem:
		em.systems[sysType] = &system.Render{}
	}
}

func (em *EntityManager) RemoveSystem() {

}

func (em *EntityManager) HasSystem() {
}

func (em *EntityManager) GetSystem(sysType int) (interface{}, error) {
	sys, ok := em.systems[sysType]
	if !ok {
		return nil, errors.New("system not found")
	}
	return sys, nil
}

func (em *EntityManager) AddEntityToSystems(ent *entity.Entity) {
	for sig, sys := range em.systems {
		// TODO: test whether entity should be added to this system based off of its component sigs
		_ = sig
		system.AddEntityToSystem(sys, ent)
	}
}

func (em *EntityManager) AddComponent() {

}
func (em *EntityManager) GetComponent() {

}

func (em *EntityManager) Update(delta uint64) {
	sys, err := em.GetSystem(system.MovementSystem)
	if err != nil {
		// TODO: error
	} else {
		sys.(*system.Movement).Update(delta)
	}
}

func (em *EntityManager) Render(renderer *sdl.Renderer, camera sdl.Rect, assetStore *assetstore.AssetStore) {
	sys, err := em.GetSystem(system.RenderSystem)
	if err != nil {
		// TODO: error
	} else {
		sys.(*system.Render).Update(renderer, camera, assetStore)
	}

}

func (em *EntityManager) AddEntityTransformComponent(ent *entity.Entity, x, y int) {
	comp := component.NewTransform(x, y)
	em.components[comp.ID] = &comp
	em.entityComponentSignatures[ent.ID] = comp.ID
}

func (em *EntityManager) TagEntity(ent *entity.Entity, tag string) {
	em.entityByTag[tag] = ent
	em.tagByEntity[ent.ID] = tag
}
func (em *EntityManager) GroupEntity(ent *entity.Entity, group string) {
	em.entityByGroup[group] = ent
	em.groupByEntity[ent.ID] = group
}
