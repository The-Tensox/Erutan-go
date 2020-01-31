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
	State map[string]Collider = make(map[string]Collider)

	// StatesMtx is used to ensure safe concurrency on the State map
	StatesMtx sync.RWMutex

	// StateUpdate is a global event to notify gameplay manager
	// that a state update has occured
	StateUpdate chan Collider = make(chan Collider, 1000)

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
	for i := 0; i < 5; i++ {
		a := NewAnimal(*f.Object.Position, erutan.NetVector3{X: rand.Float64() * 50, Y: 1, Z: rand.Float64() * 50})
		EventDispatcher.Add(a)
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
					Object: element.GetObject(),
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
			State[s.GetObject().ObjectId] = s

			/*
				switch v := State[movedObjectID].(type) {
				case *food:
				}
			*/
			checkCollisions(s)
			StatesMtx.Unlock()
		}
	}
}

func checkCollisions(collider Collider) {
	for _, element := range State {
		if element.GetObject().ObjectId == collider.GetObject().ObjectId {
			continue
		}
		//DebugLogf("Distance %v", Distance(*a.GetObject().Position, *element.GetObject().Position))
		if Distance(*collider.GetObject().Position, *element.GetObject().Position) < 5 {
			//DebugLogf("Collision a: %v, b: %v", a.GetObject().ObjectId, element.GetObject().ObjectId)
			// Notify both objects of a collision !
			// Eventually run it in goroutine?
			collider.OnCollisionEnter(element)
			element.OnCollisionEnter(collider)
		}
	}
}
