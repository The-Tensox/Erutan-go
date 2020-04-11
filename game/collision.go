package game

import (
	"github.com/The-Tensox/erutan/ecs"
	erutan "github.com/The-Tensox/erutan/protobuf"
	"github.com/The-Tensox/erutan/utils"
	"github.com/The-Tensox/octree"
	"github.com/The-Tensox/protometry"
)

type collisionObject struct {
	Id uint64
	*erutan.Component_SpaceComponent
	*erutan.Component_BehaviourTypeComponent
	*erutan.Component_PhysicsComponent
}

// CollisionSystem is a system that handle collisions
type CollisionSystem struct {
	objects octree.Octree
}

func NewCollisionSystem() *CollisionSystem {
	return &CollisionSystem{objects: *octree.NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(),
		utils.Config.GroundSize))}
}

// New is the constructor of CollisionSystem
func (c *CollisionSystem) New(w *ecs.World) {
	c.objects = *octree.NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), utils.Config.GroundSize))
}

// Add adds an entity to the CollisionSystem. To be added, the entity has to have a basic and space component.
func (c *CollisionSystem) Add(id uint64,
	space *erutan.Component_SpaceComponent,
	behaviourType *erutan.Component_BehaviourTypeComponent,
	physics *erutan.Component_PhysicsComponent) {
	co := collisionObject{id, space, behaviourType, physics}
	o := octree.NewObjectCube(co, co.Position.Get(0), co.Position.Get(1), co.Position.Get(2), 1)
	c.objects.Insert(*o)
}

// Remove removes an entity from the CollisionSystem.
func (c *CollisionSystem) Remove(object octree.Object) {
	c.objects.Remove(object)
}

// Update checks the entities for collision with eachother. Only Main entities are check for collision explicitly.
// If one of the entities are solid, the SpaceComponent is adjusted so that the other entities don't pass through it.
func (c *CollisionSystem) Update(dt float64) {
	// TODO: instead every entity handle it's own gravity ?
	//// Gravity, checking if there is an object below, otherwise we fall ! (inefficient)
	//
	//elements := c.entities.GetColliding(*protometry.NewBoxOfSize(*protometry.NewVectorN(0, 0, 0), utils.Config.GroundSize))
	////utils.DebugLogf("a: %v", c.entities.GetNumberOfObjects())
	//
	//for _, e := range elements {
	//	if a, ok := e.Data.(collisionEntity); ok {
	//		//origin := a.Position //.Dot(*protometry.NewVectorN(1, -(e.Bounds.GetSize()/2)+0.1, 1))
	//
	//		// The raycast is thrown out from just below the object
	//		//destination := origin.Cross(*protometry.NewVectorN(1, -1.1, 1))
	//		b := *protometry.NewBoxMinMax(protometry.Concatenate(*a.Position, *a.Position).Dimensions...)
	//		b.Center = *b.Center.Minus(*protometry.NewVectorN(0, 1, 1))
	//		if a.UseGravity {
	//			utils.DebugLogf("\n\n%v \n%v\n\n", e.Bounds.ToString(), b)
	//
	//			utils.DebugLogf("nb collide %v", (c.entities.GetColliding(b)))
	//		}
	//
	//		// Only fall if using gravity and nothing is below
	//		if a.UseGravity && len(c.entities.GetColliding(b)) == 0 {
	//			utils.DebugLogf("FALL")
	//			a.Position.Set(1, a.Position.Get(1)-1*dt) // TODO: mass -> heavier fall faster ...
	//		}
	//	}
	//}

}

// PhysicsUpdate will check collisions with new space and update accordingly
func (c *CollisionSystem) PhysicsUpdate(object octree.Object, newSc erutan.Component_SpaceComponent, dt float64) {
	objectsCollided := c.objects.GetColliding(object.Bounds)
	// Didn't collide anything, return
	if len(objectsCollided) == 0 {
		return
	}
	var objectCastedToCollisionObject *octree.Object

	// We need to find the current Object in collisionSystem's Octree
	for _, o := range objectsCollided {
		if o.Bounds.Equal(object.Bounds) {
			objectCastedToCollisionObject = &o
		}
	}

	// This object hasn't been added to collisionSystem or has been removed, abort
	if objectCastedToCollisionObject == nil {
		return
	}
	for _, o := range objectsCollided {
		// Notify every collided object
		ManagerInstance.Watch.Notify(utils.Event{Value: ObjectsCollided{a: &o, b: objectCastedToCollisionObject, dt: dt}})
	}
	co := objectCastedToCollisionObject.Data.(collisionObject)
	// TODO: apply translation if collision ...
	co.Position = newSc.Position // ?
	c.objects.Move(objectCastedToCollisionObject, newSc.Position.Dimensions...)
}

func (c *CollisionSystem) NotifyCallback(event utils.Event) {
	switch e := event.Value.(type) {
	case ObjectPhysicsUpdated:
		c.PhysicsUpdate(*e.object, e.newSc, e.dt)
	}
}

type ObjectsCollided struct {
	a  *octree.Object
	b  *octree.Object
	dt float64
}

type ObjectPhysicsUpdated struct {
	object    *octree.Object
	newSc erutan.Component_SpaceComponent
	dt    float64
}
