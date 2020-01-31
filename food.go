package main

import (
	"math/rand"

	"github.com/golang/protobuf/ptypes"
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
}

func (f *food) GetObject() *erutan.NetObject { return &f.Object }

func (f *food) OnCollisionEnter(other Collider) {
	// If we collided with animal
	if _, ok := other.(*animal); ok {
		f.GetObject().Position = &erutan.NetVector3{X: rand.Float64() * 50, Y: 1, Z: rand.Float64() * 50}
		EventDispatcher.Notify(Event{eventID: FoodMoved, value: *f.GetObject().Position})

		StateUpdate <- f
		Broadcast <- erutan.Packet{
			Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
			Type: &erutan.Packet_UpdatePosition{
				UpdatePosition: &erutan.Packet_UpdatePositionPacket{
					ObjectId: f.Object.ObjectId,
					Position: f.GetObject().Position,
				},
			},
		}
		/*
			Broadcast <- erutan.Packet{
				Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
				Type: &erutan.Packet_FoodEaten{
					FoodEaten: &erutan.Packet_FoodEatenPacket{
						FoodId: f.Object.ObjectId,
						Eater:  other.GetObject(),
					},
				},
			}
		*/
	}
}
