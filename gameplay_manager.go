package main

import (
	"math/rand"
	"sync"

	"github.com/golang/protobuf/ptypes"
	erutan "github.com/user/erutan_two/protos/realtime"
)

var (
	// State is used to store the world state
	State map[string]Collider = make(map[string]Collider)

	// StatesMtx is used to ensure safe concurrency on the State map
	StatesMtx sync.RWMutex

	// Movement is a global event to notify gameplay manager
	// that a physic movement occured and a collision check has to be done
	Movement chan string = make(chan string, 1000)

	// Collision is used to offer a way to communicate collision between two ObjectId across NetObjects
	// map indexed by ObjectId, the collisioned ObjectId is pushed into the channel
	Collision map[string]chan string = make(map[string]chan string, 1000)

	// CollisionMtx is used to ensure safe concurrency on the Collision map
	CollisionMtx sync.RWMutex
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
	State[ground.ObjectId] = &ObjectBehaviour{Object: ground}

	// Spawn food
	f := NewFood(erutan.NetVector3{X: rand.Float64() * 50, Y: 1, Z: rand.Float64() * 50})
	State[f.Object.ObjectId] = &food{Object: f.Object}
	StatesMtx.Unlock()
	f.Init()
	Update(update)

	// Spawn animals
	StatesMtx.Lock()
	/*
		var food erutan.NetObject
		for _, element := range State {
			if element.Type == erutan.NetObject_FOOD {
				food = *element
			}
		}
	*/
	for i := 0; i < 1; i++ {
		a := NewAnimal(f.Object, erutan.NetVector3{X: rand.Float64() * 50, Y: 1, Z: rand.Float64() * 50})
		State[a.Object.ObjectId] = &animal{Object: a.Object}

		Broadcast <- erutan.Packet{
			Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
			Type: &erutan.Packet_CreateObject{
				CreateObject: &erutan.Packet_CreateObjectPacket{
					Object: a.Object,
				},
			},
		}
		a.Init()
	}
	StatesMtx.Unlock()
	go handleMovements()
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
					Object: element.GetObject(),
				},
			},
		})
	}

	StatesMtx.RUnlock()
	return packets
}

func update(deltaTime int64) {

}

func handleMovements() {
	for {
		select {
		case movedObjectID := <-Movement:
			StatesMtx.Lock()
			// DebugLogf("%v moved to %v", movedObjectID, State[movedObjectID].Position)
			checkCollisions(movedObjectID)
			StatesMtx.Unlock()
		}
	}
}

func checkCollisions(movedObjectID string) {
	a := State[movedObjectID]
	for _, element := range State {
		if element.GetObject().ObjectId == movedObjectID {
			continue
		}
		//DebugLogf("Distance %v", Distance(*a.GetObject().Position, *element.GetObject().Position))
		if Distance(*a.GetObject().Position, *element.GetObject().Position) < 5 {
			//DebugLogf("Collision a: %v, b: %v", a.GetObject().ObjectId, element.GetObject().ObjectId)
			// Notify both objects of a collision !
			a.OnCollisionEnter(element.GetObject().ObjectId)
			element.OnCollisionEnter(a.GetObject().ObjectId)
		}
	}
}
