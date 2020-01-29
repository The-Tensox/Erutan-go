package main

import (
	"math/rand"
	"sync"

	"github.com/golang/protobuf/ptypes"
	erutan "github.com/user/erutan_two/protos/realtime"
)

var (
	State     map[string]*erutan.NetObject = make(map[string]*erutan.NetObject)
	StatesMtx sync.RWMutex
)

// Start start handling gameplay
func Start() {
	StatesMtx.Lock()
	// Spawn ground
	ground := erutan.NetObject{
		ObjectId: RandomString(),
		OwnerId:  "server",
		Position: &erutan.NetVector3{X: 0, Y: 0, Z: 0},
		Rotation: &erutan.NetQuaternion{X: 0, Y: 0, Z: 0, W: 0},
		Scale:    &erutan.NetVector3{X: 100, Y: 1, Z: 100},
		Type:     erutan.NetObject_GROUND,
	}
	State[ground.ObjectId] = &ground

	// Spawn food
	food := erutan.NetObject{ // TODO: maybe general function to spawn object
		ObjectId: RandomString(),
		OwnerId:  "server",
		Position: &erutan.NetVector3{X: rand.Float64() * 10, Y: 1, Z: rand.Float64() * 10},
		Rotation: &erutan.NetQuaternion{X: 0, Y: 0, Z: 0, W: 0},
		Scale:    &erutan.NetVector3{X: 1, Y: 1, Z: 1},
		Type:     erutan.NetObject_FOOD,
	}
	State[food.ObjectId] = &food
	StatesMtx.Unlock()
	Update(update)

	// Spawn animals

	for i := 0; i < 1; i++ {
		StatesMtx.Lock()
		var food erutan.NetObject
		for _, element := range State {
			if element.Type == erutan.NetObject_FOOD {
				food = *element
			}
		}
		a := NewAnimal(food, erutan.NetVector3{X: rand.Float64() * 50, Y: 1, Z: rand.Float64() * 50})
		State[a.Object.ObjectId] = &a.Object
		StatesMtx.Unlock()
		Broadcast <- erutan.Packet{
			Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
			Type: &erutan.Packet_CreateObject{
				CreateObject: &erutan.Packet_CreateObjectPacket{
					Object: &a.Object,
				},
			},
		}
		a.Init()
	}
}

// WorldState return the current world state
func WorldState() []*erutan.Packet {
	var packets []*erutan.Packet
	StatesMtx.RLock()

	for _, element := range State {
		packets = append(packets, &erutan.Packet{
			Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
			Type: &erutan.Packet_CreateObject{
				CreateObject: &erutan.Packet_CreateObjectPacket{
					Object: element,
				},
			},
		})
	}

	StatesMtx.RUnlock()
	return packets
}

func update() {

}
