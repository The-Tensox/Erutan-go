package erutan

import (
	"github.com/The-Tensox/Erutan-go/internal/cfg"
	"github.com/The-Tensox/Erutan-go/internal/log"
	"github.com/The-Tensox/Erutan-go/internal/mon"
	"github.com/The-Tensox/Erutan-go/internal/obs"
	"github.com/The-Tensox/octree"
	"github.com/The-Tensox/protometry"
	"go.uber.org/zap"
	"math"
	"reflect"
)

type collisionObject struct {
	*Component_SpaceComponent
	*Component_BehaviourTypeComponent
	*Component_PhysicsComponent
}

// CollisionSystem is a system that handle collisions
type CollisionSystem struct {
	objects octree.Octree
}

func NewCollisionSystem() *CollisionSystem {
	return &CollisionSystem{objects: *octree.NewOctree(protometry.NewBoxOfSize(0, 0, 0, cfg.Get().Logic.OctreeSize))}
}

func (c CollisionSystem) Priority() int {
	return math.MaxInt64
}

// Add adds an object to the CollisionSystem. To be added, the object has to have a basic and space component.
func (c *CollisionSystem) Add(object octree.Object,
	space *Component_SpaceComponent,
	behaviourType *Component_BehaviourTypeComponent,
	physics *Component_PhysicsComponent) {
	co := &collisionObject{space, behaviourType, physics}
	object.Data = co
	if !c.objects.Insert(object) {
		log.Zap.Info("Failed to insert", zap.Any("object", object))
	} else {
		mon.PhysicalObjectsGauge.Inc()
	}
}

// Remove removes an object from the CollisionSystem.
func (c *CollisionSystem) Remove(object octree.Object) {
	if !c.objects.Remove(object) {
		log.Zap.Info("Failed to remove", zap.Any("ID", object.ID()), zap.Any("data", reflect.TypeOf(object.Data)))
	}
}

// Apply gravity
func (c *CollisionSystem) Update(dt float64) {
	// Gravity
	return
	if dt > cfg.Get().UpdatesRate*1 { // FIXME: quick hack, throttle to avoid spam client when connecting (makign crash)
		c.objects.Range(func(o *octree.Object) bool {
			if co, ok := o.Data.(*collisionObject); ok {
				if co.UseGravity {
					newPosition := co.Position
					newPosition.Y = newPosition.Y - 1*dt
					ManagerInstance.NotifyAll(obs.Event{Value: PhysicsUpdateRequest{
						Object: struct {
							octree.Object
							protometry.Vector3
						}{Object: *o, Vector3: *newPosition},
						Dt: dt,
					}})
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
		ManagerInstance.NotifyAll(obs.Event{
			Value: PhysicsUpdateResponse{
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
			mon.CollisionCounter.Inc()
			// Compute new positions of both objects
			newPositionMe := newPosition
			newPositionOther := objectsCollided[i].Bounds.GetCenter()
			// Atm doesn't move in case of collision
			ManagerInstance.NotifyAll(obs.Event{Value: PhysicsUpdateResponse{
				Objects: []struct {
					octree.Object
					protometry.Vector3
				}{
					{*objectCastedToCollisionObject.Clone(), newPositionMe},
					{*objectsCollided[i].Clone(), newPositionOther},
				},
				Dt: dt}},
			)
		}
	}
}

func (c *CollisionSystem) Handle(event obs.Event) {
	switch e := event.Value.(type) {
	case PhysicsUpdateRequest:
		c.ComputePhysics(e.Object.Object, e.Object.Vector3, e.Dt)
	case PhysicsUpdateResponse:
		// Update position of every objects, if there was a collision or not
		for i := range e.Objects {
			me := c.objects.Get(e.Objects[i].Object.ID(), e.Objects[i].Object.Bounds)
			if me == nil {
				log.Zap.Info("Unable to find in system", zap.Uint64("ID", e.Objects[i].Object.ID()))
				return
			}
			if asCo, ok := me.Data.(*collisionObject); ok {
				*asCo.Position = e.Objects[i].Vector3
			}
			// Need to reinsert in the octree
			if !c.objects.Move(me, e.Objects[i].Vector3.X, e.Objects[i].Vector3.Y, e.Objects[i].Vector3.Z) {
				log.Zap.Info("Failed to move", zap.Any("object", me))
			} else {
			}
		}
	}
}
