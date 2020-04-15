package game

import (
	"github.com/The-Tensox/erutan/cfg"
	erutan "github.com/The-Tensox/erutan/protobuf"
	"github.com/The-Tensox/erutan/utils"
	"github.com/The-Tensox/octree"
	"github.com/The-Tensox/protometry"
)

// AddLife set health component life, clip it and return true if entity is dead
func AddLife(id uint64, o octree.Object, h *erutan.Component_HealthComponent, value float64) bool {
	h.Life += value
	// Clip 0, 100
	/*if h.Life > 100 {
		h.Life = 100
	} else*/if h.Life < 0 {
		h.Life = 0
	}
	if h.Life == 0 {
		ManagerInstance.World.RemoveObject(o)
		return true
	}
	return false
}

type Herbivorous struct {
	Id uint64
	*erutan.Component_SpaceComponent
	*erutan.Component_HealthComponent
	Target *AnyObject
	*erutan.Component_RenderComponent
	*erutan.Component_BehaviourTypeComponent
	*erutan.Component_SpeedComponent
	*erutan.Component_PhysicsComponent
	*erutan.Component_NetworkBehaviourComponent
}

func (h Herbivorous) ID() uint64 {
	return h.Id
}

type herbivorousObject struct {
	Id uint64
	*erutan.Component_SpaceComponent
	Target *AnyObject
	*erutan.Component_HealthComponent
	*erutan.Component_SpeedComponent
}

type HerbivorousSystem struct {
	objects octree.Octree
}

func NewHerbivorousSystem() *HerbivorousSystem {
	return &HerbivorousSystem{objects: *octree.NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(),
		cfg.Global.Logic.GroundSize*1000))}
}

func (h *HerbivorousSystem) Add(id uint64,
	space *erutan.Component_SpaceComponent,
	target *AnyObject,
	health *erutan.Component_HealthComponent,
	speed *erutan.Component_SpeedComponent) {
	ho := herbivorousObject{id, space, target,
		health, speed}
	o := octree.NewObjectCube(ho, ho.Position.Get(0), ho.Position.Get(1), ho.Position.Get(2), 1)
	if !h.objects.Insert(*o) {
		utils.DebugLogf("Failed to insert %v", o.ToString())
	}
}

// Remove removes the Object from the System. This is what most Remove methods will look like
func (h *HerbivorousSystem) Remove(o octree.Object) {
	h.objects.Remove(o)
}

func (h *HerbivorousSystem) Update(dt float64) {
	objects := h.objects.GetObjects()
	for indexObject, object := range objects {
		if ho, ok := object.Data.(herbivorousObject); ok {

			//volume := ho.Component_SpaceComponent.Scale.Get(0) * ho.Component_SpaceComponent.Scale.Get(1) * ho.Component_SpaceComponent.Scale.Get(2)
			//
			//// Every animal lose life proportional to deltatime, volume and speed
			//// So bigger and faster animals need more food
			//if AddLife(ho.Id, object, ho.Component_HealthComponent, -3*dt*volume*(ho.MoveSpeed/100)) {
			//	// Dead
			//}

			// If I don't have a target, let's find one
			if ho.Target == nil {
				h.findTarget(indexObject, &ho)
			}
			if ho.Target == nil {
				continue // There is no target / animals
			}

			distance := ho.Position.Distance(*ho.Target.Position)
			// TODO: CHECK AGAIN THIS OPERATION ...
			newPos := ho.Position.Plus(*ho.Target.Position.Minus(*ho.Position).Scale(distance).Div(dt * ho.MoveSpeed))
			newSc := *ho.Component_SpaceComponent
			newSc.Position = newPos

			//entity.Component_SpaceComponent.Update(newSc)
			ManagerInstance.Watch.NotifyAll(utils.Event{Value: utils.ObjectPhysicsUpdated{Object: &object, NewSc: newSc, Dt: dt}})
			//entity.Position = &newPos
		}
	}
}

