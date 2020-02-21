package game

import (
	"math"
	"sync"

	"github.com/aquilax/go-perlin"
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

	ClientStreams sync.Map

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
				World:     ecs.World{},
				Broadcast: make(chan erutan.Packet, 1000),
				Receive:   make(chan *erutan.Packet, 1000),
				Watch:     utils.Watch{},
			}
	})
}

// Run start handling gameplay
func (m *Manager) Run() {
	go m.Listen()

	h := &HerbivorousSystem{}
	e := &EatableSystem{}
	c := &CollisionSystem{}
	//c.New(&m.World)
	m.World.AddSystem(c)
	m.World.AddSystem(h)
	m.World.AddSystem(e)
	m.World.AddSystem(&NetworkSystem{lastUpdateTime: utils.GetProtoTime()})

	m.Watch.Add(h)
	m.Watch.Add(e)

	p := perlin.NewPerlin(1, 1, 5, 100)
	for x := 0.; x < utils.Config.GroundSize; x++ {
		for y := 0.; y < utils.Config.GroundSize; y++ {
			noise := p.Noise2D(x/10, y/10)
			//fmt.Printf("%0.0f\t%0.0f\t%0.4f\n", x, y, noise)
			m.AddGround(&erutan.NetVector3{X: x, Y: noise, Z: y}, 1)
			//m.AddHerb(&erutan.NetVector3{X: x, Y: 10, Z: y})
		}
	}

	//m.AddGround(&erutan.NetVector3{X: 0, Y: -utils.Config.GroundSize, Z: 0}, utils.Config.GroundSize)
	/*
		for i := 0; i < 10; i++ {
			m.AddGround(utils.RandomPositionInsideSphere(&erutan.NetVector3{X: 0, Y: 0, Z: 0}, 10))
		}
	*/
	// Debug thing, wait client
	/*
		nbClients := 0
		for nbClients == 0 {
			m.ClientStreams.Range(func(key interface{}, value interface{}) bool {
				nbClients++
				return true
			})
		}
	*/
	for i := 0; i < 1; i++ {
		m.AddHerb(utils.RandomPositionInsideCircle(&erutan.NetVector2{X: utils.Config.GroundSize / 2, Y: utils.Config.GroundSize / 2},
			utils.Config.GroundSize/2))
	}

	for i := 0; i < 0; i++ {
		m.AddHerbivorous(utils.RandomPositionInsideCircle(&erutan.NetVector2{X: utils.Config.GroundSize / 2, Y: utils.Config.GroundSize / 2},
			utils.Config.GroundSize/2), &erutan.NetVector3{X: 1, Y: 1, Z: 1}, -1)
	}

	// Main loop
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
				// Only handle herbivorous and client only has access to position
				var sc erutan.Component_SpaceComponent
				for _, c := range p.GetCreateEntity().Components {
					if tmp := c.GetSpace(); tmp != nil {
						sc = *tmp
					}
				}
				m.AddHerbivorous(sc.Position, &erutan.NetVector3{X: 1, Y: 1, Z: 1}, -1)
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

func (m *Manager) AddGround(position *erutan.NetVector3, sideLength float64) {
	id := ecs.NewBasic()
	ground := AnyObject{BasicEntity: &id}
	ground.Component_SpaceComponent = &erutan.Component_SpaceComponent{
		Position: position,
		Rotation: &erutan.NetQuaternion{X: 0, Y: 0, Z: 0, W: 0},
		Scale:    &erutan.NetVector3{X: 1, Y: 1, Z: 1},
		Shape:    utils.CreateCube(sideLength),
	}
	ground.Component_RenderComponent = &erutan.Component_RenderComponent{
		Red:   0,
		Green: -float32(position.Y),
		Blue:  0,
	}
	ground.Component_BehaviourTypeComponent = &erutan.Component_BehaviourTypeComponent{
		BehaviourType: erutan.Component_BehaviourTypeComponent_ANY,
	}
	ground.Component_PhysicsComponent = &erutan.Component_PhysicsComponent{
		UseGravity: false,
	}
	// Add our entity to the appropriate systems
	for _, system := range m.World.Systems() {
		switch sys := system.(type) {
		case *CollisionSystem:
			sys.Add(ground.BasicEntity,
				ground.Component_SpaceComponent,
				ground.Component_BehaviourTypeComponent,
				ground.Component_PhysicsComponent)
		case *NetworkSystem:
			sys.Add(ground.BasicEntity, []*erutan.Component{
				&erutan.Component{Type: &erutan.Component_Space{Space: ground.Component_SpaceComponent}},
				&erutan.Component{Type: &erutan.Component_Render{Render: ground.Component_RenderComponent}},
			})
		case *RenderSystem:
			sys.Add(ground.BasicEntity, ground.Component_RenderComponent)
		}
	}
}

func (m *Manager) AddHerb(position *erutan.NetVector3) {
	id := ecs.NewBasic()
	herb := AnyObject{BasicEntity: &id}
	herb.Component_SpaceComponent = &erutan.Component_SpaceComponent{
		Position: position,
		Rotation: &erutan.NetQuaternion{X: 0, Y: 0, Z: 0, W: 0},
		Scale:    &erutan.NetVector3{X: 1, Y: 1, Z: 1},
		Shape:    utils.CreateCube(1),
	}
	herb.Component_RenderComponent = &erutan.Component_RenderComponent{
		Red:   0,
		Green: 0,
		Blue:  1,
	}
	herb.Component_BehaviourTypeComponent = &erutan.Component_BehaviourTypeComponent{
		BehaviourType: erutan.Component_BehaviourTypeComponent_VEGETATION,
	}
	herb.Component_PhysicsComponent = &erutan.Component_PhysicsComponent{
		UseGravity: true,
	}
	// Add our entity to the appropriate systems
	for _, system := range m.World.Systems() {
		switch sys := system.(type) {
		case *CollisionSystem:
			sys.Add(herb.BasicEntity,
				herb.Component_SpaceComponent,
				herb.Component_BehaviourTypeComponent,
				herb.Component_PhysicsComponent)
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

func (m *Manager) AddHerbivorous(position *erutan.NetVector3, scale *erutan.NetVector3, speed float64) {
	id := ecs.NewBasic()
	herbivorous := Herbivorous{BasicEntity: &id}
	herbivorous.Component_HealthComponent = &erutan.Component_HealthComponent{Life: 40}
	herbivorous.Component_SpaceComponent = &erutan.Component_SpaceComponent{
		Position: position,
		Rotation: &erutan.NetQuaternion{X: 0, Y: 0, Z: 0, W: 0},
		Scale:    scale,
		Shape:    utils.CreateCube(1),
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
		speed = utils.RandFloats(10, 20)
	}
	herbivorous.Component_SpeedComponent = &erutan.Component_SpeedComponent{
		MoveSpeed: speed,
	}
	herbivorous.Component_PhysicsComponent = &erutan.Component_PhysicsComponent{
		UseGravity: true,
	}
	// Add our herbivorous to the appropriate systems
	for _, system := range m.World.Systems() {
		switch sys := system.(type) {
		case *CollisionSystem:
			sys.Add(herbivorous.BasicEntity,
				herbivorous.Component_SpaceComponent,
				herbivorous.Component_BehaviourTypeComponent,
				herbivorous.Component_PhysicsComponent)
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
