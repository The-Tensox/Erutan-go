package main

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	erutan "github.com/user/erutan_two/protos/realtime"
)

type animal struct {
	Object           *erutan.NetObject
	Target           erutan.NetVector3
	previousPosition erutan.NetVector3
}

// NewAnimal instanciate an animal
func NewAnimal(target erutan.NetVector3, position erutan.NetVector3) *animal {
	return &animal{
		Object: &erutan.NetObject{
			ObjectId: RandomString(),
			OwnerId:  "server",
			Position: &position,
			Rotation: &erutan.NetQuaternion{X: 0, Y: 0, Z: 0, W: 0},
			Scale:    &erutan.NetVector3{X: 1, Y: 1, Z: 1},
			Type:     erutan.NetObject_ANIMAL,
			Components: []*erutan.Component{&erutan.Component{
				Type: &erutan.Component_Animal{Animal: &erutan.Component_AnimalComponent{
					Life: 20,
				}}},
			},
		},
		Target:           target,
		previousPosition: position,
	}
}

func (a *animal) Init() {
	Update(func(timeDelta int64) bool {
		//StatesMtx.Lock()
		/*
			r := LookAtTwo(*a.Object.Position, *a.Target.Position)[3]
			yaw, pitch, roll := r[0], r[1], r[2]
			finalRotation := ToQuaternion(yaw, pitch, roll)
			State[a.Object.ObjectId].Rotation = &finalRotation
		*/

		distance := Distance(*a.Object.Position, a.Target)
		*a.Object.Position = Add(*a.Object.Position, Div(Sub(a.Target, *a.Object.Position) /*float64(timeDelta) */, distance*10))
		//State[a.Object.ObjectId].GetObject().Position = &position
		//DebugLogf("yep %v %v", a.Object.Position, a.Target)
		//StatesMtx.Unlock()

		// Let's not spam collisions check !
		//if Distance(position, a.previousPosition) > 3 { // TODO: tweak the threshold
		//a.previousPosition = *a.GetObject().Position
		//a.Object.Position = &position
		var l float64

		for _, element := range a.Object.Components {
			switch c := element.Type.(type) {
			case *erutan.Component_Animal:
				c.Animal.Life -= 0.01
				l = c.Animal.Life
			}
		}

		StateUpdate <- a
		//}

		/*
			Broadcast <- erutan.Packet{
				Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
				Type: &erutan.Packet_UpdateRotation{
					UpdateRotation: &erutan.Packet_UpdateRotationPacket{
						ObjectId: a.Object.ObjectId,
						Rotation: &finalRotation,
					},
				},
			}
		*/
		Broadcast <- erutan.Packet{
			Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
			Type: &erutan.Packet_UpdatePosition{
				UpdatePosition: &erutan.Packet_UpdatePositionPacket{
					ObjectId: a.Object.ObjectId,
					Position: a.Object.Position,
				},
			},
		}
		Broadcast <- erutan.Packet{
			Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
			Type: &erutan.Packet_UpdateAnimal{
				UpdateAnimal: &erutan.Packet_UpdateAnimalPacket{
					ObjectId: a.Object.ObjectId,
					Life:     l,
				},
			},
		}

		if l <= 0 {
			return true // Dead
		}
		return false
	})
}

func (a *animal) GetObject() *erutan.NetObject { return a.Object }

func (a *animal) OnCollisionEnter(other Collider) {
	// If we collided with food ++ life
	if _, ok := other.(*food); ok {

		var l float64
		for _, element := range a.GetObject().Components {
			if _, ok := element.Type.(*erutan.Component_Animal); ok {
				element.GetAnimal().Life += 20
				l = element.GetAnimal().Life
			}
		}

		Broadcast <- erutan.Packet{
			Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
			Type: &erutan.Packet_UpdateAnimal{
				UpdateAnimal: &erutan.Packet_UpdateAnimalPacket{
					ObjectId: a.GetObject().ObjectId,
					Life:     l,
				},
			},
		}
	}
}

// NotifyCallback implements Observer
func (a *animal) NotifyCallback(event Event) {
	switch event.eventID {
	case FoodMoved:
		//DebugLogf("food moved to %v", event.value)
		a.Target = event.value.(erutan.NetVector3)
	default:
		ServerLogf(time.Now(), "Unknown event type occured %v", event.eventID)
	}
}
