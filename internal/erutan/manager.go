package erutan

import (
	"github.com/The-Tensox/Erutan-go/internal/cfg"
	"github.com/The-Tensox/Erutan-go/internal/ecs"
	"github.com/The-Tensox/Erutan-go/internal/log"
	"github.com/The-Tensox/Erutan-go/internal/obs"
	"github.com/The-Tensox/Erutan-go/internal/utils"
	"github.com/The-Tensox/octree"
	"github.com/The-Tensox/protometry"
	"github.com/golang/protobuf/ptypes"
	"go.uber.org/zap"
	"math"
	"sync"
)

var (
	once sync.Once

	// ManagerInstance is a global singleton erutan manager
	ManagerInstance *Manager
)

// Manager ...
type Manager struct {
	// World is the structure that handle all Systems in the Object Component System design
	ecs.World

	// ClientsIn handle incoming messages, doesn't need to be a map since when we receive we know it's a connected client
	ClientsIn chan ClientPacket
	//ClientsMu sync.RWMutex
	// ClientsOut is a map that send packets to clients (map[string]chan Packet)
	ClientsOut sync.Map
	// Map a client to settings (map[string]ClientSettings)
	ClientsSettings sync.Map

	// BroadcastOut is a channel to send packets to every clients
	BroadcastOut chan Packet

	obs.Watch

	// Reference to the network octree, needed to iterate networked objects
	networkSystem *NetworkSystem // TODO: imho prob want to search into collision system, faster search based on position ?
}

// Initialize returns a thread-safe singleton instance of the erutan manager
func Initialize() {
	once.Do(func() {
		ManagerInstance =
			&Manager{
				World:        ecs.World{},
				ClientsIn:    make(chan ClientPacket, 100000),
				BroadcastOut: make(chan Packet, 1000),
				Watch:        *obs.NewWatch(),
			}
	})
}

// Run start handling gameplay
func (m *Manager) Run() {
	h := NewHerbivorousSystem()
	e := NewEatableSystem()
	c := NewCollisionSystem()
	n := NewNetworkSystem(utils.GetProtoTime())
	m.networkSystem = n
	m.AddSystem(h)
	m.AddSystem(e)
	m.AddSystem(c)
	m.AddSystem(n)

	m.Register(h)
	m.Register(e)
	m.Register(c)
	m.Register(n)

	//gs := cfg.Get().Logic.GroundSize
	//p := perlin.NewPerlin(1, 1, 1, 1337)
	//for x := 0.; x < gs; x++ {
	//	for y := 0.; y < gs; y++ {
	//		noise := p.Noise2D(x/10, y/10)
	//		//log.Zap.Info("Noise at %.1f;%.1f: %v", x, y, noise)
	//		m.AddGround(protometry.NewVector3(x, float64(int(10*noise))-5, y), 1)
	//	}
	//}
	//ds := DiamondSquareAlgorithm(int(math.Pow(2, 6))+1, 10, 1)
	//sideLength := 1.
	//for x := 0.; int(x) < len(ds); x++ {
	//	for y := 0.; int(y) < len(ds[int(x)]); y++ {
	//		m.AddGround(protometry.NewVector3(x*sideLength,
	//			float64(int(sideLength*ds[int(x)][int(y)])),
	//			y*sideLength),
	//			sideLength)
	//	}
	//}

	for i := 0; i < cfg.Get().Logic.InitialHerbs; i++ {
		p := protometry.RandomCirclePoint(0, 5, 0, 50)
		m.AddHerb(&p)
	}

	for i := 0; i < cfg.Get().Logic.Herbivorous.Quantity; i++ {
		p := protometry.RandomCirclePoint(0, 5, 0, 50)
		m.AddHerbivorous(&p, protometry.NewVector3(1, 1, 1), -1)
	}

	// FIXME: octree visualisation, make only edges draw, not faces
	//color := server.Component_RenderComponent_Color{
	//	Red:   1,
	//	Green: 1,
	//	Blue:  1,
	//	Alpha: 0.7,
	//}3
	//
	//nodes := c.objects.GetNodes()
	//for _, n := range nodes {
	//	r := n.GetRegion()
	//	center := r.GetCenter()
	//	size := r.GetSize()
	//	m.AddDebug(&center, *protometry.NewMeshSquareCuboid(size.X, true), color) // It's a cube anyway
	//}

	// Main loop
	lastUpdateTime := utils.GetProtoTime()
	for {
		dt := float64(utils.GetProtoTime()-lastUpdateTime) / math.Pow(10, 9)
		if dt > cfg.Get().UpdatesRate { // 50fps
			select {
			case msg := <-m.ClientsIn:
				m.OnClientPacket(msg.ClientToken, msg.Packet)
				//log.Zap.Info("received message %T", msg.Type)
			default:
				//log.Zap.Info("no message received")
			}
			//log.Zap.Info("blok")
			// This will usually be called within the erutan-loop, in order to update all Systems on every frame.
			m.World.Update(dt * cfg.Get().Logic.TimeScale)
			lastUpdateTime = utils.GetProtoTime()
		}
	}
}

