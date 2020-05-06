package game

import (
	"github.com/The-Tensox/Erutan-go/internal/cfg"
	"github.com/The-Tensox/Erutan-go/internal/utils"
	erutan "github.com/The-Tensox/Erutan-go/protobuf"
	"github.com/The-Tensox/octree"
	"github.com/The-Tensox/protometry"
)

// Some Manager functions to add objects to systems

// TODO: bring together common building-blocks of these stuffs, it's too redundant

// AddDebug create a debug object that will only be seen by clients with debug settings
func (m *Manager) AddDebug(position *protometry.Vector3, mesh protometry.Mesh, color erutan.Component_RenderComponent_Color) {
	obj := BasicObject{}
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
	// Add our object to the appropriate systems
	for _, system := range m.World.Systems() {
		switch sys := system.(type) {
		case *NetworkSystem:
			sys.Add(*ocObj.Clone(), // We want all systems to have their own local state, no pointers in-between
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
	obj := BasicObject{}
	ocObj := octree.NewObjectCube(nil, position.X, position.Y, position.Z, sideLength)

	obj.Component_SpaceComponent = &erutan.Component_SpaceComponent{
		Position: position,
		Rotation: protometry.NewQuaternion(0, 0, 0, 0),
		Scale:    protometry.NewVector3(1, 1, 1),
	}
	var c []*erutan.Component_RenderComponent_Color
	mesh := protometry.NewMeshSquareCuboid(sideLength, true)
	for range mesh.Vertices {
		c = append(c, &erutan.Component_RenderComponent_Color{
			Red:   0,
			Green: float32(utils.RandFloats(0, 3)),
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
	// Add our object to the appropriate systems
	for _, system := range m.World.Systems() {
		switch sys := system.(type) {
		case *CollisionSystem:
			sys.Add(*ocObj.Clone(),
				obj.Component_SpaceComponent,
				obj.Component_BehaviourTypeComponent,
				obj.Component_PhysicsComponent)
		case *NetworkSystem:
			sys.Add(*ocObj.Clone(),
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
	obj := BasicObject{}
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
	// Add our object to the appropriate systems
	for _, system := range m.World.Systems() {
		switch sys := system.(type) {
		case *CollisionSystem:
			sys.Add(*ocObj.Clone(),
				obj.Component_SpaceComponent,
				obj.Component_BehaviourTypeComponent,
				obj.Component_PhysicsComponent)
		case *EatableSystem:
			sys.Add(*ocObj.Clone(),
				obj.Component_SpaceComponent)
		case *NetworkSystem:
			sys.Add(*ocObj.Clone(),
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
	ocObj := octree.NewObjectCube(nil, position.X, position.Y, position.Z, scale.X) // TODO: Only cubes handled atm


	obj.Component_HealthComponent = &erutan.Component_HealthComponent{Life: cfg.Global.Logic.Herbivorous.Life}
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
		Tag: erutan.Component_NetworkBehaviourComponent_ALL, // TODO: handle owner, who spawned this ?
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
			//utils.DebugLogf("CollisionSystem add %d", ocObj.ID())
			sys.Add(*ocObj.Clone(),
				obj.Component_SpaceComponent,
				obj.Component_BehaviourTypeComponent,
				obj.Component_PhysicsComponent)
		case *HerbivorousSystem:
			sys.Add(*ocObj.Clone(),
				obj.Component_SpaceComponent,
				obj.Target,
				obj.Component_HealthComponent,
				obj.Component_SpeedComponent)
		case *NetworkSystem:
			//utils.DebugLogf("NetworkSystem add %d", ocObj.ID())
			sys.Add(*ocObj.Clone(),
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

// AddPlayer create a player object
func (m *Manager) AddPlayer(position *protometry.Vector3, token string) (uint64, BasicObject) {
	obj := BasicObject{}
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
			Red:   float32(utils.RandFloats(0, 1)),
			Green: float32(utils.RandFloats(0, 1)),
			Blue:  float32(utils.RandFloats(0, 1)),
			Alpha: 1,
		})
	}
	obj.Component_RenderComponent = &erutan.Component_RenderComponent{
		Mesh:   mesh,
		Colors: c,
	}
	obj.Component_BehaviourTypeComponent = &erutan.Component_BehaviourTypeComponent{
		Tag: erutan.Component_BehaviourTypeComponent_PLAYER,
	}
	obj.Component_NetworkBehaviourComponent = &erutan.Component_NetworkBehaviourComponent{
		Tag: erutan.Component_NetworkBehaviourComponent_ALL,
		OwnerToken: token, // Owned by this player obviously :)
	}
	obj.Component_PhysicsComponent = &erutan.Component_PhysicsComponent{
		UseGravity: false,
	}
	// Add our obj to the appropriate systems
	for _, system := range m.World.Systems() {
		switch sys := system.(type) {
		case *CollisionSystem:
			sys.Add(*ocObj.Clone(),
				obj.Component_SpaceComponent,
				obj.Component_BehaviourTypeComponent,
				obj.Component_PhysicsComponent)
		case *NetworkSystem:
			sys.Add(*ocObj.Clone(),
				[]*erutan.Component{
					{Type: &erutan.Component_Space{Space: obj.Component_SpaceComponent}},
					{Type: &erutan.Component_Render{Render: obj.Component_RenderComponent}},
					{Type: &erutan.Component_NetworkBehaviour{NetworkBehaviour: obj.Component_NetworkBehaviourComponent}},
				})
		}
	}
	return ocObj.ID(), obj
}
