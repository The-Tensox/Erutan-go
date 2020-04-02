package game

import (
	"github.com/The-Tensox/erutan/ecs"
	erutan "github.com/The-Tensox/erutan/protobuf"
	"github.com/The-Tensox/erutan/utils"
	"github.com/The-Tensox/octree"
	"github.com/The-Tensox/protometry"
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
	c.entities = *octree.NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), utils.Config.GroundSize))
}

// Add adds an entity to the CollisionSystem. To be added, the entity has to have a basic and space component.
func (c *CollisionSystem) Add(basic *ecs.BasicEntity,
	space *erutan.Component_SpaceComponent,
	behaviourType *erutan.Component_BehaviourTypeComponent,
	physics *erutan.Component_PhysicsComponent) {
	e := collisionEntity{basic, space, behaviourType, physics}
	c.entities.Insert(*octree.NewObjectCube(e, e.Position.Get(0), e.Position.Get(1), e.Position.Get(2), 1))
}

// Remove removes an entity from the CollisionSystem.
func (c *CollisionSystem) Remove(basic ecs.BasicEntity) {
	// TODO: remove with arg interface data instead
	/*
		for _, entity := range c.entities.ElementsIn(vector.GetBoxOfSize(protometry.VectorN{X: 0, Y: 0, Z: 0}, utils.Config.GroundSize)) {
			if val, ok := entity.(collisionEntity); ok {
				if val.ID() == basic.ID() {
					c.entities.Remove(val)
				}
			}
		}
	*/
}

// Update checks the entities for collision with eachother. Only Main entities are check for collision explicitly.
// If one of the entities are solid, the SpaceComponent is adjusted so that the other entities don't pass through it.
func (c *CollisionSystem) Update(dt float64) {
	// Gravity, checking if there is an object below, otherwise we fall ! (inefficient)

	elements := c.entities.GetColliding(*protometry.NewBoxOfSize(*protometry.NewVectorN(0, 0, 0), utils.Config.GroundSize))

	for _, e := range elements {
		if a, ok := e.Data.(collisionEntity); ok {
			origin := a.Position

			// The raycast is thrown out from just below the object
			origin.Set(1, origin.Get(1)-(e.Bounds.GetSize()/2)+0.1)
			destination := origin.Sub(*protometry.NewVectorN(0, 1, 0))
			b := *protometry.NewBox(*protometry.Concatenate(origin, destination))
			utils.DebugLogf("yo %v", b)

			// Only fall if using gravity and nothing is below
			if a.UseGravity && len(c.entities.GetColliding(b)) == 0 {
				a.Position.Set(1, a.Position.Get(1)-1*dt) // TODO: mass -> heavier fall faster ...
			}
		}
	}

}

// PhysicsUpdate will check collisions with new space and update accordingly
func (c *CollisionSystem) PhysicsUpdate(id uint64, newSc erutan.Component_SpaceComponent, dt float64) {
	// elements := c.entities.GetColliding(*protometry.NewBoxOfSize(*protometry.NewVectorN(0, 0, 0), utils.Config.GroundSize))

	// for i := range elements {
	// 	// Find the element requesting physics update
	// 	if a, ok := elements[i].Data.(collisionEntity); ok && a.ID() == id {
	// 		iBox := elements[i].Bounds
	// 		for j := 0; j < len(elements); j++ {
	// 			// Ignore self collision
	// 			if b, ok := elements[j].Data.(collisionEntity); ok && b.ID() != id {
	// 				jBox := elements[j].Bounds
	// 				if iBox.Intersects(&jBox) {
	// 					/*
	// 						translation := vector.MinimumTranslation(iBox, jBox)
	// 						if translation.X > c.entities[j].Scale.X/2 || translation.Y > c.entities[j].Scale.Y/2 || translation.Z > c.entities[j].Scale.Z/2 {
	// 							translation = vector.Div(translation, 2)
	// 							*c.entities[j].Position = vector.Add(*c.entities[j].Position, translation)
	// 						}
	// 					*/
	// 					ManagerInstance.Watch.Notify(utils.Event{Value: EntitiesCollided{a: a, b: b, dt: dt}})
	// 				}
	// 			}
	// 		}
	// 		*a.Component_SpaceComponent = newSc
	// 	}
	// }
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
