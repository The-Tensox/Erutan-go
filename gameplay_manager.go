package main

import (
	"math/rand"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	erutan "github.com/user/erutan_two/protos/realtime"
)

var (
	state     map[string]erutan.NetObject = make(map[string]erutan.NetObject)
	statesMtx sync.RWMutex
)

// Start start handling gameplay
func Start() {
	// Spawn ground
	ground := &erutan.NetObject{
		ObjectId: RandomString(),
		OwnerId:  "server",
		Position: &erutan.NetVector3{X: 0, Y: 0, Z: 0},
		Rotation: &erutan.NetQuaternion{X: 0, Y: 0, Z: 0, W: 0},
		Scale:    &erutan.NetVector3{X: 100, Y: 1, Z: 100},
		Type:     erutan.NetObject_GROUND,
	}
	state[ground.ObjectId] = *ground

	// Spawn food
	food := &erutan.NetObject{ // TODO: maybe general function to spawn object
		ObjectId: RandomString(),
		OwnerId:  "server",
		Position: &erutan.NetVector3{X: rand.Float32() * 10, Y: 1, Z: rand.Float32() * 10},
		Rotation: &erutan.NetQuaternion{X: 0, Y: 0, Z: 0, W: 0},
		Scale:    &erutan.NetVector3{X: 1, Y: 1, Z: 1},
		Type:     erutan.NetObject_FOOD,
	}
	state[food.ObjectId] = *food

	for range time.Tick(20000 * time.Millisecond) {
		update()
	}
}

// SendWorldState will broadcast the current world state
func SendWorldState() {
	statesMtx.RLock()
	for _, element := range state {
		DebugLogf("Object: %v", element)
		Broadcast <- erutan.Packet{
			Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
			Type: &erutan.Packet_CreateObject{
				CreateObject: &erutan.Packet_CreateObjectPacket{
					Object: &element,
				},
			},
		}
	}
	statesMtx.RUnlock()
}

func update() {
	// Spawn animals
	statesMtx.RLock()
	DebugLogf("Current state: %v", state)
	statesMtx.RUnlock()
}
