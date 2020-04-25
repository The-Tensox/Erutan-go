package game

import (
	"github.com/The-Tensox/erutan/internal/cfg"
	"github.com/The-Tensox/erutan/internal/mon"
	"github.com/The-Tensox/erutan/internal/obs"
	"github.com/The-Tensox/erutan/internal/utils"
	erutan "github.com/The-Tensox/erutan/protobuf"
	"github.com/The-Tensox/octree"
	"github.com/The-Tensox/protometry"
)


type Herbivorous struct {
	*erutan.Component_SpaceComponent
	*erutan.Component_HealthComponent
	Target *AnyObject
	*erutan.Component_RenderComponent
	*erutan.Component_BehaviourTypeComponent
	*erutan.Component_SpeedComponent
	*erutan.Component_PhysicsComponent
	*erutan.Component_NetworkBehaviourComponent
}


type herbivorousObject struct {
	*erutan.Component_SpaceComponent
	Target *AnyObject
	*erutan.Component_HealthComponent
	*erutan.Component_SpeedComponent
}

// AddLife set health component life, clip it and return true if entity is dead
func (h *herbivorousObject) addLife(value float64) bool {
	h.Life += value
	// Clip 0, 100
	/*if h.Life > 100 {
		h.Life = 100
	} else*/if h.Life < 0 {
		h.Life = 0
	}
	mon.LifeGauge.Set(h.Life)
	if h.Life == 0 {
		return true
	}
	return false
}

type HerbivorousSystem struct {
	objects octree.Octree
}


func NewHerbivorousSystem() *HerbivorousSystem {
	return &HerbivorousSystem{objects: *octree.NewOctree(protometry.NewBoxOfSize(0, 0, 0,
		cfg.Global.Logic.GroundSize*1000))}
}

func (h *HerbivorousSystem) Priority() int {
	return 2
}

func (h *HerbivorousSystem) Add(object octree.Object,
	space *erutan.Component_SpaceComponent,
	target *AnyObject,
	health *erutan.Component_HealthComponent,
	speed *erutan.Component_SpeedComponent) {
	ho := &herbivorousObject{space, target,
		health, speed}
	object.Data = ho
	if !h.objects.Insert(object) {
		utils.DebugLogf("Failed to insert %v", object)
	}
	mon.SpeedGauge.Set(speed.MoveSpeed)
}

// Remove removes the Object from the System. This is what most Remove methods will look like
func (h *HerbivorousSystem) Remove(object octree.Object) {
	utils.DebugLogf("removed %v", object.ID())
	if !h.objects.Remove(object) {
		utils.DebugLogf( "Failed to remove")
	}
}

func (h *HerbivorousSystem) Update(dt float64) {
	//FIXME: super inefficient, RIP laptops
	h.objects.Range(func(object *octree.Object) bool {
		if ho, ok := object.Data.(*herbivorousObject); ok {

			//volume := ho.Component_SpaceComponent.Scale.X * ho.Component_SpaceComponent.Scale.Y * ho.Component_SpaceComponent.Scale.Z
			//
			////Every animal lose life proportional to deltatime, volume and speed
			////So bigger and faster animals need more food
			//if ho.addLife(*object, -3*dt*volume*(ho.MoveSpeed/100)) {
			//	// Dead
			//}

			// If I don't have a target, let's find one
			if ho.Target == nil {
				h.findTarget(ho)
				//utils.DebugLogf("my target %v", ho.Target)
			}
			if ho.Target != nil { // Can still be nil (no target on the map)
				distance := ho.Position.Distance(*ho.Target.Position)

				// yolo random direction change for stochastic behaviour :D
				//rnd := protometry.RandomCirclePoint(ho.Target.Position.X, ho.Target.Position.Z, 50)
				//newPos := rnd.Minus(*ho.Position)
				newPos := ho.Target.Position.Minus(*ho.Position)
				newPos.Divide(distance)
				newPos.Scale(dt*ho.MoveSpeed)
				newPos.Add(ho.Position)


				//utils.DebugLogf("I'm here %v, I want to move to %v", object.Bounds.GetCenter(), newPos)
				ManagerInstance.Watch.NotifyAll(obs.Event{Value: obs.OnPhysicsUpdateRequest{Object: *object, NewPosition: newPos, Dt: dt}})
			}
		}
		return true
	})
}

func (h *HerbivorousSystem) findTarget(ho *herbivorousObject) {
	//utils.DebugLogf("trying to find a target", ho.Life)
	// Reproduction mode
	if ho.Life > 80 {
		h.objects.Range(func(object *octree.Object) bool {
			if otherHo, ok := object.Data.(*herbivorousObject); ok {
				if otherHo.Life > 80 {
					newTarget := &AnyObject{}
					newTarget.Component_SpaceComponent = otherHo.Component_SpaceComponent
					ho.Target = newTarget

					newTargetTwo := &AnyObject{}
					newTargetTwo.Component_SpaceComponent = ho.Component_SpaceComponent
					otherHo.Target = newTargetTwo
					//utils.DebugLogf("found a juicy animal target")
					// Found a target (another animal)
				}
			}
			return true
		})
	}
	// Eating mode
	for _, e := range ManagerInstance.World.Systems() {
		if e, ok := e.(*EatableSystem); ok {
			// Currently look for eatable on the whole map
			// Is there any eatable on the map?
			e.objects.Range(func(object *octree.Object) bool {
				//var minPosition *octree.Object
				//for _, eatable := range eatables {
				//	if eo, ok := eatable.Data.(*eatableObject); ok {
				//		mp := minPosition.Data.(eatableObject)
				//		if minPosition == nil || ho.Position.Distance(*eo.Position) < ho.Position.Distance(*mp.Position) {
				//			minPosition = &eatable
				//		}
				//	}
				//}
				//utils.DebugLogf("found a good broccoli target")

				minPosition := object.Data.(*eatableObject)
				newTarget := &AnyObject{}
				newTarget.Component_SpaceComponent = minPosition.Component_SpaceComponent
				ho.Target = newTarget
				return false // Taking just the first object
			})
		}
	}
}


