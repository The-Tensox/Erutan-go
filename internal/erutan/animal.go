package erutan

import (
	"github.com/The-Tensox/Erutan-go/internal/cfg"
	"github.com/The-Tensox/Erutan-go/internal/log"
	"github.com/The-Tensox/Erutan-go/internal/mon"
	"github.com/The-Tensox/Erutan-go/internal/obs"
	"github.com/The-Tensox/Erutan-go/internal/utils"
	"github.com/The-Tensox/octree"
	"github.com/The-Tensox/protometry"
	"go.uber.org/zap"
	"reflect"
)


type Herbivorous struct {
	*Component_SpaceComponent
	*Component_HealthComponent
	Target *BasicObject
	*Component_RenderComponent
	*Component_BehaviourTypeComponent
	*Component_SpeedComponent
	*Component_PhysicsComponent
	*Component_NetworkBehaviourComponent
}


type herbivorousObject struct {
	*Component_SpaceComponent
	Target *BasicObject
	*Component_HealthComponent
	*Component_SpeedComponent
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
	//log.Zap.Info("my life %v", h.Life)
	if h.Life == 0 {
		return true
	}
	return false
}

type HerbivorousSystem struct {
	objects octree.Octree
}


func NewHerbivorousSystem() *HerbivorousSystem {
	return &HerbivorousSystem{objects: *octree.NewOctree(protometry.NewBoxOfSize(0, 0, 0, cfg.Get().Logic.OctreeSize))}
}

func (h HerbivorousSystem) Priority() int {
	return 0
}

func (h *HerbivorousSystem) Add(object octree.Object,
	space *Component_SpaceComponent,
	target *BasicObject,
	health *Component_HealthComponent,
	speed *Component_SpeedComponent) {
	ho := &herbivorousObject{space, target,
		health, speed}
	object.Data = ho
	if !h.objects.Insert(object) {
		log.Zap.Info("Failed to insert", zap.Any("object", object))
	}
	mon.SpeedGauge.Set(speed.MoveSpeed)
	mon.VolumeGauge.Set(space.Scale.Sum())
}

// Remove removes the Object from the System. This is what most Remove methods will look like
func (h *HerbivorousSystem) Remove(object octree.Object) {
	if !h.objects.Remove(object) {
		log.Zap.Info("Failed to remove", zap.Any("ID", object.ID()), zap.Any("data", reflect.TypeOf(object.Data)))
	}
}

func (h *HerbivorousSystem) Update(dt float64) {
	for _, object := range h.objects.GetAllObjects() {
		if ho, ok := object.Data.(*herbivorousObject); ok {
			val := Clip((-dt*ho.Life)/2000, -0.1, -0.0001)
			if ho.addLife(val) {
				// Dead
				ManagerInstance.RemoveObject(object)
				continue
			}

			// If I don't have a target, let's find one
			if ho.Target == nil {
				h.findTarget(ho)
			}
			if ho.Target != nil { // Can still be nil (no target on the map)
				distance := ho.Position.Distance(*ho.Target.Position)
				newPos := ho.Target.Position.Minus(*ho.Position)
				newPos.Divide(distance)
				newPos.Scale(dt*ho.MoveSpeed)
				newPos.Add(ho.Position)
				ManagerInstance.NotifyAll(obs.Event{Value: PhysicsUpdateRequest{
					Object: struct{octree.Object;protometry.Vector3}{Object: *object.Clone(), Vector3: newPos},
					Dt: dt}})
			}
		}
	}
}

func (h *HerbivorousSystem) findTarget(ho *herbivorousObject) {
	// Reproduction mode
	if ho.Life > cfg.Get().Logic.Herbivorous.ReproductionThreshold {
		for _, object := range h.objects.GetAllObjects() {
			if otherHo, ok := object.Data.(*herbivorousObject); ok {
				if otherHo.Life > cfg.Get().Logic.Herbivorous.ReproductionThreshold {
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
	for _, e := range ManagerInstance.Systems() {
		if e, ok := e.(*EatableSystem); ok {
			// Currently look for eatable on the whole map
			// Is there any eatable on the map?
			var minPosition *eatableObject // TODO: fix this doesn't work pick random target it seems xD
			// Next step behaviour would be to check around using getcolliding instead faster, more realistic
			for _, object := range e.objects.GetAllObjects() {
				if eo, ok := object.Data.(*eatableObject); ok {
					if minPosition == nil || ho.Position.Distance(*eo.Position) < ho.Position.Distance(*minPosition.Position) {
						minPosition = eo
					}
				}
			}
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
	case PhysicsUpdateResponse:
		// No collision here
		if len(e.Objects) == 1 {
			me := h.objects.Get(e.Objects[0].Object.ID(), e.Objects[0].Object.Bounds)

			if me == nil {
				//log.Zap.Info("Unable to find %v in system %T", e.Me.ID(), h)
				return
			}
			if asHo, ok := me.Data.(*herbivorousObject); ok {
				*asHo.Position = e.Objects[0].Vector3
			}
			// Need to reinsert in the octree
			if !h.objects.Move(me, e.Objects[0].Vector3.X, e.Objects[0].Vector3.Y, e.Objects[0].Vector3.Z) {
				log.Zap.Info("Failed to move", zap.Any("object", me))
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

			if meCo.Tag == Component_BehaviourTypeComponent_VEGETATION &&
				otherCo.Tag == Component_BehaviourTypeComponent_ANIMAL {
				// me is a vegetation
				// other is an animal

				if otherHo != nil && otherHsObject != nil {
					if otherHo.addLife(cfg.Get().Logic.Herbivorous.EatLifeGain) {
						ManagerInstance.RemoveObject(*otherHsObject)
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

			} else if otherCo.Tag == Component_BehaviourTypeComponent_VEGETATION &&
				meCo.Tag == Component_BehaviourTypeComponent_ANIMAL {

				// other is a vegetation
				// me is an animal

				if meHo != nil && meHsObject != nil {
					if meHo.addLife(cfg.Get().Logic.Herbivorous.EatLifeGain) {
						ManagerInstance.RemoveObject(*meHsObject)
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
					meHo.Life > cfg.Get().Logic.Herbivorous.ReproductionThreshold &&
					otherHo.Life > cfg.Get().Logic.Herbivorous.ReproductionThreshold {
					if meHo.Target != nil {
						meHo.Target = nil
					}
					if otherHo.Target != nil {
						otherHo.Target = nil
					}

					mon.ReproductionCounter.Inc()
					if meHo.addLife(-cfg.Get().Logic.Herbivorous.ReproductionLifeLoss) {
						ManagerInstance.RemoveObject(*meHsObject)
					}
					if otherHo.addLife(-cfg.Get().Logic.Herbivorous.ReproductionLifeLoss) {
						ManagerInstance.RemoveObject(*otherHsObject)
					}

					speed := ((meHo.MoveSpeed + otherHo.MoveSpeed) / 2) * utils.RandFloats(0.5, 1.5)
					scale := meHo.Scale.Plus(*otherHo.Scale).Times(0.5).Times(utils.RandFloats(0.5, 1.5))

					scale.X = Clip(scale.X, 0.1, 5)
					scale.Y = Clip(scale.Y, 0.1, 5)
					scale.Z = Clip(scale.Z, 0.1, 5)
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