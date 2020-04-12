package game

import (
	erutan "github.com/The-Tensox/erutan/protobuf"
	"github.com/The-Tensox/erutan/utils"
	"github.com/The-Tensox/octree"
	"github.com/The-Tensox/protometry"
	"time"
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
		utils.Config.GroundSize*100))}
}

// Add adds an entity to the CollisionSystem. To be added, the entity has to have a basic and space component.
func (c *CollisionSystem) Add(id uint64,
	size float64,
	space *erutan.Component_SpaceComponent,
	behaviourType *erutan.Component_BehaviourTypeComponent,
	physics *erutan.Component_PhysicsComponent) {
	co := collisionObject{id, space, behaviourType, physics}
	o := octree.NewObjectCube(co, co.Position.Get(0), co.Position.Get(1), co.Position.Get(2), size)
	if !c.objects.Insert(*o) {
		utils.ServerLogf(time.Now(), "Failed to insert %v", o.ToString())
	}
}

// Remove removes an entity from the CollisionSystem.
func (c *CollisionSystem) Remove(object octree.Object) {
	c.objects.Remove(object)
}

// Update checks the entities for collision with eachother. Only Main entities are check for collision explicitly.
// If one of the entities are solid, the SpaceComponent is adjusted so that the other entities don't pass through it.
func (c *CollisionSystem) Update(dt float64) {
	// TODO: instead every entity handle it's own gravity ?
	// Gravity, checking if there is an object below, otherwise we fall ! (inefficient)
	objects := c.objects.GetColliding(*protometry.NewBoxOfSize(*protometry.NewVectorN(0, 0, 0), utils.Config.GroundSize*10))
	//utils.DebugLogf("len %v", len(objects))
	for _, o := range objects {
		if co, ok := o.Data.(collisionObject); ok {
			min := o.Bounds.GetMin()
			// Get collision under the object
			b := protometry.Box{ // TODO: use object size instead
				Center: *protometry.NewVectorN(o.Bounds.Center.Get(0), min.Get(1)-0.25, o.Bounds.Center.Get(2)),
				Extents: *protometry.NewVectorN(0, 0.24, 0),
			}
			//utils.DebugLogf("b : %v\n%v", o.Bounds.ToString(), b.ToString())
			// Only fall if using gravity and nothing is below
			if co.UseGravity && len(c.objects.GetColliding(b)) == 0 {
				//utils.DebugLogf("FALL")
				//_ = co.Position.Set(1, co.Position.Get(1)-1*dt) // TODO: mass -> heavier fall faster ...
				newSc := *co.Component_SpaceComponent
				_ = newSc.Position.Set(1, co.Position.Get(1)-10*dt)
				//utils.DebugLogf("old pos: %v\nnew pos: %v", co.Position.ToString(), newSc.Position.ToString())
				ManagerInstance.Watch.Notify(utils.Event{Value: ObjectPhysicsUpdated{object: &o, newSc: newSc, dt: dt}})
			}
		}
	}

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
		if o.Data == object.Data { // Could instead compare ids
			objectCastedToCollisionObject = &o
		}
	}

	// This object hasn't been added to collisionSystem or has been removed, abort
	if objectCastedToCollisionObject == nil {
		return
	}
	for _, o := range objectsCollided {
		// Ignore self-collision
		if o.Data != objectCastedToCollisionObject.Data {
			//utils.DebugLogf("collision between %v and\n%v", objectCastedToCollisionObject.ToString(), o.ToString())
			// Notify every collided object
			ManagerInstance.Watch.Notify(utils.Event{Value: ObjectsCollided{a: &o, b: objectCastedToCollisionObject, dt: dt}})
		}
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
