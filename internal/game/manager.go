package game

import (
	"github.com/The-Tensox/erutan/internal/cfg"
	"github.com/The-Tensox/erutan/internal/obs"
	"github.com/The-Tensox/octree"
	"math"
	"sync"

	ecs "github.com/The-Tensox/erutan/internal/ecs"
	utils "github.com/The-Tensox/erutan/internal/utils"
	erutan "github.com/The-Tensox/erutan/protobuf"
	"github.com/The-Tensox/protometry"
)

var (
	once sync.Once

	// ManagerInstance is a global singleton game manager
	ManagerInstance *Manager
)

// Manager ...
type Manager struct {
	// World is the structure that handle all Systems in the Entity Component System design
	World ecs.World

	// ClientsOut is a map that send packets to clients (map[string]chan Packet)
	ClientsOut sync.Map
	// Map a client to settings (map[string]ClientSettings)
	ClientsSettings sync.Map

	// Broadcast is a channel to send packets to every clients
	Broadcast chan erutan.Packet

	Watch obs.Watch
}

// Initialize returns a thread-safe singleton instance of the game manager
func Initialize() {
	once.Do(func() {
		ManagerInstance =
			&Manager{
				World:     ecs.World{},
				Broadcast: make(chan erutan.Packet, 1000),
				Watch:     *obs.NewWatch(),
			}
	})
}

// Run start handling gameplay
func (m *Manager) Run() {
	h := NewHerbivorousSystem()
	e := NewEatableSystem()
	c := NewCollisionSystem()
	n := NewNetworkSystem(utils.GetProtoTime())
	m.World.AddSystem(h)
	m.World.AddSystem(e)
	m.World.AddSystem(c)
	m.World.AddSystem(n)

	m.Watch.Register(h)
	m.Watch.Register(e)
	m.Watch.Register(c)
	m.Watch.Register(n)

	//gs := cfg.Global.Logic.GroundSize - 1
	//p := perlin.NewPerlin(1, 10, 10, 100)
	//for x := 0.; x < gs; x++ {
	//	for y := 0.; y < gs; y++ {
	//		noise := p.Noise2D(x/10, y/10)
	//		//fmt.Printf("%0.0f\t%0.0f\t%0.4f\n", x, y, noise)
	//		m.AddGround(protometry.NewVector3(x, noise-5, y), 1)
	//		//m.AddHerb(protometry.NewVector3(x, 5, y))
	//	}
	//}

	//m.AddGround(protometry.NewVector3(0, -gs/2, 0), gs / 4)

	//for x := 0.; x < gs; x++ {
	//	for z := 0.; z < gs; z++ {
	//		m.AddGround(protometry.NewVector3(x, -10, z), 1)
	//	}
	//}

	//for i := 0; i < cfg.Global.Logic.InitialHerbs; i++ {
	//	p := protometry.RandomCirclePoint(gs/4, gs/4, gs/8)
	//	m.AddHerb(&p)
	//}
	//
	//for i := 0; i < cfg.Global.Logic.InitialHerbivorous; i++ {
	//	p := protometry.RandomCirclePoint(gs/4, gs/4, gs/8)
	//	m.AddHerbivorous(&p, protometry.NewVector3(1, 1, 1), -1)
	//}

	for i := 0.; i < float64(cfg.Global.Logic.InitialHerbs); i++ {
		m.AddHerb(protometry.NewVector3(i*2, 0, i*2))
	}

	for i := 0.; i < float64(cfg.Global.Logic.InitialHerbivorous); i++ {
		m.AddHerbivorous(protometry.NewVector3(i*2+100, 0, i*2+100), protometry.NewVector3(1, 1, 1), -1)
	}

	color := erutan.Component_RenderComponent_Color{
		Red:   1,
		Green: 1,
		Blue:  1,
		Alpha: 0.7,
	}
	nodes := c.objects.GetNodes()
	for _, n := range nodes {
		r := n.GetRegion()
		center := r.GetCenter()
		size := r.GetSize()
		m.AddDebug(&center, *protometry.NewMeshSquareCuboid(size.X, true), color) // It's a cube anyway
	}

	//center := protometry.NewVector3(0, -10, 0)
	//mesh := protometry.NewMeshRectangularCuboid(*center, *protometry.NewVector3(10, 1, 1))
	//
	//m.AddDebug(center, *mesh, color)

	// Main loop
	lastUpdateTime := utils.GetProtoTime()
	for {
		dt := float64(utils.GetProtoTime()-lastUpdateTime) / math.Pow(10, 9)
		if dt > cfg.Global.FramesPerSecond { // 50fps
			// This will usually be called within the game-loop, in order to update all Systems on every frame.
			m.World.Update(dt * cfg.Global.Logic.TimeScale)
			lastUpdateTime = utils.GetProtoTime()
		}
	}
}

