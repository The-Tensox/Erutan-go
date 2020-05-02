package game

import (
	"github.com/The-Tensox/Erutan-go/internal/cfg"
	"github.com/The-Tensox/Erutan-go/internal/mon"
	"github.com/The-Tensox/Erutan-go/internal/obs"
	"github.com/The-Tensox/Erutan-go/internal/utils"
	erutan "github.com/The-Tensox/Erutan-go/protobuf"
	"github.com/The-Tensox/octree"
	"github.com/The-Tensox/protometry"
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
	return &CollisionSystem{objects: *octree.NewOctree(protometry.NewBoxOfSize(0, 0, 0, cfg.Global.Logic.GroundSize*1000))}
}

func (c *CollisionSystem) Priority() int {
	return 0
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
	if dt > cfg.Global.FramesPerSecond*1 { // FIXME: quick hack, throttle to avoid spam client when connecting (makign crash)
		c.objects.Range(func(o *octree.Object) bool {
			if co, ok := o.Data.(*collisionObject); ok {
				if co.UseGravity {
					newPosition := co.Position
					newPosition.Y = newPosition.Y - 1*dt
					ManagerInstance.Watch.NotifyAll(obs.Event{Value: obs.PhysicsUpdateRequest{Object: *o, NewPosition: *newPosition, Dt: dt}})
					//utils.DebugLogf("%v %v", newPosition, o.Bounds.GetCenter())
				}
				//	// TODO: mass -> heavier fall faster ...
			}
			return true
		})
	}

}

// PhysicsUpdate will check collisions with new space and update accordingly
func (c *CollisionSystem) PhysicsUpdate(object octree.Object, newPosition protometry.Vector3, dt float64) {
	// We need to find the current Object in collisionSystem's Octree
	// FIXME: redundant with getcolliding
	objectCastedToCollisionObject := Find(c.objects, object)
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
				Me:          objectCastedToCollisionObject,
				NewPosition: newPosition,
				Other:       nil, // No collision here !
				Dt:          dt,
			},
		})
		return
	}

	for _, o := range objectsCollided {
		// Ignore self-collision
		if !o.Equal(*objectCastedToCollisionObject) {
			//utils.DebugLogf("collision between id:%v, pos:%v and \nid:%v, pos:%v",
			//	objectCastedToCollisionObject.ID(),
			//	objectCastedToCollisionObject.Bounds.GetCenter(),
			//	o.ID(),
			//	o.Bounds.GetCenter())
			mon.CollisionCounter.Inc()
			// TODO: maybe we should also apply translation to collided object depending on some physical stuff
			// TODO: OK, atm lets just do this: no collision = ok u can move, collision = don't move
			// Somehow a computation of the place it should be after collision: FIXME
			//newMin := objectCastedToCollisionObject.Bounds.Min.Minus(*o.Bounds.Min)
			//newMax := objectCastedToCollisionObject.Bounds.Max.Minus(*o.Bounds.Max)
			objectCastedToCollisionObject.Bounds.Min = objectCastedToCollisionObject.Bounds.Min.Lerp(o.Bounds.Min, 0.5)
			objectCastedToCollisionObject.Bounds.Max = objectCastedToCollisionObject.Bounds.Max.Lerp(o.Bounds.Max, 0.5)

			// Kinda redundant, is it ok to use at two places position ? (octree object + ecs component synced)
			//co := objectCastedToCollisionObject.Data.(*collisionObject)
			//co2 := o.Data.(*collisionObject)
			//co.UseGravity = false
			//co2.UseGravity = false
			//center := objectCastedToCollisionObject.Bounds.GetCenter()
			//co.Position = &center
			//
			//// Notify every collided object with their new positions
			//c.objects.Move(objectCastedToCollisionObject, center.X, center.Y, center.Z)

			// debug raycast thing

			//c1 := objectCastedToCollisionObject.Bounds.GetCenter()
			//c2 := o.Bounds.GetCenter()
			// Middle point between two objects center ?
			//possibleNaivePositionAfterCollision := c1.Lerp(&c2, 0.5)
			//ManagerInstance.AddDebug(&newPosition,
			//	*protometry.NewMeshSquareCuboid(1, true),
			//	/**protometry.NewMeshRectangularCuboid(*possibleNaivePositionAfterCollision, *protometry.NewVector3(1, 1, 1)),*/
			//	erutan.Component_RenderComponent_Color{
			//		Red:   1,
			//		Green: 0,
			//		Blue:  1,
			//		Alpha: 0.7,
			//	})
			ManagerInstance.Watch.NotifyAll(obs.Event{Value: obs.PhysicsUpdateResponse{Me: &o, Other: objectCastedToCollisionObject, Dt: dt}})
		}
	}
}

func (c *CollisionSystem) Handle(event obs.Event) {
	switch e := event.Value.(type) {
	case obs.PhysicsUpdateRequest:
		c.PhysicsUpdate(e.Object, e.NewPosition, e.Dt)
	case obs.PhysicsUpdateResponse:
		// No collision here
		if e.Other == nil {
			me := Find(c.objects, *e.Me)
			if me == nil {
				utils.DebugLogf("Unable to find %v in system %T", e.Me.ID(), c)
				return
			}
			asCo := me.Data.(*collisionObject)
			*asCo.Position = e.NewPosition
			// Need to reinsert in the octree
			if !c.objects.Move(me, e.NewPosition.X, e.NewPosition.Y, e.NewPosition.Z) {
				utils.DebugLogf("Failed to move %v", me)
			}
			// Over
			return
		}
	}
}
