package main

import (
	"math/rand"

	"github.com/golang/protobuf/ptypes"
	erutan "github.com/user/erutan_two/protos/realtime"
)

type Food struct {
	*AbstractBehaviour
}

// NewFood instanciate a food
func NewFood(position erutan.NetVector3) *Food {
	b := &AbstractBehaviour{
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
	f := &Food{b}
	f.Behaviour = f
	return f
}

// Start is used to initialize your object
func (f *Food) Start() {
	f.Update()
}

// Update is used to handle this object life loop
func (f *Food) Update() {
	Update(func(deltaTime int64) bool {
		return true
	})
}

func (f *Food) OnCollisionEnter(other erutan.NetObject) {
	// If we collided with animal
	if other.Type == erutan.NetObject_ANIMAL {
		f.AbstractBehaviour.Object.Position = &erutan.NetVector3{X: rand.Float64() * 50, Y: 1, Z: rand.Float64() * 50}
		//EventDispatcher.Notify(Event{eventID: FoodMoved, value: *f.AbstractBehaviour.Object.Position})
		DebugLogf("%v", f.AbstractBehaviour.Object.Position)
		StateUpdate <- f.AbstractBehaviour
		Broadcast <- erutan.Packet{
			Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
			Type: &erutan.Packet_UpdatePosition{
				UpdatePosition: &erutan.Packet_UpdatePositionPacket{
					ObjectId: f.AbstractBehaviour.Object.ObjectId,
					Position: f.AbstractBehaviour.Object.Position,
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

// OnDestroy is called before getting destroyed
func (f *Food) OnDestroy() {}