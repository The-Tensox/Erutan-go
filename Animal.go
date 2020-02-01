package main

import (
	"github.com/golang/protobuf/ptypes"
	erutan "github.com/user/erutan_two/protos/realtime"
)

type Animal struct {
	*AbstractBehaviour
	Target           string
	previousPosition erutan.NetVector3
}

// NewAnimal instanciate an animal
func NewAnimal(target string, position erutan.NetVector3) *Animal {
	b := &AbstractBehaviour{
		Object: erutan.NetObject{
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
	}
	a := &Animal{
		AbstractBehaviour: b,
		Target:            target,
		previousPosition:  position,
	}
	a.Behaviour = a
	return a
}

// Start is used to initialize your object
func (a *Animal) Start() {
	a.Update()
}

// Update is used to handle this object life loop
func (a *Animal) Update() {
	Update(func(deltaTime int64) bool {
		StatesMtx.RLock()
		target := State[a.Target].Object.Position
		StatesMtx.RUnlock()
		//StatesMtx.Lock()
		/*
			r := LookAtTwo(*a.Object.Position, *a.Target.Position)[3]
			yaw, pitch, roll := r[0], r[1], r[2]
			finalRotation := ToQuaternion(yaw, pitch, roll)
			State[a.Object.ObjectId].Rotation = &finalRotation
		*/

		distance := Distance(*a.AbstractBehaviour.Object.Position, *target)
		*a.AbstractBehaviour.Object.Position = Add(*a.AbstractBehaviour.Object.Position,
			Div(Sub(*target, *a.AbstractBehaviour.Object.Position) /*float64(timeDelta) */, distance*10))
		//State[a.Object.ObjectId].GetObject().Position = &position
		//DebugLogf("yep %v %v", a.Object.Position, a.Target)
		//StatesMtx.Unlock()

		// Let's not spam collisions check !
		//if Distance(position, a.previousPosition) > 3 { // TODO: tweak the threshold
		//a.previousPosition = *a.GetObject().Position
		//a.Object.Position = &position
		var l float64

		for _, element := range a.AbstractBehaviour.Object.Components {
			switch c := element.Type.(type) {
			case *erutan.Component_Animal:
				c.Animal.Life -= 0.01
				l = c.Animal.Life
			}
		}

		StateUpdate <- a.AbstractBehaviour
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
					ObjectId: a.AbstractBehaviour.Object.ObjectId,
					Position: a.AbstractBehaviour.Object.Position,
				},
			},
		}
		Broadcast <- erutan.Packet{
			Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
			Type: &erutan.Packet_UpdateAnimal{
				UpdateAnimal: &erutan.Packet_UpdateAnimalPacket{
					ObjectId: a.AbstractBehaviour.Object.ObjectId,
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

func (a *Animal) OnCollisionEnter(other erutan.NetObject) {
	// If we collided with food ++ life
	if other.Type == erutan.NetObject_FOOD {

		var l float64
		for _, element := range a.AbstractBehaviour.Object.Components {
			if _, ok := element.Type.(*erutan.Component_Animal); ok {
				element.GetAnimal().Life += 20
				l = element.GetAnimal().Life
			}
		}

		Broadcast <- erutan.Packet{
			Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
			Type: &erutan.Packet_UpdateAnimal{
				UpdateAnimal: &erutan.Packet_UpdateAnimalPacket{
					ObjectId: a.AbstractBehaviour.Object.ObjectId,
					Life:     l,
				},
			},
		}
	}
}

/*
// NotifyCallback implements Observer
func (a *Animal) NotifyCallback(event Event) {
	switch event.eventID {
	case FoodMoved:
		DebugLogf("food moved to %v", event.value)
		a.Target = event.value.(erutan.NetVector3)
	default:
		ServerLogf(time.Now(), "Unknown event type occured %v", event.eventID)
	}
}
*/

// OnDestroy is called before getting destroyed
func (a *Animal) OnDestroy() {}