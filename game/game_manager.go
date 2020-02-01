package game

import (
	"math/rand"
	"sync"

	"github.com/golang/protobuf/ptypes"
	erutan "github.com/user/erutan/protos/realtime"
	utils "github.com/user/erutan/utils"
)

var once sync.Once

type GameManager struct {
	// Broadcast is a global channel to send packets to clients
	Broadcast chan erutan.Packet

	// State is used to store the world state
	State map[string]*AbstractBehaviour

	// StatesMtx is used to ensure safe concurrency on the State map
	StatesMtx sync.RWMutex

	// StateUpdate is a global event to notify gameplay manager
	// that a state update has occured
	StateUpdate chan *AbstractBehaviour

	// Collision is used to offer a way to communicate collision between two ObjectId across NetObjects
	// map indexed by ObjectId, the collisioned ObjectId is pushed into the channel
	Collision map[string]chan string

	// CollisionMtx is used to ensure safe concurrency on the Collision map
	CollisionMtx sync.RWMutex

	// Receive receive packets from clients
	Receive chan *erutan.Packet
}

var GameManagerInstance *GameManager

// InitializeGame returns a thread-safe singleton instance of the game manager
func InitializeGame() {
	once.Do(func() {
		GameManagerInstance =
			&GameManager{
				Broadcast:   make(chan erutan.Packet, 1000),
				State:       make(map[string]*AbstractBehaviour),
				StateUpdate: make(chan *AbstractBehaviour, 1000), // TODO: maybe shouldnt be buffered ...
				Collision:   make(map[string]chan string, 1000),
				Receive:     make(chan *erutan.Packet, 1000),
			}
	})
}

// RunGame start handling gameplay
func (g *GameManager) RunGame() {
	go g.handleClientPackets()
	g.StatesMtx.Lock()
	// Spawn ground
	ground := erutan.NetObject{
		ObjectId: utils.RandomString(),
		OwnerId:  "server",
		Position: &erutan.NetVector3{X: 0, Y: 0, Z: 0},
		Rotation: &erutan.NetQuaternion{X: 0, Y: 0, Z: 0, W: 0},
		Scale:    &erutan.NetVector3{X: 100, Y: 1, Z: 100},
		Type:     erutan.NetObject_GROUND,
	}
	g.State[ground.ObjectId] = &AbstractBehaviour{Object: ground}

	// Spawn food
	f := NewFood(erutan.NetVector3{X: rand.Float64() * 50, Y: 1, Z: rand.Float64() * 50})
	g.State[f.AbstractBehaviour.Object.ObjectId] = f.AbstractBehaviour
	g.StatesMtx.Unlock()
	f.Start()
	// Spawn animals
	g.StatesMtx.Lock()
	for i := 0; i < 1; i++ {
		a := NewAnimal(f.AbstractBehaviour.Object.ObjectId, erutan.NetVector3{X: rand.Float64() * 50, Y: 1, Z: rand.Float64() * 50})
		g.State[a.AbstractBehaviour.Object.ObjectId] = a.AbstractBehaviour

		g.Broadcast <- erutan.Packet{
			Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
			Type: &erutan.Packet_CreateObject{
				CreateObject: &erutan.Packet_CreateObjectPacket{
					Object: &a.AbstractBehaviour.Object,
				},
			},
		}
		a.Start()
	}
	g.StatesMtx.Unlock()
	go g.handleStateUpdates()
}

// WorldState return the current world state
func (g *GameManager) WorldState() []*erutan.Packet {
	var packets []*erutan.Packet
	g.StatesMtx.RLock()

	for _, element := range g.State {
		packets = append(packets, &erutan.Packet{
			Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
			Type: &erutan.Packet_CreateObject{
				CreateObject: &erutan.Packet_CreateObjectPacket{
					Object: &element.Object,
				},
			},
		})
	}

	g.StatesMtx.RUnlock()
	return packets
}

func (g *GameManager) handleClientPackets() {
	for {
		select {
		case p := <-g.Receive:
			// TODO: handle client messages here
			switch t := p.Type.(type) {
			case *erutan.Packet_UpdateParameters:
				for _, element := range t.UpdateParameters.Parameters {
					switch param := element.Type.(type) {
					case *erutan.Packet_UpdateParametersPacket_Parameter_Tickrate:
						utils.Config.TickRate = param.Tickrate
					}
				}
			case *erutan.Packet_CreateObject:
				var o *AbstractBehaviour
				g.StatesMtx.Lock()
				switch t.CreateObject.Object.Type {
				case erutan.NetObject_ANIMAL:
					// Find a target (first found, TODO: nearest, w/e)
					for _, e := range g.State {
						if f, ok := e.Behaviour.(*Food); ok {
							o = NewAnimal(f.Object.ObjectId, *t.CreateObject.Object.Position).AbstractBehaviour
						}
					}
				}
				g.State[o.Object.ObjectId] = o
				g.StatesMtx.Unlock()
				g.Broadcast <- erutan.Packet{
					Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
					Type: &erutan.Packet_CreateObject{
						CreateObject: &erutan.Packet_CreateObjectPacket{
							Object: &o.Object,
						},
					},
				}
				o.Start()

			default:
				utils.DebugLogf("Client sent unimplemented packet: %v", t)
			}
		}
	}
}

func (g *GameManager) handleStateUpdates() {
	for {
		select {
		case s := <-g.StateUpdate:
			g.StatesMtx.Lock()
			g.State[s.Object.ObjectId] = s
			g.checkCollisions(*s)
			g.StatesMtx.Unlock()
		}
	}
}

func (g *GameManager) checkCollisions(collider AbstractBehaviour) {
	for _, element := range g.State {
		if element.Object.ObjectId == collider.Object.ObjectId {
			continue
		}
		if utils.Distance(*collider.Object.Position, *element.Object.Position) < 1 {
			// Notify both objects of a collision !
			// Eventually run it in goroutine?
			collider.Behaviour.OnCollisionEnter(element.Object)
			element.Behaviour.OnCollisionEnter(collider.Object)
		}
	}
}
