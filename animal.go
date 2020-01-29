package main

import (
	"github.com/golang/protobuf/ptypes"
	erutan "github.com/user/erutan_two/protos/realtime"
)

type animal struct {
	Object erutan.NetObject
}

// NewAnimal instanciate an animal
func NewAnimal(food erutan.NetObject, position erutan.NetVector3) *animal {
	return &animal{
		Object: erutan.NetObject{
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
	}
}

func (a *animal) Init() {
	Update(func() {
		StatesMtx.Lock()
		//angleWithFood := Cos(a.Object.Position.X)
		var animalComponent erutan.Component_Animal
		for _, c := range a.Object.Components {
			switch u := c.Type.(type) {
			case *erutan.Component_Animal:
				animalComponent = *u
			}
		}
		rotation := LookAt(*a.Object.Position, *animalComponent.Animal.Food.Position)
		DebugLogf("Rotation %v", rotation)
		*State[a.Object.ObjectId].Rotation = rotation
		position := a.Object.Position
		position = &erutan.NetVector3{X: position.X, Y: position.Y, Z: position.Z + 0.1}
		*State[a.Object.ObjectId].Position = *position
		StatesMtx.Unlock()
		Broadcast <- erutan.Packet{
			Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
			Type: &erutan.Packet_UpdateRotation{
				UpdateRotation: &erutan.Packet_UpdateRotationPacket{
					ObjectId: a.Object.ObjectId,
					Rotation: &rotation,
				},
			},
		}
		Broadcast <- erutan.Packet{
			Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
			Type: &erutan.Packet_UpdatePosition{
				UpdatePosition: &erutan.Packet_UpdatePositionPacket{
					ObjectId: a.Object.ObjectId,
					Position: position,
				},
			},
		}
	})
}
