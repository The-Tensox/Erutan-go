package erutan

import (
	"github.com/The-Tensox/Erutan-go/internal/cfg"
	"github.com/The-Tensox/Erutan-go/internal/log"
	"github.com/The-Tensox/Erutan-go/internal/obs"
	"github.com/The-Tensox/octree"
	"github.com/The-Tensox/protometry"
	"go.uber.org/zap"
	"reflect"
)

// TODO: I think there is some changes to this design to be done, not sure it's clean to mix objects, objects ...
// TODO: or maybe it's ok idk
type BasicObject struct {
	*Component_SpaceComponent
	*Component_RenderComponent
	*Component_BehaviourTypeComponent
	*Component_PhysicsComponent
	*Component_NetworkBehaviourComponent
}

type eatableObject struct {
	*Component_SpaceComponent
}

type EatableSystem struct {
	objects octree.Octree
}

func NewEatableSystem() *EatableSystem {
	return &EatableSystem{objects: *octree.NewOctree(protometry.NewBoxOfSize(0, 0, 0, cfg.Get().Logic.OctreeSize))}
}

func (e *EatableSystem) Add(object octree.Object,
	space *Component_SpaceComponent) {
	eo := &eatableObject{space}
	object.Data = eo
	if !e.objects.Insert(object) {
		log.Zap.Info("Failed to insert", zap.Any("object", object))
	}
}

// Remove removes the Object from the System. This is what most Remove methods will look like
func (e *EatableSystem) Remove(object octree.Object) {
	if !e.objects.Remove(object) {
		log.Zap.Info("Failed to remove", zap.Any("ID", object.ID()), zap.Any("data", reflect.TypeOf(object.Data)))
	}
}

func (e *EatableSystem) Update(_ float64) {
}

func (e *EatableSystem) Handle(event obs.Event) {
	switch u := event.Value.(type) {
	case PhysicsUpdateResponse:
		// No collision here
		if len(u.Objects) == 1 {
			me := e.objects.Get(u.Objects[0].Object.ID(), u.Objects[0].Object.Bounds)
			if me == nil {
				//log.Zap.Info("Unable to find %v in system %T", u.Me.ID(), u)
				return
			}
			if asEo, ok := me.Data.(*eatableObject); ok {
				*asEo.Position = u.Objects[0].Vector3
			}
			// Need to reinsert in the octree
			if !e.objects.Move(me, u.Objects[0].Vector3.X, u.Objects[0].Vector3.Y, u.Objects[0].Vector3.Z) {
				log.Zap.Info("Failed to move", zap.Any("object", me))
			}
		} else if len(u.Objects) == 2 { // Means collision, shouldn't be > 2 imho
			me := u.Objects[0].Data.(*collisionObject)
			other := u.Objects[1].Data.(*collisionObject)
			var newSpotToTeleport *protometry.Vector3
			tries := 0
			for newSpotToTeleport == nil || tries == 20 {
				p := protometry.RandomCirclePoint(0, 0, 0, 50)
				if collisions := e.objects.GetColliding(
					*protometry.NewBoxOfSize(p.X, p.Y, p.Z, u.Objects[0].Bounds.GetSize().Sum()/3)); len(collisions) == 0 {
					newSpotToTeleport = &p
				}
				tries++
			}
			if tries == 20  {
				log.Zap.Info("Couldn't find an empty spot to teleport !!")
				return
			}
			// If an animal collided with me
			// TODO: FIXME
			if me.Tag == Component_BehaviourTypeComponent_ANIMAL &&
				other.Tag == Component_BehaviourTypeComponent_VEGETATION {
				// Teleport somewhere else
				ManagerInstance.NotifyAll(obs.Event{Value: PhysicsUpdateRequest{
					Object: struct{octree.Object;protometry.Vector3}{Object: *u.Objects[1].Object.Clone(),
						Vector3: *newSpotToTeleport},
					Dt: u.Dt}})
			} else if other.Tag == Component_BehaviourTypeComponent_ANIMAL &&
				me.Tag == Component_BehaviourTypeComponent_VEGETATION {
				// Teleport somewhere else
				ManagerInstance.NotifyAll(obs.Event{Value: PhysicsUpdateRequest{
					Object: struct{octree.Object;protometry.Vector3}{Object: *u.Objects[0].Object.Clone(),
						Vector3: *newSpotToTeleport},
					Dt: u.Dt}})
			}
		}
	}
}
