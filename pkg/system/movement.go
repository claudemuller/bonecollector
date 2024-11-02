package system

import (
	"bonecollector/pkg/component"
	"bonecollector/pkg/entity"
)

type Movement struct {
	componentSigature string // ??
	entities          []*entity.Entity
}

func (rs *Movement) Update(delta uint64) {
	for _, e := range rs.entities {
		transform := e.GetComponent(component.TransformCompenent).(component.Transform)
		transform.Pos.X += 1
		transform.Pos.Y += 1
	}

	// Update
}
