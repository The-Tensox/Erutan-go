package game

import (
	"math"
	"math/rand"
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

	h := &HerbivorousSystem{}

	e := &EatableSystem{}
	m.World.AddSystem(&CollisionSystem{})
	m.World.AddSystem(h)
	m.World.AddSystem(e)
	m.World.AddSystem(&NetworkSystem{lastUpdateTime: utils.GetProtoTime()})

	m.Watch.Add(h)
	m.Watch.Add(e)

	id := ecs.NewBasic()
	ground := AnyObject{BasicEntity: &id}
	ground.Component_SpaceComponent = &erutan.Component_SpaceComponent{
		Position: &erutan.NetVector3{X: 0, Y: 0, Z: 0},
		Rotation: &erutan.NetQuaternion{X: 0, Y: 0, Z: 0, W: 0},
		Scale:    &erutan.NetVector3{X: utils.Config.GroundSize, Y: 1, Z: utils.Config.GroundSize},
	}
	ground.Component_RenderComponent = &erutan.Component_RenderComponent{
		Red:   0,
		Green: 0,
		Blue:  1,
	}

	for i := 0; i < 20; i++ {
		id := ecs.NewBasic()
		herb := AnyObject{BasicEntity: &id}
		herb.Component_SpaceComponent = &erutan.Component_SpaceComponent{
			Position: utils.RandomPositionInsideCircle(utils.Config.GroundSize / 2),
			Rotation: &erutan.NetQuaternion{X: 0, Y: 0, Z: 0, W: 0},
			Scale:    &erutan.NetVector3{X: 1, Y: 1, Z: 1},
		}
		herb.Component_RenderComponent = &erutan.Component_RenderComponent{
			Red:   0,
			Green: 1,
			Blue:  0,
		}
		herb.Component_BehaviourTypeComponent = &erutan.Component_BehaviourTypeComponent{
			BehaviourType: erutan.Component_BehaviourTypeComponent_VEGETATION,
		}
		// Add our entity to the appropriate systems
		for _, system := range m.World.Systems() {
			switch sys := system.(type) {
			case *CollisionSystem:
				sys.Add(herb.BasicEntity, herb.Component_SpaceComponent, herb.Component_BehaviourTypeComponent)
			case *EatableSystem:
				sys.Add(herb.BasicEntity, herb.Component_SpaceComponent)
			case *NetworkSystem:
				sys.Add(herb.BasicEntity, []*erutan.Component{
					&erutan.Component{Type: &erutan.Component_Space{Space: herb.Component_SpaceComponent}},
					&erutan.Component{Type: &erutan.Component_Render{Render: herb.Component_RenderComponent}},
				})
			case *RenderSystem:
				sys.Add(herb.BasicEntity, herb.Component_RenderComponent)
			}
		}
	}

	for i := 0; i < 5; i++ {
		m.AddHerbivorous(utils.RandomPositionInsideCircle(utils.Config.GroundSize/2), &erutan.NetVector3{X: 1, Y: 1, Z: 1}, -1)
	}
	// Add our entity to the appropriate systems
	for _, system := range m.World.Systems() {
		switch sys := system.(type) {
		case *NetworkSystem:
			sys.Add(ground.BasicEntity, []*erutan.Component{
				&erutan.Component{Type: &erutan.Component_Space{Space: ground.Component_SpaceComponent}},
				&erutan.Component{Type: &erutan.Component_Render{Render: ground.Component_RenderComponent}},
			})
		case *RenderSystem:
			sys.Add(ground.BasicEntity, ground.Component_RenderComponent)
		}
	}
	lastUpdateTime := utils.GetProtoTime()
	for {
		dt := float64(utils.GetProtoTime()-lastUpdateTime) / math.Pow(10, 9)
		if dt > 0.0001 { // 50fps
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

func (m *Manager) AddHerbivorous(position *erutan.NetVector3, scale *erutan.NetVector3, speed float64) {
	id := ecs.NewBasic()
	herbivorous := Herbivorous{BasicEntity: &id}
	herbivorous.Component_HealthComponent = &erutan.Component_HealthComponent{Life: 40}
	herbivorous.Component_SpaceComponent = &erutan.Component_SpaceComponent{
		Position: position,
		Rotation: &erutan.NetQuaternion{X: 0, Y: 0, Z: 0, W: 0},
		Scale:    scale,
	}
	herbivorous.Target = nil // target
	herbivorous.Component_RenderComponent = &erutan.Component_RenderComponent{
		Red:   1,
		Green: 0,
		Blue:  0,
	}
	herbivorous.Component_BehaviourTypeComponent = &erutan.Component_BehaviourTypeComponent{
		BehaviourType: erutan.Component_BehaviourTypeComponent_ANIMAL,
	}
	// Default param
	if speed == -1 {
		speed = 10 + rand.Float64()*10
	}
	herbivorous.Component_SpeedComponent = &erutan.Component_SpeedComponent{
		MoveSpeed: speed,
	}
	// Add our herbivorous to the appropriate systems
	for _, system := range m.World.Systems() {
		switch sys := system.(type) {
		case *CollisionSystem:
			sys.Add(herbivorous.BasicEntity, herbivorous.Component_SpaceComponent, herbivorous.Component_BehaviourTypeComponent)
		case *HerbivorousSystem:
			sys.Add(herbivorous.BasicEntity,
				herbivorous.Component_SpaceComponent,
				herbivorous.Target,
				herbivorous.Component_HealthComponent,
				herbivorous.Component_SpeedComponent)
		case *NetworkSystem:
			sys.Add(herbivorous.BasicEntity, []*erutan.Component{
				&erutan.Component{Type: &erutan.Component_Space{Space: herbivorous.Component_SpaceComponent}},
				&erutan.Component{Type: &erutan.Component_Render{Render: herbivorous.Component_RenderComponent}},
				&erutan.Component{Type: &erutan.Component_Health{Health: herbivorous.Component_HealthComponent}},
				&erutan.Component{Type: &erutan.Component_Speed{Speed: herbivorous.Component_SpeedComponent}},
			})
		case *RenderSystem:
			sys.Add(herbivorous.BasicEntity, herbivorous.Component_RenderComponent)
		}
	}
}
