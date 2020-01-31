package main

import erutan "github.com/user/erutan_two/protos/realtime"

type Collider interface {
	GetObject() *erutan.NetObject
	OnCollisionEnter(other Collider)
}

type ObjectBehaviour struct {
	Object erutan.NetObject
}

func (o *ObjectBehaviour) OnCollisionEnter(other Collider) {
}

func (o *ObjectBehaviour) GetObject() *erutan.NetObject { return &o.Object }
