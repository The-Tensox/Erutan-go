package game

import (
	"math/rand"

	"github.com/golang/protobuf/ptypes"
	erutan "github.com/user/erutan/protos/realtime"
	"github.com/user/erutan/utils"
)

type Food struct {
	*AbstractBehaviour
}

// NewFood instanciate a food
func NewFood(position erutan.NetVector3) *Food {
	b := &AbstractBehaviour{
		Object: erutan.NetObject{
			ObjectId:   utils.RandomString(),
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
	utils.Update(func(deltaTime int64) bool {
		return true
	})
}

func (f *Food) OnCollisionEnter(other erutan.NetObject) {
	// If we collided with animal
	if other.Type == erutan.NetObject_ANIMAL {
		f.AbstractBehaviour.Object.Position = &erutan.NetVector3{X: rand.Float64() * 50, Y: 1, Z: rand.Float64() * 50}
		GameManagerInstance.StateUpdate <- f.AbstractBehaviour
		GameManagerInstance.Broadcast <- erutan.Packet{
			Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
			Type: &erutan.Packet_UpdatePosition{
				UpdatePosition: &erutan.Packet_UpdatePositionPacket{
					ObjectId: f.AbstractBehaviour.Object.ObjectId,
					Position: f.AbstractBehaviour.Object.Position,
				},
			},
		}
	}
}

// OnDestroy is called before getting destroyed
func (f *Food) OnDestroy() {}