func (h *HerbivorousSystem) Handle(event obs.Event) {
	switch e := event.Value.(type) {
	// In the occurrence of this event we want to check if the animal collided
	// with a vegetation or another animal and take appropriate actions
	case obs.OnPhysicsUpdateResponse:
		// No collision here
		if e.Other == nil {
			me := Find(h.objects, *e.Me)
			if me == nil {
				//utils.DebugLogf("Unable to find %v in system %T", e.Me.ID(), h)
				return
			}
			asHo := me.Data.(*herbivorousObject)
			*asHo.Position = e.NewPosition
			// Need to reinsert in the octree
			if !h.objects.Move(me, e.NewPosition.X, e.NewPosition.Y, e.NewPosition.Z) {
				utils.DebugLogf("Failed to move %v", me)
			}
			//utils.DebugLogf("move %v %v", center, asHo.Position)

			// Over
			return
		}

		var meHsObject, otherHsObject *octree.Object
		// We have to retrieve the object in this system
		// Eventually there can be several object in this object bounds (shouldn't happen though if collision are ON + translation)
		for _, o := range h.objects.GetColliding(e.Me.Bounds) {
			if o.Equal(*e.Me) {
				meHsObject = &o
			} else if o.Equal(*e.Other) { // Else if (we shouldn't receive self-collision)
				otherHsObject = &o
			}
		}

		// Both objects are not in herbivorous system
		if meHsObject == nil && otherHsObject == nil {
			return
		}


		meCo := e.Me.Data.(*collisionObject)
		otherCo := e.Other.Data.(*collisionObject)
		var meHo, otherHo *herbivorousObject
		if meHsObject != nil {
			meHo = meHsObject.Data.(*herbivorousObject)
		}
		if otherHsObject != nil {
			otherHo = otherHsObject.Data.(*herbivorousObject)
		}

		//utils.DebugLogf("collision %v %v", meCo, otherCo)
		if meCo.Tag == erutan.Component_BehaviourTypeComponent_VEGETATION &&
			otherCo.Tag == erutan.Component_BehaviourTypeComponent_ANIMAL {
			//utils.DebugLogf("collision meCo veg + otherCo ani")
			// me is a vegetation
			// other is an animal

			if otherHo != nil && otherHsObject != nil {
				//utils.DebugLogf("before %v", otherHo.Life)
				if otherHo.addLife(40) {
					ManagerInstance.World.RemoveObject(*otherHsObject)
				}
				//utils.DebugLogf("after %v", otherHo.Life)
				mon.EatCounter.Inc()
				// Reset target for everyone that had this target
				h.objects.Range(func(e *octree.Object) bool {
					eho := e.Data.(*herbivorousObject)
					if eho != nil && eho.Target == otherHo.Target {
						eho.Target = nil
					}
					return true
				})
			}

		} else if otherCo.Tag == erutan.Component_BehaviourTypeComponent_VEGETATION &&
			meCo.Tag == erutan.Component_BehaviourTypeComponent_ANIMAL {
			//utils.DebugLogf("collision otherCo veg + meCo ani")

			// other is a vegetation
			// me is an animal

			if meHo != nil && meHsObject != nil {
				//utils.DebugLogf("before %v", meHo.Life)
				if meHo.addLife(40) {
					ManagerInstance.World.RemoveObject(*meHsObject)
				}
				//utils.DebugLogf("after %v", meHo.Life)
				mon.EatCounter.Inc()
				// Reset target for everyone that had this target
				h.objects.Range(func(e *octree.Object) bool {
					eho := e.Data.(*herbivorousObject)
					if eho != nil && eho.Target == meHo.Target {
						eho.Target = nil
					}
					return true
				})
			}
		} else { // Both are animals
			if meHo != nil && otherHo != nil && meHsObject != nil && otherHsObject != nil &&
				meHo.Life > 80  && otherHo.Life > 80 {
				if meHo.Target != nil {
					meHo.Target = nil
				}
				if otherHo.Target != nil {
					otherHo.Target = nil
				}
				utils.DebugLogf("reproduction")

				mon.ReproductionCounter.Inc()
				//utils.DebugLogf("before %v %v", meHo.Life, otherHo.Life)
				if meHo.addLife(50) {
					ManagerInstance.World.RemoveObject(*meHsObject)
				}
				if otherHo.addLife(50) {
					ManagerInstance.World.RemoveObject(*otherHsObject)
				}
				//utils.DebugLogf("after %v %v", meHo.Life, otherHo.Life)

				speed := ((meHo.MoveSpeed + otherHo.MoveSpeed) / 2) * utils.RandFloats(0.5, 1.5)
				scale := meHo.Scale.Plus(*otherHo.Scale).Times(0.5).Times(utils.RandFloats(0.5, 1.5))

				// Clipping scale ... TODO: add min & max scale somewhere
				clip := func(val float64, min float64, max float64) float64 {
					if val < min {
						return min
					} else if val > max {
						return max
					} else {
						return val
					}
				}
				scale.X = clip(scale.X, 0.1, 5)
				scale.Y = clip(scale.Y, 0.1, 5)
				scale.Z = clip(scale.Z, 0.1, 5)
				speed = clip(speed, 5, 80)
				position := meHo.Position
				position.Y = scale.Y // To stay above ground
				ManagerInstance.AddHerbivorous(position, &scale, speed)
			}
		}
	}
}