// Handle ...
func (m *Manager) Handle(tkn string, p erutan.Packet) {
	switch t := p.Type.(type) {
	case *erutan.Packet_UpdateParameters:
		// Update the client settings
		for _, p := range t.UpdateParameters.Parameters {
			switch p.Type.(type) {
			// No need to notify a change in timescale
			case *erutan.Packet_UpdateParametersPacket_Parameter_Debug:
				ManagerInstance.Watch.NotifyAll(obs.Event{Value: obs.OnClientSettingsUpdate{ClientToken: tkn, Settings: *t}})
			}
		}

		// Then do some global (dangerous :p)) logic
		for _, element := range t.UpdateParameters.Parameters {
			switch param := element.Type.(type) {
			case *erutan.Packet_UpdateParametersPacket_Parameter_TimeScale:
				utils.DebugLogf("[%s] changed global timescale from %v to %v",
					tkn,
					cfg.Global.Logic.TimeScale,
					param.TimeScale)
				cfg.Global.Logic.TimeScale = param.TimeScale
			}
		}
	case *erutan.Packet_UpdateEntity: // FIXME
		// Only handle herbivorous and client only has access to position
		var sc erutan.Component_SpaceComponent
		for _, c := range p.GetUpdateEntity().Components {
			if tmp := c.GetSpace(); tmp != nil {
				sc = *tmp
			}
		}
		m.AddHerbivorous(sc.Position, protometry.NewVector3(1, 1, 1), -1)
		utils.DebugLogf("Entity created at %v", sc.Position)
	default:
		utils.DebugLogf("Client sent unimplemented packet: %v", t)
	}
}

// TODO: bring together common building-blocks of these stuffs, it's too redundant

// AddDebug create a debug object that will only be seen by clients with debug settings
func (m *Manager) AddDebug(position *protometry.Vector3, mesh protometry.Mesh, color erutan.Component_RenderComponent_Color) {
	obj := AnyObject{}
	ocObj := octree.NewObjectCube(nil, position.X, position.Y, position.Z, 1)
	obj.Component_SpaceComponent = &erutan.Component_SpaceComponent{
		Position: position,
		Rotation: protometry.NewQuaternion(0, 0, 0, 0),
		Scale:    protometry.NewVector3(1, 1, 1),
	}
	var c []*erutan.Component_RenderComponent_Color
	for range mesh.Vertices {
		c = append(c, &color)
	}
	obj.Component_RenderComponent = &erutan.Component_RenderComponent{
		Mesh:   &mesh,
		Colors: c,
	}
	obj.Component_NetworkBehaviourComponent = &erutan.Component_NetworkBehaviourComponent{
		Tag: erutan.Component_NetworkBehaviourComponent_DEBUG,
	}
	// Add our entity to the appropriate systems
	for _, system := range m.World.Systems() {
		switch sys := system.(type) {
		case *NetworkSystem:
			sys.Add(*ocObj,
				[]*erutan.Component{
					{Type: &erutan.Component_Space{Space: obj.Component_SpaceComponent}},
					{Type: &erutan.Component_Render{Render: obj.Component_RenderComponent}},
					{Type: &erutan.Component_NetworkBehaviour{NetworkBehaviour: obj.Component_NetworkBehaviourComponent}},
				})
		}
	}
}