// OnClientPacket handle packets coming from clients
func (m *Manager) OnClientPacket(tkn string, p Packet) {
	//log.Zap.Info("OnClientPacket %v", p.Type)
	switch t := p.Type.(type) {
	case *Packet_UpdateParameters:
		for _, element := range t.UpdateParameters.Parameters {
			switch param := element.Type.(type) {
			case *Packet_UpdateParametersPacket_Parameter_CullingArea,
				*Packet_UpdateParametersPacket_Parameter_Debug:
				ManagerInstance.NotifyAll(obs.Event{Value: ClientSettingsUpdate{ClientToken: tkn, Settings: *t}})
			case *Packet_UpdateParametersPacket_Parameter_TimeScale:
				cfg.Get().Logic.TimeScale = param.TimeScale
			}
		}
	case *Packet_UpdateObject:
		// This case is actually misleading, UpdateObject client -> server is used to create new objects, not update
		var sc Component_SpaceComponent
		var behaviour Component_BehaviourTypeComponent
		for _, ct := range p.GetUpdateObject().Components {
			switch c := ct.Type.(type) {
			case *Component_Space:
				sc = *c.Space
			case *Component_BehaviourType:
				behaviour = *c.BehaviourType
			}
		}
		if sc.Position == nil {
			log.Zap.Info("Client requested object update without required args")
			return
		}
		// Create object
		switch behaviour.Tag {
		case Component_BehaviourTypeComponent_ANY:
		case Component_BehaviourTypeComponent_ANIMAL:
			m.AddHerbivorous(sc.Position, protometry.NewVector3(1, 1, 1), -1)
		case Component_BehaviourTypeComponent_VEGETATION: // TODO
		case Component_BehaviourTypeComponent_PLAYER: // Are clients allowed to spawn players ?
		}

	case *Packet_UpdateSpaceRequest:
		// Update requested object after applying physics
		if t.UpdateSpaceRequest.ActualSpace.Position == nil ||
			t.UpdateSpaceRequest.NewSpace.Position == nil {
			log.Zap.Info("Client requested object space update with incorrect args")
			return
		}

		// Let's find this object in the state
		b := protometry.NewBoxOfSize(t.UpdateSpaceRequest.ActualSpace.Position.X,
			t.UpdateSpaceRequest.ActualSpace.Position.Y,
			t.UpdateSpaceRequest.ActualSpace.Position.Z,
			cfg.Get().Logic.OctreeSize)
		o := m.networkSystem.objects.Get(t.UpdateSpaceRequest.ObjectId, *b)

		if o != nil {
			// Ignores physics
			ManagerInstance.NotifyAll(obs.Event{Value: PhysicsUpdateResponse{
				Objects: []struct {
					octree.Object
					protometry.Vector3
				}{{Object: *o.Clone(), Vector3: *t.UpdateSpaceRequest.NewSpace.Position}}}})
		} else { // Update object
			log.Zap.Info("Client tried to update an in-existent object", zap.String("client", tkn),
				zap.Uint64("ID", t.UpdateSpaceRequest.ObjectId))
		}
		//log.Zap.Info("Object created at %v", sc.Position)
	case *Packet_Armageddon:
		log.Zap.Info("Start armageddon")
		// My theory is that this function is ran in another goroutine so it's better to get all objects at this given time
		// One shot and then act instead of iterating
		objs := m.networkSystem.objects.GetAllObjects() // TODO: disable suicide ?
		for _, obj := range objs {
			m.World.RemoveObject(obj)
		}
	case *Packet_DestroyObject:
		destroy := p.GetDestroyObject()
		o := m.networkSystem.objects.Get(destroy.ObjectId, *destroy.Region)
		if o != nil {
			log.Zap.Info("Client destroy", zap.String("client", tkn), zap.Uint64("ID", destroy.ObjectId))
			m.World.RemoveObject(*o)
		} else {
			log.Zap.Info("Client tried to destroy in-existent object", zap.String("client", tkn), zap.Uint64("ID", destroy.ObjectId))
		}
	default:
		log.Zap.Info("Client sent unimplemented packet", zap.Any("packet", t))
	}
}

func (m *Manager) Send(clientToken string, packet Packet) {
	packet.Metadata = &Metadata{Timestamp: ptypes.TimestampNow()}
	// Sync map useful in case a client disconnect while we're trying to send him something
	if inter, ok := m.ClientsOut.Load(clientToken); ok {
		if channel, ok2 := inter.(chan Packet); ok2 {
			channel <- packet
		}
	}
}

func (m *Manager) Broadcast(packet Packet) {
	packet.Metadata = &Metadata{Timestamp: ptypes.TimestampNow()}
	m.BroadcastOut <- packet
}
