package game

import (
	"math"
	"sync"

	ecs "github.com/user/erutan/ecs"
	erutan "github.com/user/erutan/protos/realtime"
	utils "github.com/user/erutan/utils"
)

var (
	once sync.Once

	// ManagerInstance is a global singleton game manager
	ManagerInstance *Manager
)

// Manager ...
type Manager struct {
	World ecs.World

	ClientNames          map[string]string
	ClientStreams        map[string]chan erutan.Packet
	NamesMtx, StreamsMtx sync.RWMutex

	// Broadcast is a global channel to send packets to clients
	Broadcast chan erutan.Packet

	// Receive receive packets from clients
	Receive chan *erutan.Packet

	Watch utils.Watch
}

// Initialize returns a thread-safe singleton instance of the game manager
func Initialize() {
	once.Do(func() {
		ManagerInstance =
			&Manager{
				World:         ecs.World{},
				ClientNames:   make(map[string]string),
				ClientStreams: make(map[string]chan erutan.Packet),
				Broadcast:     make(chan erutan.Packet, 1000),
				Receive:       make(chan *erutan.Packet, 1000),
				Watch:         utils.Watch{},
			}
	})
}

// Run start handling gameplay
func (m *Manager) Run() {
	go m.Listen()

	r := &ReachTargetSystem{}
	e := &EatableSystem{}
	m.World.AddSystem(&CollisionSystem{})
	m.World.AddSystem(r)
	m.World.AddSystem(e)
	m.World.AddSystem(&NetworkSystem{})

	m.Watch.Add(r)
	m.Watch.Add(e)

	ground := AnyObject{BasicEntity: ecs.NewBasic()}
	ground.Component_SpaceComponent = erutan.Component_SpaceComponent{
		Position: &erutan.NetVector3{X: 0, Y: 0, Z: 0},
		Rotation: &erutan.NetQuaternion{X: 0, Y: 0, Z: 0, W: 0},
		Scale:    &erutan.NetVector3{X: 100, Y: 1, Z: 100},
	}
	ground.Component_RenderComponent = erutan.Component_RenderComponent{
		Red:   0,
		Green: 0,
		Blue:  1,
	}

	herb := AnyObject{BasicEntity: ecs.NewBasic()}
	herb.Component_SpaceComponent = erutan.Component_SpaceComponent{
		Position: utils.RandomPositionInsideCircle(50),
		Rotation: &erutan.NetQuaternion{X: 0, Y: 0, Z: 0, W: 0},
		Scale:    &erutan.NetVector3{X: 1, Y: 1, Z: 1},
	}
	herb.Component_RenderComponent = erutan.Component_RenderComponent{
		Red:   0,
		Green: 1,
		Blue:  0,
	}

	herbivorous := Herbivorous{BasicEntity: ecs.NewBasic()}
	herbivorous.Component_HealthComponent = erutan.Component_HealthComponent{Life: 40}
	herbivorous.Component_SpaceComponent = erutan.Component_SpaceComponent{
		Position: utils.RandomPositionInsideCircle(50),
		Rotation: &erutan.NetQuaternion{X: 0, Y: 0, Z: 0, W: 0},
		Scale:    &erutan.NetVector3{X: 1, Y: 1, Z: 1},
	}
	herbivorous.Component_TargetComponent = erutan.Component_TargetComponent{Target: herb.Position}
	herbivorous.Component_RenderComponent = erutan.Component_RenderComponent{
		Red:   1,
		Green: 0,
		Blue:  0,
	}
	// Add our entity to the appropriate systems
	for _, system := range m.World.Systems() {
		switch sys := system.(type) {
		case *CollisionSystem:
			sys.Add(&herbivorous.BasicEntity, &herbivorous.Component_SpaceComponent)
			sys.Add(&herb.BasicEntity, &herb.Component_SpaceComponent)
			sys.Add(&ground.BasicEntity, &ground.Component_SpaceComponent)
		case *ReachTargetSystem:
			sys.Add(&herbivorous.BasicEntity,
				&herbivorous.Component_SpaceComponent,
				&herbivorous.Component_TargetComponent,
				&herbivorous.Component_HealthComponent)
		case *EatableSystem:
			sys.Add(&herb.BasicEntity, &herb.Component_SpaceComponent)
		case *NetworkSystem:
			sys.Add(&herbivorous.BasicEntity, []*erutan.Component{
				&erutan.Component{Type: &erutan.Component_Space{Space: &herbivorous.Component_SpaceComponent}},
				&erutan.Component{Type: &erutan.Component_Render{Render: &herbivorous.Component_RenderComponent}},
				&erutan.Component{Type: &erutan.Component_Health{Health: &herbivorous.Component_HealthComponent}},
			})
			sys.Add(&herb.BasicEntity, []*erutan.Component{
				&erutan.Component{Type: &erutan.Component_Space{Space: &herb.Component_SpaceComponent}},
				&erutan.Component{Type: &erutan.Component_Render{Render: &herb.Component_RenderComponent}},
			})
			sys.Add(&ground.BasicEntity, []*erutan.Component{
				&erutan.Component{Type: &erutan.Component_Space{Space: &ground.Component_SpaceComponent}},
				&erutan.Component{Type: &erutan.Component_Render{Render: &ground.Component_RenderComponent}},
			})
		case *RenderSystem:
			sys.Add(&herbivorous.BasicEntity, &herbivorous.Component_RenderComponent)
			sys.Add(&herb.BasicEntity, &herb.Component_RenderComponent)
			sys.Add(&ground.BasicEntity, &ground.Component_RenderComponent)
		}
	}
	lastUpdateTime := utils.GetProtoTime()
	for {
		dt := float64(utils.GetProtoTime()-lastUpdateTime) / math.Pow(10, 9)
		if dt > 0.02 { // 50fps
			// This will usually be called within the game-loop, in order to update all Systems on every frame.
			m.World.Update(dt * utils.Config.TimeScale)
			lastUpdateTime = utils.GetProtoTime()
		}
	}
}

// Listen ...
func (m *Manager) Listen() {
	for {
		select {
		case p := <-m.Receive:
			switch t := p.Type.(type) {
			case *erutan.Packet_UpdateParameters:
				for _, element := range t.UpdateParameters.Parameters {
					switch param := element.Type.(type) {
					case *erutan.Packet_UpdateParametersPacket_Parameter_TimeScale:
						utils.Config.TimeScale = param.TimeScale
					}
				}
			case *erutan.Packet_CreateEntity:
			default:
				utils.DebugLogf("Client sent unimplemented packet: %v", t)
			}
		}
	}
}

func (m *Manager) SyncNewClient(tkn string) {
	for _, system := range m.World.Systems() {
		switch sys := system.(type) {
		case *NetworkSystem:
			sys.SyncNewClient(tkn)
		}
	}
}