// AddGround create a ground object
func (m *Manager) AddGround(position *protometry.Vector3, sideLength float64) {
	obj := AnyObject{}
	ocObj := octree.NewObjectCube(nil, position.X, position.Y, position.Z, sideLength)
	obj.Component_SpaceComponent = &erutan.Component_SpaceComponent{
		Position: position,
		Rotation: protometry.NewQuaternion(0, 0, 0, 0),
		Scale:    protometry.NewVector3(1, 1, 1),
	}
	var c []*erutan.Component_RenderComponent_Color
	mesh := protometry.NewMeshSquareCuboid(1, true)
	for range mesh.Vertices {
		c = append(c, &erutan.Component_RenderComponent_Color{
			Red:   0,
			Green: -float32(position.Y),
			Blue:  0,
			Alpha: 1,
		})
	}
	obj.Component_RenderComponent = &erutan.Component_RenderComponent{
		Mesh:   mesh,
		Colors: c,
	}
	obj.Component_BehaviourTypeComponent = &erutan.Component_BehaviourTypeComponent{
		Tag: erutan.Component_BehaviourTypeComponent_ANY,
	}
	obj.Component_NetworkBehaviourComponent = &erutan.Component_NetworkBehaviourComponent{
		Tag: erutan.Component_NetworkBehaviourComponent_ALL,
	}
	obj.Component_PhysicsComponent = &erutan.Component_PhysicsComponent{
		UseGravity: false,
	}
	// Add our entity to the appropriate systems
	for _, system := range m.World.Systems() {
		switch sys := system.(type) {
		case *CollisionSystem:
			sys.Add(*ocObj,
				obj.Component_SpaceComponent,
				obj.Component_BehaviourTypeComponent,
				obj.Component_PhysicsComponent)
		case *NetworkSystem:
			sys.Add(*ocObj,
				[]*erutan.Component{
					{Type: &erutan.Component_Space{Space: obj.Component_SpaceComponent}},
					{Type: &erutan.Component_Render{Render: obj.Component_RenderComponent}},
					{Type: &erutan.Component_NetworkBehaviour{NetworkBehaviour: obj.Component_NetworkBehaviourComponent}},
				})
		}
	}
}

// AddHerb create an herb object
func (m *Manager) AddHerb(position *protometry.Vector3) {
	obj := AnyObject{}
	ocObj := octree.NewObjectCube(nil, position.X, position.Y, position.Z, 1)
	obj.Component_SpaceComponent = &erutan.Component_SpaceComponent{
		Position: position,
		Rotation: protometry.NewQuaternion(0, 0, 0, 0),
		Scale:    protometry.NewVector3(1, 1, 1),
	}
	var c []*erutan.Component_RenderComponent_Color
	mesh := protometry.NewMeshSquareCuboid(1, true)
	for range mesh.Vertices {
		c = append(c, &erutan.Component_RenderComponent_Color{
			Red:   0,
			Green: 0,
			Blue:  1,
			Alpha: 1,
		})
	}
	obj.Component_RenderComponent = &erutan.Component_RenderComponent{
		Mesh:   mesh,
		Colors: c,
	}
	obj.Component_BehaviourTypeComponent = &erutan.Component_BehaviourTypeComponent{
		Tag: erutan.Component_BehaviourTypeComponent_VEGETATION,
	}
	obj.Component_NetworkBehaviourComponent = &erutan.Component_NetworkBehaviourComponent{
		Tag: erutan.Component_NetworkBehaviourComponent_ALL,
	}
	obj.Component_PhysicsComponent = &erutan.Component_PhysicsComponent{
		UseGravity: true,
	}
	// Add our entity to the appropriate systems
	for _, system := range m.World.Systems() {
		switch sys := system.(type) {
		case *CollisionSystem:
			sys.Add(*ocObj,
				obj.Component_SpaceComponent,
				obj.Component_BehaviourTypeComponent,
				obj.Component_PhysicsComponent)
		case *EatableSystem:
			sys.Add(*ocObj,
				obj.Component_SpaceComponent)
		case *NetworkSystem:
			sys.Add(*ocObj,
				[]*erutan.Component{
					{Type: &erutan.Component_Space{Space: obj.Component_SpaceComponent}},
					{Type: &erutan.Component_Render{Render: obj.Component_RenderComponent}},
					{Type: &erutan.Component_NetworkBehaviour{NetworkBehaviour: obj.Component_NetworkBehaviourComponent}},
				})
		}
	}
}

