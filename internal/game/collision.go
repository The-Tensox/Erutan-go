package game

import (
	"github.com/The-Tensox/Erutan-go/internal/cfg"
	"github.com/The-Tensox/Erutan-go/internal/mon"
	"github.com/The-Tensox/Erutan-go/internal/obs"
	"github.com/The-Tensox/Erutan-go/internal/utils"
	erutan "github.com/The-Tensox/Erutan-go/protobuf"
	"github.com/The-Tensox/octree"
	"github.com/The-Tensox/protometry"
	"math"
)

type collisionObject struct {
	*erutan.Component_SpaceComponent
	*erutan.Component_BehaviourTypeComponent
	*erutan.Component_PhysicsComponent
}

// CollisionSystem is a system that handle collisions
type CollisionSystem struct {
	objects octree.Octree
}

func NewCollisionSystem() *CollisionSystem {
	return &CollisionSystem{objects: *octree.NewOctree(protometry.NewBoxOfSize(0, 0, 0, cfg.Global.Logic.OctreeSize))}
}

func (c CollisionSystem) Priority() int {
	return math.MaxInt64
}

// Add adds an object to the CollisionSystem. To be added, the object has to have a basic and space component.
func (c *CollisionSystem) Add(object octree.Object,
	space *erutan.Component_SpaceComponent,
	behaviourType *erutan.Component_BehaviourTypeComponent,
	physics *erutan.Component_PhysicsComponent) {
	co := &collisionObject{space, behaviourType, physics}
	object.Data = co
	if !c.objects.Insert(object) {
		utils.DebugLogf("Failed to insert %v", object)
	} else {
		mon.PhysicalObjectsGauge.Inc()
	}
}

// Remove removes an object from the CollisionSystem.
func (c *CollisionSystem) Remove(object octree.Object) {
	if !c.objects.Remove(object) {
		utils.DebugLogf("Failed to remove %d, data: %T", object.ID(), object.Data)
	}
}

// Update checks the objects for collision with eachother. Only Main objects are check for collision explicitly.
// If one of the objects are solid, the SpaceComponent is adjusted so that the other objects don't pass through it.
func (c *CollisionSystem) Update(dt float64) {
	// Gravity
	return
	if dt > cfg.Global.UpdatesRate*1 { // FIXME: quick hack, throttle to avoid spam client when connecting (makign crash)
		c.objects.Range(func(o *octree.Object) bool {
			if co, ok := o.Data.(*collisionObject); ok {
				if co.UseGravity {
					newPosition := co.Position
					newPosition.Y = newPosition.Y - 1*dt
					ManagerInstance.Watch.NotifyAll(obs.Event{Value: obs.PhysicsUpdateRequest{
						Object: struct {
							octree.Object
							protometry.Vector3
						}{Object: *o, Vector3: *newPosition},
						Dt: dt,
					}})
					//utils.DebugLogf("%v %v", newPosition, o.Bounds.GetCenter())
				}
				//	// TODO: mass -> heavier fall faster ...
			}
			return true
		})
	}

}

