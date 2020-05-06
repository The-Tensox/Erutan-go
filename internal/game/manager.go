package game

import (
	"github.com/The-Tensox/Erutan-go/internal/cfg"
	"github.com/The-Tensox/Erutan-go/internal/ecs"
	"github.com/The-Tensox/Erutan-go/internal/obs"
	"github.com/The-Tensox/Erutan-go/internal/utils"
	erutan "github.com/The-Tensox/Erutan-go/protobuf"
	"github.com/The-Tensox/octree"
	"github.com/The-Tensox/protometry"
	"github.com/golang/protobuf/ptypes"
	"math"
	"sync"
)

var (
	once sync.Once

	// ManagerInstance is a global singleton game manager
	ManagerInstance *Manager
)

// Manager ...
type Manager struct {
	// World is the structure that handle all Systems in the Object Component System design
	World ecs.World

	// ClientsIn handle incoming messages, doesn't need to be a map since when we receive we know it's a connected client
	ClientsIn chan erutan.ClientPacket
	//ClientsMu sync.RWMutex
	// ClientsOut is a map that send packets to clients (map[string]chan Packet)
	ClientsOut sync.Map
	// Map a client to settings (map[string]ClientSettings)
	ClientsSettings sync.Map

	// BroadcastOut is a channel to send packets to every clients
	BroadcastOut chan erutan.Packet

	Watch obs.Watch

	// Reference to the network octree, needed to iterate networked objects
	networkSystem *NetworkSystem
}

