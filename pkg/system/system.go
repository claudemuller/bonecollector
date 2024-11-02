package system

import "bonecollector/pkg/entity"

const (
	MovementSystem = iota
	RenderSystem
)

func AddEntityToSystem(sys any, ent *entity.Entity) {
	switch s := sys.(type) {
	case Movement:
		s.entities = append(s.entities, ent)
	case Render:
	}
}
