package main

import (
	"github.com/golang/protobuf/ptypes"
	erutan "github.com/user/erutan_two/protos/realtime"
)

type animal struct {
	Object           *erutan.NetObject
	Target           *erutan.NetObject
	previousPosition erutan.NetVector3
}

// NewAnimal instanciate an animal
func NewAnimal(food erutan.NetObject, position erutan.NetVector3) *animal {
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
					Life:   20,
					Food:   &food,
					Target: &food,
				}}},
			},
		},
		Target:           &food,
		previousPosition: position,
	}
}

func (a *animal) Init() {
	Update(func(timeDelta int64) {
		StatesMtx.Lock()
		/*
			r := LookAtTwo(*a.Object.Position, *a.Target.Position)[3]
			yaw, pitch, roll := r[0], r[1], r[2]
			finalRotation := ToQuaternion(yaw, pitch, roll)
			State[a.Object.ObjectId].Rotation = &finalRotation
		*/

		distance := Distance(*a.Object.Position, *a.Target.Position)
		position := Add(*a.Object.Position, Div(Sub(*a.Target.Position, *a.Object.Position) /*float64(timeDelta) */, distance*10))
		State[a.Object.ObjectId].GetObject().Position = &position
		//DebugLogf("yep %v", position)
		StatesMtx.Unlock()

		// Let's not spam collisions check !
		if Distance(position, a.previousPosition) > 3 { // TODO: tweak the threshold
			a.previousPosition = *State[a.Object.ObjectId].GetObject().Position
			Movement <- a.Object.ObjectId
		}

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
					Position: &position,
				},
			},
		}
	})
}

func (a *animal) GetObject() *erutan.NetObject { return a.Object }

func (a *animal) OnCollisionEnter(collisionedObjectID string) {
	DebugLogf("I %v got collisioned with %v", a.Object.ObjectId, collisionedObjectID)
}
