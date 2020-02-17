package game

import (
	"github.com/user/erutan/ecs"
	erutan "github.com/user/erutan/protos/realtime"
	"github.com/user/erutan/utils"
)

type collisionEntity struct {
	*ecs.BasicEntity
	*erutan.Component_SpaceComponent
	*erutan.Component_BehaviourTypeComponent
}

// CollisionSystem is a system that detects collisions between entities, sends a message if collisions
// are detected, and updates their SpaceComponent so entities cannot pass through Solids.
type CollisionSystem struct {
	entities []collisionEntity
}

// Add adds an entity to the CollisionSystem. To be added, the entity has to have a basic and space component.
func (c *CollisionSystem) Add(basic *ecs.BasicEntity,
	space *erutan.Component_SpaceComponent,
	behaviourType *erutan.Component_BehaviourTypeComponent) {
	c.entities = append(c.entities, collisionEntity{basic, space, behaviourType})
}

// Remove removes an entity from the CollisionSystem.
func (c *CollisionSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range c.entities {
		if e.BasicEntity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		c.entities = append(c.entities[:delete], c.entities[delete+1:]...)
	}
}

// Update checks the entities for collision with eachother. Only Main entities are check for collision explicitly.
// If one of the entities are solid, the SpaceComponent is adjusted so that the other entities don't pass through it.
func (c *CollisionSystem) Update(dt float64) {
	// O(n²)
	for i := 0; i < len(c.entities); i++ {
		for j := i + 1; j < len(c.entities); j++ {
			// Overlap cube only collider atm
			if Overlap(c.entities[i].Component_SpaceComponent, c.entities[j].Component_SpaceComponent) { //utils.Distance(*c.entities[i].Position, *c.entities[j].Position) < 1 {
				// Collide
				//utils.DebugLogf("a: %v, b: %v", c.entities[i].ID(), c.entities[j].ID())
				ManagerInstance.Watch.Notify(utils.Event{Value: EntitiesCollided{a: c.entities[i], b: c.entities[j], dt: dt}})
			}
		}
	}
}

type EntitiesCollided struct {
	a  collisionEntity
	b  collisionEntity
	dt float64
}

type AABB struct {
	minX float64
	maxX float64
	minY float64
	maxY float64
	minZ float64
	maxZ float64
}

func GetAABB(sc *erutan.Component_SpaceComponent) AABB {
	return AABB{minX: sc.Position.X - sc.Scale.X/2,
		maxX: sc.Position.X + sc.Scale.X/2,
		minY: sc.Position.Y - sc.Scale.Y/2,
		maxY: sc.Position.Y + sc.Scale.Y/2,
		minZ: sc.Position.Z - sc.Scale.Z/2,
		maxZ: sc.Position.Z + sc.Scale.Z/2,
	}
}

func Overlap(scA *erutan.Component_SpaceComponent, scB *erutan.Component_SpaceComponent) bool {
	a := GetAABB(scA)
	b := GetAABB(scB)
	return (a.minX <= b.maxX && a.maxX >= b.minX) &&
		(a.minY <= b.maxY && a.maxY >= b.minY) &&
		(a.minZ <= b.maxZ && a.maxZ >= b.minZ)
}
