package main

import (
	"time"

	erutan "github.com/user/erutan_two/protos/realtime"
)

// Behaviour gives the base of what any object should be able to do
type Behaviour interface {
	Start()
	Update()
	OnCollisionEnter(other erutan.NetObject)
	OnDestroy()
}

// AbstractBehaviour is a somewhat similar replica of Unity3d MonoBehaviour
// Golang abstract class, store variables ...
type AbstractBehaviour struct {
	Behaviour
	Object erutan.NetObject
}

// Start is used to initialize your object
func (o *AbstractBehaviour) Start() {
	o.Update()
}

// Update is used to handle this object life loop
func (o *AbstractBehaviour) Update() {
	Update(func(deltaTime int64) bool {
		return true
	})
}

// OnCollisionEnter is used to detect entering collision
func (o *AbstractBehaviour) OnCollisionEnter(other erutan.NetObject) {}

// GetObject returns this object's NetObject
//func GetObject(b Behaviour) *erutan.NetObject { return b.(ObjectBehaviour).Object }

// SetObject set object
//func SetObject(b Behaviour, netObject *erutan.NetObject) { b.(ObjectBehaviour).Object = netObject }

// Destroy despawn the current object
func Destroy(o *AbstractBehaviour, delay time.Duration) {}

// OnDestroy is called before getting destroyed
func (o *AbstractBehaviour) OnDestroy() {}
