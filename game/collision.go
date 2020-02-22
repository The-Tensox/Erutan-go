package game

import (
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
	c.entities = *octree.NewOctree(vector.GetBoxOfSize(*erutan.NewNetVector3(0, 0, 0), utils.Config.GroundSize))
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
	for _, entity := range c.entities.ElementsIn(*vector.GetBoxOfSize(erutan.NetVector3{X: 0, Y: 0, Z: 0}, utils.Config.GroundSize)) {
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

}

// PhysicsUpdate will check collisions with new space and update accordingly
func (c *CollisionSystem) PhysicsUpdate(id uint64, newSc erutan.Component_SpaceComponent, dt float64) {
	//c.entities.Range(func(elements []interface{}) {
	elements := c.entities.ElementsIn(*vector.GetBoxOfSize(erutan.NetVector3{X: 0, Y: 0, Z: 0}, utils.Config.GroundSize))
	/*
		if v, ok := elements[rand.Intn(len(elements)-1)].(collisionEntity); ok {
			utils.DebugLogf("%v", v.UseGravity, v.BehaviourType)

		}
		return
	*/
	for i := range elements {
		// Find the element requesting physics update
		if a, ok := elements[i].(collisionEntity); ok && a.ID() == id {
			iBox := *vector.GetBox(*a.Position, *a.Scale)
			// Gravity, checking if there is an object below, otherwise we fall ! (inefficient)
			/*
				origin := *a.Position
				origin.Y -= ((iBox.Size().Y / 2) + 0.1)
				//utils.DebugLogf("raycast origin %v, me %v", origin.Y, a.Position.Y)
				direction := erutan.NetVector3{X: 0, Y: -1, Z: 0}
				if a.UseGravity && c.entities.Raycast(origin, direction, 0.1) == nil {
					a.Position.Y -= (1 * dt) // TODO: mass -> heavier fall faster ...
				}
			*/
			for j := 0; j < len(elements); j++ {
				// Ignore self collision
				if b, ok := elements[j].(collisionEntity); ok && b.ID() != id {
					jBox := *vector.GetBox(*b.Position, *b.Scale)
					if iBox.Intersects(&jBox) {
						/*
							translation := vector.MinimumTranslation(iBox, jBox)
							if translation.X > c.entities[j].Scale.X/2 || translation.Y > c.entities[j].Scale.Y/2 || translation.Z > c.entities[j].Scale.Z/2 {
								translation = vector.Div(translation, 2)
								*c.entities[j].Position = vector.Add(*c.entities[j].Position, translation)
							}
						*/
						ManagerInstance.Watch.Notify(utils.Event{Value: EntitiesCollided{a: a, b: b, dt: dt}})
					}
				}
			}
			*a.Component_SpaceComponent = newSc
		}
	}
	//})
}

func (c *CollisionSystem) NotifyCallback(event utils.Event) {
	switch e := event.Value.(type) {
	case EntityPhysicsUpdated:
		c.PhysicsUpdate(e.id, e.newSc, e.dt)
	}
}

type EntitiesCollided struct {
	a  collisionEntity
	b  collisionEntity
	dt float64
}

type EntityPhysicsUpdated struct {
	id    uint64
	newSc erutan.Component_SpaceComponent
	dt    float64
}
