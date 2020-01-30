package main

import (
	erutan "github.com/user/erutan_two/protos/realtime"
)

type food struct {
	Object erutan.NetObject
}

// NewFood instanciate a food
func NewFood(position erutan.NetVector3) *food {
	return &food{
		Object: erutan.NetObject{
			ObjectId:   RandomString(),
			OwnerId:    "server",
			Position:   &position,
			Rotation:   &erutan.NetQuaternion{X: 0, Y: 0, Z: 0, W: 0},
			Scale:      &erutan.NetVector3{X: 1, Y: 1, Z: 1},
			Type:       erutan.NetObject_FOOD,
			Components: []*erutan.Component{},
		},
	}
}

func (f *food) Init() {
	Update(func(timeDelta int64) {
	})
}

func (f *food) GetObject() *erutan.NetObject { return &f.Object }

func (f *food) OnCollisionEnter(collisionedObjectID string) {
	DebugLogf("I %v got collisioned with %v", f.Object.ObjectId, collisionedObjectID)
}
