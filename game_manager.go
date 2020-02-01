package main

import (
	"math/rand"
	"sync"

	"github.com/golang/protobuf/ptypes"

	erutan "github.com/user/erutan_two/protos/realtime"
)

var (
	// EventDispatcher is a global event dispatcher (Observable)
	EventDispatcher Watch = Watch{}

	// State is used to store the world state
	State map[string]*AbstractBehaviour = make(map[string]*AbstractBehaviour)

	// StatesMtx is used to ensure safe concurrency on the State map
	StatesMtx sync.RWMutex

	// StateUpdate is a global event to notify gameplay manager
	// that a state update has occured
	StateUpdate chan *AbstractBehaviour = make(chan *AbstractBehaviour, 1000)

	// Collision is used to offer a way to communicate collision between two ObjectId across NetObjects
	// map indexed by ObjectId, the collisioned ObjectId is pushed into the channel
	Collision map[string]chan string = make(map[string]chan string, 1000)

	// CollisionMtx is used to ensure safe concurrency on the Collision map
	CollisionMtx sync.RWMutex
)

// RunGame start handling gameplay
func RunGame() {
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
	State[ground.ObjectId] = &AbstractBehaviour{Object: ground}

	// Spawn food
	f := NewFood(erutan.NetVector3{X: rand.Float64() * 50, Y: 1, Z: rand.Float64() * 50})
	State[f.AbstractBehaviour.Object.ObjectId] = f.AbstractBehaviour
	StatesMtx.Unlock()
	f.Start()
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
		a := NewAnimal(f.AbstractBehaviour.Object.ObjectId, erutan.NetVector3{X: rand.Float64() * 50, Y: 1, Z: rand.Float64() * 50})
		//EventDispatcher.Add(a)
		State[a.AbstractBehaviour.Object.ObjectId] = a.AbstractBehaviour

		Broadcast <- erutan.Packet{
			Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
			Type: &erutan.Packet_CreateObject{
				CreateObject: &erutan.Packet_CreateObjectPacket{
					Object: &a.AbstractBehaviour.Object,
				},
			},
		}
		a.Start()
	}
	StatesMtx.Unlock()
	go handleStateUpdates()
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
					Object: &element.Object,
				},
			},
		})
	}

	StatesMtx.RUnlock()
	return packets
}

func handleStateUpdates() {
	for {
		select {
		case s := <-StateUpdate:
			StatesMtx.Lock()
			//DebugLogf("%v moved to %v", s.GetObject().ObjectId, s.GetObject().Position)
			State[s.Object.ObjectId] = s
			checkCollisions(*s)
			StatesMtx.Unlock()
		}
	}
}

func checkCollisions(collider AbstractBehaviour) {
	for _, element := range State {
		if element.Object.ObjectId == collider.Object.ObjectId {
			continue
		}
		//DebugLogf("Distance %v", Distance(*collider.Object.Position, *element.Object.Position))
		if Distance(*collider.Object.Position, *element.Object.Position) < 5 {
			//DebugLogf("Collision a: %v, b: %v", collider.Object, element.Object)
			// Notify both objects of a collision !
			// Eventually run it in goroutine?
			collider.Behaviour.OnCollisionEnter(element.Object)
			element.Behaviour.OnCollisionEnter(collider.Object)
		}
	}
}