// AddHerbivorous create an herbivorous object
func (m *Manager) AddHerbivorous(position *protometry.Vector3, scale *protometry.Vector3, speed float64) {
	obj := Herbivorous{}
	ocObj := octree.NewObjectCube(nil, position.X, position.Y, position.Z, 1)
	obj.Component_HealthComponent = &erutan.Component_HealthComponent{Life: cfg.Global.Logic.InitialHerbivorousLife}
	obj.Component_SpaceComponent = &erutan.Component_SpaceComponent{
		Position: position,
		Rotation: protometry.NewQuaternion(0, 0, 0, 0),
		Scale:    scale,
	}

	obj.Target = nil // target
	var c []*erutan.Component_RenderComponent_Color
	mesh := protometry.NewMeshSquareCuboid(1, true)
	for range mesh.Vertices {
		c = append(c, &erutan.Component_RenderComponent_Color{
			Red:   1,
			Green: 0,
			Blue:  0,
			Alpha: 1,
		})
	}
	obj.Component_RenderComponent = &erutan.Component_RenderComponent{
		Mesh:   mesh,
		Colors: c,
	}
	obj.Component_BehaviourTypeComponent = &erutan.Component_BehaviourTypeComponent{
		Tag: erutan.Component_BehaviourTypeComponent_ANIMAL,
	}
	obj.Component_NetworkBehaviourComponent = &erutan.Component_NetworkBehaviourComponent{
		Tag: erutan.Component_NetworkBehaviourComponent_ALL,
	}
	// Default param
	if speed == -1 {
		speed = utils.RandFloats(10, 20)
	}
	obj.Component_SpeedComponent = &erutan.Component_SpeedComponent{
		MoveSpeed: speed,
	}
	obj.Component_PhysicsComponent = &erutan.Component_PhysicsComponent{
		UseGravity: true,
	}
	// Add our obj to the appropriate systems
	for _, system := range m.World.Systems() {
		switch sys := system.(type) {
		case *CollisionSystem:
			sys.Add(*ocObj,
				obj.Component_SpaceComponent,
				obj.Component_BehaviourTypeComponent,
				obj.Component_PhysicsComponent)
		case *HerbivorousSystem:
			sys.Add(*ocObj,
				obj.Component_SpaceComponent,
				obj.Target,
				obj.Component_HealthComponent,
				obj.Component_SpeedComponent)
		case *NetworkSystem:
			sys.Add(*ocObj,
				[]*erutan.Component{
					{Type: &erutan.Component_Space{Space: obj.Component_SpaceComponent}},
					{Type: &erutan.Component_Render{Render: obj.Component_RenderComponent}},
					{Type: &erutan.Component_Health{Health: obj.Component_HealthComponent}},
					{Type: &erutan.Component_Speed{Speed: obj.Component_SpeedComponent}},
					{Type: &erutan.Component_NetworkBehaviour{NetworkBehaviour: obj.Component_NetworkBehaviourComponent}},
				})
		}
	}
}