// ComputePhysics will check collisions with new space and update accordingly
func (c *CollisionSystem) ComputePhysics(object octree.Object, newPosition protometry.Vector3, dt float64) {
	// We need to find the current Object in collisionSystem's Octree
	// FIXME: redundant with getcolliding
	//objectCastedToCollisionObject := Find(c.objects, object)
	objectCastedToCollisionObject := c.objects.Get(object.ID(), object.Bounds)
	if objectCastedToCollisionObject == nil {
		return
	}
	size := object.Bounds.GetSize() // FIXME: atm assuming cube
	objectsCollided := c.objects.GetColliding(*protometry.NewBoxOfSize(newPosition.X, newPosition.Y, newPosition.Z, size.X))

	// Didn't collide anything or self-collision return
	if len(objectsCollided) == 0 || objectsCollided[0].Equal(object) {
		//utils.DebugLogf("\ncheck collision of id:%v | pos: %v \nwanting to move to %v\nresult: %v collisions",
		//	object.ID(), objectCastedToCollisionObject.Bounds.GetCenter(), newPosition, objectsCollided)
		ManagerInstance.Watch.NotifyAll(obs.Event{
			Value: obs.PhysicsUpdateResponse{
				Objects: []struct {
					octree.Object
					protometry.Vector3
				}{{*objectCastedToCollisionObject.Clone(), newPosition}},
				Dt: dt,
			},
		})
		return
	}

	// Handle all collisions
	for i := range objectsCollided {
		// Ignore self-collision
		if !objectsCollided[i].Equal(*objectCastedToCollisionObject) {
			//utils.DebugLogf("collision between id:%v, pos:%v and \nid:%v, pos:%v",
			//	objectCastedToCollisionObject.ID(),
			//	objectCastedToCollisionObject.Bounds.GetCenter(),
			//	objectsCollided[i].ID(),
			//	objectsCollided[i].Bounds.GetCenter())
			mon.CollisionCounter.Inc()
			// Compute new positions of both objects
			//newPositionMe := objectCastedToCollisionObject.Bounds.GetCenter() // Sure?
			//newPositionOther := objectsCollided[i].Bounds.GetCenter()
			//
			//meCo := objectCastedToCollisionObject.Data.(*collisionObject)
			//otherCo := objectsCollided[i].Data.(*collisionObject)
			//if otherCo.IsKinematic {
			//
			//} else {
			//	//newPositionOther.Add(meCo.)
			//}
			//
			//// Compute new position of the object that moved
			//// Check if it can climb on that object (Y axis)
			//if objectsCollided[i].Bounds.Max.Y-objectCastedToCollisionObject.Bounds.Min.Y > 0.25 {
			//	newPosition.Y+=0.25
			//}
			//objectCastedToCollisionObject.Bounds.Min = objectCastedToCollisionObject.Bounds.Min.Lerp(o.Bounds.Min, 0.5)
			//objectCastedToCollisionObject.Bounds.Max = objectCastedToCollisionObject.Bounds.Max.Lerp(o.Bounds.Max, 0.5)
			// Atm doesn't move in case of collision
			ManagerInstance.Watch.NotifyAll(obs.Event{Value: obs.PhysicsUpdateResponse{
				Objects: []struct {
					octree.Object
					protometry.Vector3
				}{
					{*objectCastedToCollisionObject.Clone(), objectCastedToCollisionObject.Bounds.GetCenter()},
					{*objectsCollided[i].Clone(), objectsCollided[i].Bounds.GetCenter()}, // First object with its new position
				},
				Dt: dt}},
			)
		}
	}
}

func (c *CollisionSystem) Handle(event obs.Event) {
	switch e := event.Value.(type) {
	case obs.PhysicsUpdateRequest:
		c.ComputePhysics(e.Object.Object, e.Object.Vector3, e.Dt)
	case obs.PhysicsUpdateResponse:
		// Update position of every objects, if there was a collision or not
		//if len(e.Objects) == 2 {
		//	utils.DebugLogf("col %v - %v", e.Objects[0], e.Objects[1])
		//}
		for i := range e.Objects {
			//utils.DebugLogf("need to move %v; %v to %v", obj.ID(), obj.Bounds.GetCenter(), obj.Vector3)
			me := c.objects.Get(e.Objects[i].Object.ID(), e.Objects[i].Object.Bounds)
			//me := Find(c.objects, obj.Object)
			if me == nil {
				utils.DebugLogf("Unable to find %v in system %T", e.Objects[i].Object.ID(), c)
				return
			}
			if asCo, ok := me.Data.(*collisionObject); ok {
				*asCo.Position = e.Objects[i].Vector3
			}
			// Need to reinsert in the octree
			if !c.objects.Move(me, e.Objects[i].Vector3.X, e.Objects[i].Vector3.Y, e.Objects[i].Vector3.Z) {
				utils.DebugLogf("Failed to move %v", me)
			} else {
				//utils.DebugLogf("moved %v to %v", me.ID(), me.Bounds.GetCenter())
			}
		}
	}
}
