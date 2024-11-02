package component

const (
	TransformCompenent = iota
)

type Vec2 struct {
	X int
	Y int
}

type Transform struct {
	ID       int
	Pos      Vec2
	Scale    Vec2
	Rotation float32
}

func NewTransform(x, y int) Transform {
	return Transform{
		ID: 1,
		Pos: Vec2{
			X: x,
			Y: y,
		},
	}
}