func (h *HerbivorousSystem) findTarget(indexEntity int, ho *herbivorousObject) {
	// Super brute-force inefficient implementations :)

	// Reproduction mode
	if ho.Life > 80 {
		objects := h.objects.GetObjects()
		for j := indexEntity + 1; j < len(objects); j++ {
			if otherHo, ok := objects[j].Data.(herbivorousObject); ok {
				if ho.Life > 80 {
					newTarget := &AnyObject{}
					newTarget.Id = otherHo.Id
					newTarget.Component_SpaceComponent = otherHo.Component_SpaceComponent
					ho.Target = newTarget

					newTargetTwo := &AnyObject{}
					newTargetTwo.Id = ho.Id
					newTargetTwo.Component_SpaceComponent = ho.Component_SpaceComponent
					otherHo.Target = newTargetTwo
					// Found a target (another animal)
					return
				}
			}
		}
	}
	// Eating mode
	for _, e := range ManagerInstance.World.Systems() {
		if f, ok := e.(*EatableSystem); ok {
			// Currently look for eatable on the whole map
			eatables := f.objects.GetObjects()
			// Is there any eatable on the map?
			if len(eatables) > 0 {
				//var minPosition *octree.Object
				//for _, eatable := range eatables {
				//	if eo, ok := eatable.Data.(*eatableObject); ok {
				//		mp := minPosition.Data.(eatableObject)
				//		if minPosition == nil || ho.Position.Distance(*eo.Position) < ho.Position.Distance(*mp.Position) {
				//			minPosition = &eatable
				//		}
				//	}
				//}
				minPosition := eatables[0].Data.(eatableObject)
				newTarget := &AnyObject{}
				newTarget.Component_SpaceComponent = minPosition.Component_SpaceComponent
				newTarget.Id = minPosition.Id
				ho.Target = newTarget
			}
		}
	}
}

func (h *HerbivorousSystem) Handle(event utils.Event) {
	switch event.Value.(type) {
	case utils.ObjectsCollided:
	}
	//switch e := event.Value.(type) {
	//case EntitiesCollided:
	//	a := h.objects.GetColliding(*protometry.NewBoxOfSize(e.a.Bounds.Center, 1))[0]
	//	b := h.objects.GetColliding(*protometry.NewBoxOfSize(e.b.Bounds.Center, 1))[0]
	//	// TODO: clean this ugly as hell function
	//
	//	if e.a.BehaviourType == erutan.Component_BehaviourTypeComponent_VEGETATION &&
	//		e.b.BehaviourType == erutan.Component_BehaviourTypeComponent_ANIMAL {
	//		AddLife(e.b.Id, b, b.Component_HealthComponent, 40)
	//
	//		// Reset target for everyone that had this target
	//		for _, e := range h.entities {
	//			if (e.Target != nil && a != nil) && (e.Target.ID() == a.ID()) {
	//				e.Target = nil
	//			}
	//		}
	//	} else if e.b.BehaviourType == erutan.Component_BehaviourTypeComponent_VEGETATION &&
	//		e.a.BehaviourType == erutan.Component_BehaviourTypeComponent_ANIMAL {
	//		AddLife(*a.BasicEntity, a.Component_HealthComponent, 40)
	//
	//		// Reset target for everyone that had this target
	//		for _, e := range h.entities {
	//			if (e.Target != nil && b != nil) && (e.Target.ID() == b.ID()) {
	//				e.Target = nil
	//			}
	//		}
	//	} else { // Both are animals
	//		if (a != nil && b != nil) && (a.Life > 80 && b.Life > 80) {
	//			if a.Target != nil {
	//				a.Target = nil
	//			}
	//			if b.Target != nil {
	//				b.Target = nil
	//			}
	//			AddLife(*a.BasicEntity, a.Component_HealthComponent, -50)
	//			AddLife(*b.BasicEntity, b.Component_HealthComponent, -50)
	//			speed := ((a.MoveSpeed + b.MoveSpeed) / 2) * utils.RandFloats(0.5, 1.5)
	//			scale := a.Scale.Plus(*b.Scale).Div(2).Scale(utils.RandFloats(0.5, 1.5))
	//
	//			// Clipping scale ... TODO: add min & max scale somewhere
	//			clip := func(val float64, min float64, max float64) float64 {
	//				if val < min {
	//					return min
	//				} else if val > max {
	//					return max
	//				} else {
	//					return val
	//				}
	//			}
	//			scale.Set(0, clip(scale.Get(0), 0.1, 5))
	//			scale.Set(1, clip(scale.Get(1), 0.1, 5))
	//			scale.Set(2, clip(scale.Get(2), 0.1, 5))
	//			speed = clip(speed, 5, 80)
	//			position := a.Position
	//			position.Set(1, scale.Get(1)) // To stay above ground
	//			//utils.DebugLogf("Scale: %v", scale)
	//			ManagerInstance.AddHerbivorous(position, scale, speed)
	//		}
	//	}
	//}
}
