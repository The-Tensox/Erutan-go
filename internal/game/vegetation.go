package game

import (
	"github.com/The-Tensox/erutan/internal/cfg"
	"github.com/The-Tensox/erutan/internal/obs"
	"github.com/The-Tensox/erutan/internal/utils"
	erutan "github.com/The-Tensox/erutan/protobuf"
	"github.com/The-Tensox/octree"
	"github.com/The-Tensox/protometry"
)


type AnyObject struct {
	*erutan.Component_SpaceComponent
	*erutan.Component_RenderComponent
	*erutan.Component_BehaviourTypeComponent
	*erutan.Component_PhysicsComponent
	*erutan.Component_NetworkBehaviourComponent
}

type eatableObject struct {
	*erutan.Component_SpaceComponent
}

type EatableSystem struct {
	objects octree.Octree
}

func NewEatableSystem() *EatableSystem {
	return &EatableSystem{objects: *octree.NewOctree(protometry.NewBoxOfSize(0, 0, 0,
		cfg.Global.Logic.GroundSize*1000))}
}

func (e *EatableSystem) Add(object octree.Object,
	space *erutan.Component_SpaceComponent) {
	eo := &eatableObject{space}
	object.Data = eo
	if !e.objects.Insert(object) {
		utils.DebugLogf("Failed to insert %v", object)
	}
}

// Remove removes the Object from the System. This is what most Remove methods will look like
func (e *EatableSystem) Remove(object octree.Object) {
	if !e.objects.Remove(object) {
		utils.DebugLogf("Failed to remove")
	}
}

func (e *EatableSystem) Update(dt float64) {
}

func (e *EatableSystem) Handle(event obs.Event) {
	switch u := event.Value.(type) {
	case obs.OnPhysicsUpdateResponse:
		// No collision here
		if u.Other == nil {
			me := Find(e.objects, *u.Me)
			if me == nil {
				//utils.DebugLogf("Unable to find %v in system %T", u.Me.ID(), u)
				return
			}
			asEo := me.Data.(*eatableObject)
			*asEo.Position = u.NewPosition
			// Need to reinsert in the octree
			if !e.objects.Move(me, u.NewPosition.X, u.NewPosition.Y, u.NewPosition.Z) {
				utils.DebugLogf("Failed to move %v", me)
			}
			// Over
			return
		}

		me := u.Me.Data.(*collisionObject)
		other := u.Other.Data.(*collisionObject)
		// If an animal collided with me
		// TODO: FIXME
		if me.Tag == erutan.Component_BehaviourTypeComponent_ANIMAL &&
			other.Tag == erutan.Component_BehaviourTypeComponent_VEGETATION {
			// Teleport somewhere else
			p := protometry.RandomCirclePoint(0, 0, cfg.Global.Logic.GroundSize)
			ManagerInstance.Watch.NotifyAll(obs.Event{Value: obs.OnPhysicsUpdateRequest{Object: *u.Other, NewPosition: p, Dt: u.Dt}})
		}

		if other.Tag == erutan.Component_BehaviourTypeComponent_ANIMAL &&
			me.Tag == erutan.Component_BehaviourTypeComponent_VEGETATION {
			// Teleport somewhere else
			p := protometry.RandomCirclePoint(0, 0, cfg.Global.Logic.GroundSize)
			ManagerInstance.Watch.NotifyAll(obs.Event{Value: obs.OnPhysicsUpdateRequest{Object: *u.Me, NewPosition: p, Dt: u.Dt}})
		}
	}
}
