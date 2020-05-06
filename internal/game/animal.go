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


type Herbivorous struct {
	*erutan.Component_SpaceComponent
	*erutan.Component_HealthComponent
	Target *BasicObject
	*erutan.Component_RenderComponent
	*erutan.Component_BehaviourTypeComponent
	*erutan.Component_SpeedComponent
	*erutan.Component_PhysicsComponent
	*erutan.Component_NetworkBehaviourComponent
}


type herbivorousObject struct {
	*erutan.Component_SpaceComponent
	Target *BasicObject
	*erutan.Component_HealthComponent
	*erutan.Component_SpeedComponent
}

// AddLife set health component life, clip it and return true if object is dead
func (h *herbivorousObject) addLife(value float64) bool {
	h.Life += value
	// Clip 0, 100
	/*if h.Life > 100 {
		h.Life = 100
	} else*/if h.Life < 0 {
		h.Life = 0
	}
	mon.LifeGauge.Set(h.Life)
	//utils.DebugLogf("my life %v", h.Life)
	if h.Life == 0 {
		return true
	}
	return false
}

type HerbivorousSystem struct {
	objects octree.Octree
}


func NewHerbivorousSystem() *HerbivorousSystem {
	return &HerbivorousSystem{objects: *octree.NewOctree(protometry.NewBoxOfSize(0, 0, 0, cfg.Global.Logic.OctreeSize))}
}

func (h HerbivorousSystem) Priority() int {
	return 0
}

func (h *HerbivorousSystem) Add(object octree.Object,
	space *erutan.Component_SpaceComponent,
	target *BasicObject,
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
	if !h.objects.Remove(object) {
		utils.DebugLogf("Failed to remove %d, data: %T", object.ID(), object.Data)
	}
}

func (h *HerbivorousSystem) Update(dt float64) {
	for _, object := range h.objects.GetAllObjects() {
		if ho, ok := object.Data.(*herbivorousObject); ok {

			//volume := ho.Component_SpaceComponent.Scale.X * ho.Component_SpaceComponent.Scale.Y * ho.Component_SpaceComponent.Scale.Z

			//Every animal lose life proportional to deltatime, volume and speed
			//So bigger and faster animals need more food
			//if ho.addLife(-3*dt*volume*(ho.MoveSpeed/cfg.Global.Logic.Herbivorous.LifeLossRate)) {
			//	// Dead
			//	ManagerInstance.World.RemoveObject(*object)
			//}
			val := Clip((-dt*ho.Life)/2000, -0.1, -0.0001)
			//utils.DebugLogf("Life loss: %v", val)
			if ho.addLife(val) {
				// Dead
				ManagerInstance.World.RemoveObject(object)
				continue
			}

			// If I don't have a target, let's find one
			if ho.Target == nil {
				h.findTarget(ho)
				//utils.DebugLogf("my target %v", ho.Target)
			}
			if ho.Target != nil { // Can still be nil (no target on the map)
				distance := ho.Position.Distance(*ho.Target.Position)

				//rnd := protometry.RandomCirclePoint(ho.Target.Position.X, ho.Target.Position.Z, 50)
				//newPos := rnd.Minus(*ho.Position)
				newPos := ho.Target.Position.Minus(*ho.Position)
				newPos.Divide(distance)
				newPos.Scale(dt*ho.MoveSpeed)
				newPos.Add(ho.Position)


				//utils.DebugLogf("I'm here %v, I want to move to %v", object.Bounds.GetCenter(), newPos)
				//ManagerInstance.Watch.NotifyAll(obs.Event{Value: obs.PhysicsUpdateRequest{Object: *object, NewPosition: newPos, Dt: dt}})
				ManagerInstance.Watch.NotifyAll(obs.Event{Value: obs.PhysicsUpdateRequest{
					Object: struct{octree.Object;protometry.Vector3}{Object: object, Vector3: newPos},
					Dt: dt}})
			}
		}
	}
}

func (h *HerbivorousSystem) findTarget(ho *herbivorousObject) {
	//utils.DebugLogf("trying to find a target", ho.Life)
	// Reproduction mode
	if ho.Life > cfg.Global.Logic.Herbivorous.ReproductionThreshold {
		for _, object := range h.objects.GetAllObjects() {
			if otherHo, ok := object.Data.(*herbivorousObject); ok {
				if otherHo.Life > cfg.Global.Logic.Herbivorous.ReproductionThreshold {
					newTarget := &BasicObject{}
					newTarget.Component_SpaceComponent = otherHo.Component_SpaceComponent
					ho.Target = newTarget

					newTargetTwo := &BasicObject{}
					newTargetTwo.Component_SpaceComponent = ho.Component_SpaceComponent
					otherHo.Target = newTargetTwo
					// Found a target (another animal)
				}
			}
		}
	}
	// Eating mode
	for _, e := range ManagerInstance.World.Systems() {
		if e, ok := e.(*EatableSystem); ok {
			// Currently look for eatable on the whole map
			// Is there any eatable on the map?
			var minPosition *eatableObject // TODO: fix this doesn't work pick random target it seems xD
			// Next step behaviour would be to check around using getcolliding instead faster, more realistic
			for _, object := range e.objects.GetAllObjects() {
				if eo, ok := object.Data.(*eatableObject); ok {
					if minPosition == nil || ho.Position.Distance(*eo.Position) < ho.Position.Distance(*minPosition.Position) {
						//utils.DebugLogf("eo.position %v distances %v", eo.Position, ho.Position.Distance(*eo.Position))
						minPosition = eo
					}
				}
			}
			//utils.DebugLogf("found a herb target %v", minPosition.Position)

			// Potentially no eatable around
			if minPosition != nil {
				newTarget := &BasicObject{}
				newTarget.Component_SpaceComponent = minPosition.Component_SpaceComponent
				ho.Target = newTarget
			}
		}
	}
}