// Initialize returns a thread-safe singleton instance of the game manager
func Initialize() {
	once.Do(func() {
		ManagerInstance =
			&Manager{
				World:      ecs.World{},
				ClientsIn:  make(chan erutan.ClientPacket, 100000),
				BroadcastOut:  make(chan erutan.Packet, 1000),
				Watch:      *obs.NewWatch(),
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
	m.World.AddSystem(h)
	m.World.AddSystem(e)
	m.World.AddSystem(c)
	m.World.AddSystem(n)

	m.Watch.Register(h)
	m.Watch.Register(e)
	m.Watch.Register(c)
	m.Watch.Register(n)

	//gs := cfg.Global.Logic.GroundSize
	//p := perlin.NewPerlin(1, 1, 1, 1337)
	//for x := 0.; x < gs; x++ {
	//	for y := 0.; y < gs; y++ {
	//		noise := p.Noise2D(x/10, y/10)
	//		//utils.DebugLogf("Noise at %.1f;%.1f: %v", x, y, noise)
	//		m.AddGround(protometry.NewVector3(x, float64(int(10*noise))-5, y), 1)
	//	}
	//}
	//ds := DiamondSquareAlgorithm(int(math.Pow(2, 8))+1, 40, 1)
	//sideLength := 40.
	//for x := 0.; int(x) < len(ds); x++ {
	//	for y := 0.; int(y) < len(ds[int(x)]); y++ {
	//		m.AddGround(protometry.NewVector3(x*sideLength,
	//			float64(int(sideLength*ds[int(x)][int(y)])),
	//			y*sideLength),
	//			sideLength)
	//	}
	//}

	for i := 0; i < cfg.Global.Logic.InitialHerbs; i++ {
		p := protometry.RandomCirclePoint(0, 0,0, 50)
		m.AddHerb(&p)
	}

	for i := 0; i < cfg.Global.Logic.Herbivorous.Quantity; i++ {
		p := protometry.RandomCirclePoint(0, 0,0, 50)
		m.AddHerbivorous(&p, protometry.NewVector3(1, 1, 1), -1)
	}

	// FIXME: octree visualisation, make only edges draw, not faces
	//color := erutan.Component_RenderComponent_Color{
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
		if dt > cfg.Global.UpdatesRate { // 50fps
			select {
			case msg := <-m.ClientsIn:
				m.OnClientPacket(msg.ClientToken, msg.Packet)
				//utils.DebugLogf("received message %T", msg.Type)
			default:
				//utils.DebugLogf("no message received")
			}
			//utils.DebugLogf("blok")
			// This will usually be called within the game-loop, in order to update all Systems on every frame.
			m.World.Update(dt * cfg.Global.Logic.TimeScale)
			lastUpdateTime = utils.GetProtoTime()
		}
	}
}

// OnClientPacket handle packets coming from clients
func (m *Manager) OnClientPacket(tkn string, p erutan.Packet) {
	//utils.DebugLogf("OnClientPacket %v", p.Type)
	switch t := p.Type.(type) {
	case *erutan.Packet_UpdateParameters:
		for _, element := range t.UpdateParameters.Parameters {
			switch param := element.Type.(type) {
			case *erutan.Packet_UpdateParametersPacket_Parameter_CullingArea,
				*erutan.Packet_UpdateParametersPacket_Parameter_Debug:
				ManagerInstance.Watch.NotifyAll(obs.Event{Value: obs.ClientSettingsUpdate{ClientToken: tkn, Settings: *t}})
			case *erutan.Packet_UpdateParametersPacket_Parameter_TimeScale:
				utils.DebugLogf("[%s] changed global timescale from %v to %v",
					tkn,
					cfg.Global.Logic.TimeScale,
					param.TimeScale)
				cfg.Global.Logic.TimeScale = param.TimeScale
			}
		}
	case *erutan.Packet_UpdateObject:
		// This case is actually misleading, UpdateObject client -> server is used to create new objects, not update
		var sc erutan.Component_SpaceComponent
		var behaviour erutan.Component_BehaviourTypeComponent
		for _, ct := range p.GetUpdateObject().Components {
			switch c := ct.Type.(type) {
			case *erutan.Component_Space:
				sc = *c.Space
			case *erutan.Component_BehaviourType:
				behaviour = *c.BehaviourType
			}
		}
		if sc.Position == nil {
			utils.DebugLogf("Client requested object update without required args")
			return
		}
		// Create object
		switch behaviour.Tag {
		case erutan.Component_BehaviourTypeComponent_ANY:
		case erutan.Component_BehaviourTypeComponent_ANIMAL:
			m.AddHerbivorous(sc.Position, protometry.NewVector3(1, 1, 1), -1)
		case erutan.Component_BehaviourTypeComponent_VEGETATION: // TODO
		case erutan.Component_BehaviourTypeComponent_PLAYER: // Are clients allowed to spawn players ?
		}

	case *erutan.Packet_UpdateSpaceRequest:
		// Update requested object after applying physics
		if t.UpdateSpaceRequest.ActualSpace.Position == nil ||
			t.UpdateSpaceRequest.NewSpace.Position == nil {
			utils.DebugLogf("Client requested object space update with incorrect args")
			return
		}

		// Let's find this object in the state
		b := protometry.NewBoxOfSize(t.UpdateSpaceRequest.ActualSpace.Position.X,
			t.UpdateSpaceRequest.ActualSpace.Position.Y,
			t.UpdateSpaceRequest.ActualSpace.Position.Z,
			cfg.Global.Logic.OctreeSize)
		//min := protometry.Min(*t.UpdateSpaceRequest.ActualSpace.Position, *t.UpdateSpaceRequest.NewSpace.Position)
		//max := protometry.Min(*t.UpdateSpaceRequest.ActualSpace.Position, *t.UpdateSpaceRequest.NewSpace.Position)
		//b := protometry.NewBoxMinMax(min.X, min.Y, min.Z, max.X, max.Y, max.Z)
		o := m.networkSystem.objects.Get(t.UpdateSpaceRequest.ObjectId, *b)

		if o != nil {
			//utils.DebugLogf("Client [%s] request update to %v, actual %v", tkn,
			//	t.UpdateSpaceRequest.NewSpace.Position,
			//	t.UpdateSpaceRequest.ActualSpace.Position)

			// Ignores physics
			ManagerInstance.Watch.NotifyAll(obs.Event{Value: obs.PhysicsUpdateResponse{
				Objects: []struct{octree.Object;protometry.Vector3}{{Object: *o.Clone(), Vector3: *t.UpdateSpaceRequest.NewSpace.Position}}}})
		} else { // Update object
			utils.DebugLogf("Client [%s] tried to update an in-existent object %d", tkn, t.UpdateSpaceRequest.ObjectId)
		}
		//utils.DebugLogf("Object created at %v", sc.Position)
	case *erutan.Packet_Armageddon:
		utils.DebugLogf("Start armageddon")
		// My theory is that this function is ran in another goroutine so it's better to get all objects at this given time
		// One shot and then act instead of iterating
		objs := m.networkSystem.objects.GetAllObjects() // TODO: disable suicide ?
		for _, obj := range objs {
			m.World.RemoveObject(obj)
		}
		//m.networkOctree.Range(func(object *octree.Object) bool {
		//	m.World.RemoveObject(*object)
		//	return true
		//})
	case *erutan.Packet_DestroyObject:
		destroy := p.GetDestroyObject()
		//utils.DebugLogf("lol %v", m.networkOctree.GetColliding(*destroy.Region))
		o := m.networkSystem.objects.Get(destroy.ObjectId, *destroy.Region)
		if o != nil {
			utils.DebugLogf("Client %s destroy %d", tkn, destroy.ObjectId)
			m.World.RemoveObject(*o)
		} else {
			utils.DebugLogf("Client %s tried to destroy in-existent object %d", tkn, destroy.ObjectId)
			//utils.DebugLogf("network objets : %v", m.networkOctree.GetAllObjects())
		}
	default:
		utils.DebugLogf("Client sent unimplemented packet: %v", t)
	}
}


func (m *Manager) Send(clientToken string, packet erutan.Packet) {
	packet.Metadata = &erutan.Metadata{Timestamp: ptypes.TimestampNow()}
	// Sync map useful in case a client disconnect while we're trying to send him something
	if inter, ok := m.ClientsOut.Load(clientToken); ok {
		if channel, ok2 := inter.(chan erutan.Packet); ok2 {
			channel <- packet
		}
	}
}

func (m *Manager) Broadcast(packet erutan.Packet) {
	packet.Metadata = &erutan.Metadata{Timestamp: ptypes.TimestampNow()}
	m.BroadcastOut<-packet
}