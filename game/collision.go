package game

import (
	"github.com/user/erutan/ecs"
	erutan "github.com/user/erutan/protos/realtime"
	"github.com/user/erutan/utils"
)

type collisionEntity struct {
	*ecs.BasicEntity
	*erutan.Component_SpaceComponent
}

// CollisionSystem is a system that detects collisions between entities, sends a message if collisions
// are detected, and updates their SpaceComponent so entities cannot pass through Solids.
type CollisionSystem struct {
	entities []collisionEntity
}

// Add adds an entity to the CollisionSystem. To be added, the entity has to have a basic and space component.
func (c *CollisionSystem) Add(basic *ecs.BasicEntity, space *erutan.Component_SpaceComponent) {
	c.entities = append(c.entities, collisionEntity{basic, space})
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
	for i1, e1 := range c.entities {
		for i2, e2 := range c.entities {
			if i1 == i2 {
				continue // with other entities, because we won't collide with ourselves
			}
			// Naïve collision distance < 1
			if utils.Distance(*e1.Position, *e2.Position) < 1 {
				// Collide
			}
		}
	}
}