func (h HerbivorousSystem) Handle(event obs.Event) {
	switch e := event.Value.(type) {
	// In the occurrence of this event we want to check if the animal collided
	// with a vegetation or another animal and take appropriate actions
	case obs.PhysicsUpdateResponse:
		// No collision here
		if len(e.Objects) == 1 {
			//utils.DebugLogf("need to move %v; %v to %v", e.Objects[0].ID(), e.Objects[0].Bounds.GetCenter(), e.Objects[0].Vector3)

			//me := Find(h.objects, e.Objects[0].Object)
			me := h.objects.Get(e.Objects[0].Object.ID(), e.Objects[0].Object.Bounds)

			if me == nil {
				//utils.DebugLogf("Unable to find %v in system %T", e.Me.ID(), h)
				return
			}
			if asHo, ok := me.Data.(*herbivorousObject); ok {
				*asHo.Position = e.Objects[0].Vector3
			}
			// Need to reinsert in the octree
			if !h.objects.Move(me, e.Objects[0].Vector3.X, e.Objects[0].Vector3.Y, e.Objects[0].Vector3.Z) {
				utils.DebugLogf("Failed to move %v", me)
			} else {
				//utils.DebugLogf("moved %v to %v", me.ID(), me.Bounds.GetCenter())
			}
		} else if len(e.Objects) == 2 { // Means collision, shouldn't be > 2 imho
			var meHsObject, otherHsObject *octree.Object
			// We have to retrieve the object in this system
			// Eventually there can be several object in this object bounds (shouldn't happen though if collision are ON + translation)
			for _, o := range h.objects.GetColliding(e.Objects[0].Bounds) {
				if o.Equal(e.Objects[0].Object) {
					meHsObject = &o
				} else if o.Equal(e.Objects[1].Object) { // Else if (we shouldn't receive self-collision)
					otherHsObject = &o
				}
			}

			// Both objects are not in herbivorous system
			if meHsObject == nil && otherHsObject == nil {
				return
			}

			meCo := e.Objects[0].Data.(*collisionObject)
			otherCo := e.Objects[1].Data.(*collisionObject)
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
				// me is a vegetation
				// other is an animal

				if otherHo != nil && otherHsObject != nil {
					if otherHo.addLife(cfg.Global.Logic.Herbivorous.EatLifeGain) {
						ManagerInstance.World.RemoveObject(*otherHsObject)
					}
					mon.EatCounter.Inc()
					// Reset target for everyone that had this target
					for _, e := range h.objects.GetAllObjects() {
						eho, ok := e.Data.(*herbivorousObject)
						if ok && eho.Target == otherHo.Target {
							eho.Target = nil
						}
					}
				}

			} else if otherCo.Tag == erutan.Component_BehaviourTypeComponent_VEGETATION &&
				meCo.Tag == erutan.Component_BehaviourTypeComponent_ANIMAL {

				// other is a vegetation
				// me is an animal

				if meHo != nil && meHsObject != nil {
					if meHo.addLife(cfg.Global.Logic.Herbivorous.EatLifeGain) {
						ManagerInstance.World.RemoveObject(*meHsObject)
					}
					mon.EatCounter.Inc()
					// Reset target for everyone that had this target
					for _, e := range h.objects.GetAllObjects() {
						eho, ok := e.Data.(*herbivorousObject)
						if ok && eho.Target == meHo.Target {
							eho.Target = nil
						}
					}
				}
			} else { // Both are animals
				if meHo != nil && otherHo != nil && meHsObject != nil && otherHsObject != nil &&
					meHo.Life > cfg.Global.Logic.Herbivorous.ReproductionThreshold &&
					otherHo.Life > cfg.Global.Logic.Herbivorous.ReproductionThreshold {
					if meHo.Target != nil {
						meHo.Target = nil
					}
					if otherHo.Target != nil {
						otherHo.Target = nil
					}
					//utils.DebugLogf("reproduction")

					mon.ReproductionCounter.Inc()
					//utils.DebugLogf("before %v %v", meHo.Life, otherHo.Life)
					if meHo.addLife(-cfg.Global.Logic.Herbivorous.ReproductionLifeLoss) {
						ManagerInstance.World.RemoveObject(*meHsObject)
					}
					if otherHo.addLife(-cfg.Global.Logic.Herbivorous.ReproductionLifeLoss) {
						ManagerInstance.World.RemoveObject(*otherHsObject)
					}
					//utils.DebugLogf("after %v %v", meHo.Life, otherHo.Life)

					speed := ((meHo.MoveSpeed + otherHo.MoveSpeed) / 2) * utils.RandFloats(0.5, 1.5)
					scale := meHo.Scale.Plus(*otherHo.Scale).Times(0.5).Times(utils.RandFloats(0.5, 1.5))

					scale.X = Clip(scale.X, 0.1, 5)
					scale.Y = Clip(scale.Y, 0.1, 5)
					scale.Z = Clip(scale.Z, 0.1, 5)
					mon.VolumeGauge.Set(scale.Sum())
					speed = Clip(speed, 5, 80)
					position := protometry.RandomCirclePoint(meHo.Position.X, 0, meHo.Position.Z, 10)
					position.Y = scale.Y // To stay above ground
					ManagerInstance.AddHerbivorous(&position, &scale, speed)
				}
			}
		}
	}
}

// Clip, clip value between min and max
func Clip(val float64, min float64, max float64) float64 {
	if val < min {
		return min
	} else if val > max {
		return max
	} else {
		return val
	}
}