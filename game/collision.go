package game

import (
	"math"

	"github.com/user/erutan/ecs"
	erutan "github.com/user/erutan/protos/realtime"
	"github.com/user/erutan/utils"
	"github.com/user/erutan/utils/octree"
	"github.com/user/erutan/utils/vector"
)

type collisionEntity struct {
	*ecs.BasicEntity
	*erutan.Component_SpaceComponent
	*erutan.Component_BehaviourTypeComponent
	*erutan.Component_PhysicsComponent
}

// CollisionSystem is a system that detects collisions between entities, sends a message if collisions
// are detected, and updates their SpaceComponent so entities cannot pass through Solids.
type CollisionSystem struct {
	entities octree.Octree
}

// New is the constructor of CollisionSystem
func (c *CollisionSystem) New(w *ecs.World) {
	c.entities = *octree.NewOctree(vector.GetBoxOfSize(erutan.NetVector3{X: 0, Y: 0, Z: 0}, utils.Config.GroundSize))
}

// Add adds an entity to the CollisionSystem. To be added, the entity has to have a basic and space component.
func (c *CollisionSystem) Add(basic *ecs.BasicEntity,
	space *erutan.Component_SpaceComponent,
	behaviourType *erutan.Component_BehaviourTypeComponent,
	physics *erutan.Component_PhysicsComponent) {
	e := collisionEntity{basic, space, behaviourType, physics}
	c.entities.Add(e, *e.Position)
}

// Remove removes an entity from the CollisionSystem.
func (c *CollisionSystem) Remove(basic ecs.BasicEntity) {
	// TODO: take advantage of tree DS to do O(logn) removal instead of this O(n)+?
	for _, entity := range c.entities.ElementsIn(vector.GetBoxOfSize(erutan.NetVector3{X: 0, Y: 0, Z: 0}, utils.Config.GroundSize)) {
		if val, ok := entity.(collisionEntity); ok {
			if val.ID() == basic.ID() {
				c.entities.Remove(val)
			}
		}
	}
}

// RemoveByPosition TODO: just search in this box, may be faster?
func (c *CollisionSystem) RemoveByPosition(position erutan.NetVector3) {
}

// Update checks the entities for collision with eachother. Only Main entities are check for collision explicitly.
// If one of the entities are solid, the SpaceComponent is adjusted so that the other entities don't pass through it.
func (c *CollisionSystem) Update(dt float64) {
	c.entities.Range(func(elements []interface{}) {
		for i := range elements {
			if a, ok := elements[i].(collisionEntity); ok {
				// Gravity, checking if there is an object below, otherwise we fall ! (inefficient)
				if a.UseGravity && c.Raycast(*a.Position, erutan.NetVector3{X: 0, Y: -1, Z: 0}, 0.1) == nil {
					a.Position.Y -= (1 * dt) // TODO: mass -> heavier fall faster ...
				}
				iBox := vector.GetBox(*a.Position, *a.Scale)
				for j := i + 1; j < len(elements); j++ {
					if b, ok := elements[j].(collisionEntity); ok {
						if a.ID() == b.ID() {
							continue
						}
						jBox := vector.GetBox(*b.Position, *b.Scale)
						if iBox.Intersects(&jBox) {
							// TODO: fix all this translation thing
							/*
								translation := vector.MinimumTranslation(iBox, jBox)
								// TODO: if overlap of more than 20% of its volume ...
								if translation.X > c.entities[j].Scale.X/2 || translation.Y > c.entities[j].Scale.Y/2 || translation.Z > c.entities[j].Scale.Z/2 {
									translation = vector.Div(translation, 2)
									*c.entities[j].Position = vector.Add(*c.entities[j].Position, translation)
								}
							*/
							ManagerInstance.Watch.Notify(utils.Event{Value: EntitiesCollided{a: a, b: b, dt: dt}})
						}
					}
				}
			}
		}
	})
}

type EntitiesCollided struct {
	a  collisionEntity
	b  collisionEntity
	dt float64
}

// Raycast will try to find collisionEntity in a direction from an origin with a maximum distance
func (c *CollisionSystem) Raycast(origin erutan.NetVector3, direction erutan.NetVector3, maxDistance float64) *collisionEntity {
	// Default value in go ...
	if maxDistance == -1 {
		maxDistance = math.MaxFloat64
	}
	b := vector.Box{
		Min: origin,
		Max: vector.Mul(direction, maxDistance),
	}
	// TODO: optimization: should just stop at first element met, we don't care for others
	hits := c.entities.ElementsIn(b)
	if len(hits) == 0 {
		return nil
	}

	if hit, ok := hits[0].(collisionEntity); ok {
		return &hit
	}
	return nil
}
