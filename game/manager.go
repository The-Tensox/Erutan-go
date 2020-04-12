package game

import (
	"math"
	"sync"

	ecs "github.com/The-Tensox/erutan/ecs"
	erutan "github.com/The-Tensox/erutan/protobuf"
	utils "github.com/The-Tensox/erutan/utils"
	"github.com/The-Tensox/protometry"
)

var (
	once sync.Once

	// ManagerInstance is a global singleton game manager
	ManagerInstance *Manager
)

// Manager ...
// TODO: later create settings / config for manager ?
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

	h := NewHerbivorousSystem()
	e := NewEatableSystem()
	c := NewCollisionSystem()
	m.World.AddSystem(c)
	m.World.AddSystem(h)
	m.World.AddSystem(e)
	m.World.AddSystem(NewRenderSystem())
	m.World.AddSystem(NewNetworkSystem(utils.GetProtoTime()))

	m.Watch.Add(h)
	m.Watch.Add(e)
	m.Watch.Add(c)

	gs := utils.Config.GroundSize - 1
	//p := perlin.NewPerlin(1, 10, 10, 100)
	//for x := 0.; x < gs; x++ {
	//	for y := 0.; y < gs; y++ {
	//		noise := p.Noise2D(x/10, y/10)
	//		//fmt.Printf("%0.0f\t%0.0f\t%0.4f\n", x, y, noise)
	//		m.AddGround(protometry.NewVectorN(x, noise, y), 1)
	//		m.AddHerb(protometry.NewVectorN(x, 5, y))
	//	}
	//}

	//m.AddGround(protometry.NewVectorN(0, -gs/2, 0), gs / 4)

	for x := 0.; x < gs; x++ {
		for z := 0.; z < gs; z++ {
			m.AddGround(protometry.NewVectorN(x, -10, z), 1)
		}
	}

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
	for i := 0; i < 5; i++ {
		p := protometry.RandomCirclePoint(*protometry.NewVectorN(gs/4, gs/4), gs/8)
		m.AddHerb(&p)
	}

	for i := 0; i < 5; i++ {
		p := protometry.RandomCirclePoint(*protometry.NewVectorN(gs/4, gs/4), gs/8)
		m.AddHerbivorous(&p, protometry.NewVectorN(1, 1, 1), -1)
	}

	// Main loop
	lastUpdateTime := utils.GetProtoTime()
	for {
		dt := float64(utils.GetProtoTime() - lastUpdateTime)/math.Pow(10, 9)
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
				m.AddHerbivorous(sc.Position, protometry.NewVectorN(1, 1, 1), -1)
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

func (m *Manager) AddGround(position *protometry.VectorN, sideLength float64) {
	id := ecs.NewId()
	ground := AnyObject{Id: id}
	ground.Component_SpaceComponent = &erutan.Component_SpaceComponent{
		Position: position,
		Rotation: protometry.NewQuaternion(0, 0, 0, 0),
		Scale:    protometry.NewVectorN(1, 1, 1),
		Shape:    utils.CreateCube(sideLength),
	}
	ground.Component_RenderComponent = &erutan.Component_RenderComponent{
		Red:   0,
		Green: -float32(position.Get(1)),
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
			sys.Add(id,
				sideLength / 2,
				ground.Component_SpaceComponent,
				ground.Component_BehaviourTypeComponent,
				ground.Component_PhysicsComponent)
		case *NetworkSystem:
			sys.Add(id, []*erutan.Component{
				{Type: &erutan.Component_Space{Space: ground.Component_SpaceComponent}},
				{Type: &erutan.Component_Render{Render: ground.Component_RenderComponent}},
			})
		case *RenderSystem:
			sys.Add(id, ground.Component_RenderComponent)
		}
	}
}

func (m *Manager) AddHerb(position *protometry.VectorN) {
	id := ecs.NewId()
	herb := AnyObject{Id: id}
	herb.Component_SpaceComponent = &erutan.Component_SpaceComponent{
		Position: position,
		Rotation: protometry.NewQuaternion(0, 0, 0, 0),
		Scale:    protometry.NewVectorN(1, 1, 1),
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
			sys.Add(id,
				0.5,
				herb.Component_SpaceComponent,
				herb.Component_BehaviourTypeComponent,
				herb.Component_PhysicsComponent)
		case *EatableSystem:
			sys.Add(id, herb.Component_SpaceComponent)
		case *NetworkSystem:
			sys.Add(id, []*erutan.Component{
				{Type: &erutan.Component_Space{Space: herb.Component_SpaceComponent}},
				{Type: &erutan.Component_Render{Render: herb.Component_RenderComponent}},
			})
		case *RenderSystem:
			sys.Add(id, herb.Component_RenderComponent)
		}
	}
}

func (m *Manager) AddHerbivorous(position *protometry.VectorN, scale *protometry.VectorN, speed float64) {
	id := ecs.NewId()
	herbivorous := Herbivorous{Id: id}
	herbivorous.Component_HealthComponent = &erutan.Component_HealthComponent{Life: 40}
	herbivorous.Component_SpaceComponent = &erutan.Component_SpaceComponent{
		Position: position,
		Rotation: protometry.NewQuaternion(0, 0, 0, 0),
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
			sys.Add(id,
				0.5,
				herbivorous.Component_SpaceComponent,
				herbivorous.Component_BehaviourTypeComponent,
				herbivorous.Component_PhysicsComponent)
		case *HerbivorousSystem:
			sys.Add(id,
				herbivorous.Component_SpaceComponent,
				herbivorous.Target,
				herbivorous.Component_HealthComponent,
				herbivorous.Component_SpeedComponent)
		case *NetworkSystem:
			sys.Add(id, []*erutan.Component{
				{Type: &erutan.Component_Space{Space: herbivorous.Component_SpaceComponent}},
				{Type: &erutan.Component_Render{Render: herbivorous.Component_RenderComponent}},
				{Type: &erutan.Component_Health{Health: herbivorous.Component_HealthComponent}},
				{Type: &erutan.Component_Speed{Speed: herbivorous.Component_SpeedComponent}},
			})
		case *RenderSystem:
			sys.Add(id, herbivorous.Component_RenderComponent)
		}
	}
}
