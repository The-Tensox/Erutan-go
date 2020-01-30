package main

import erutan "github.com/user/erutan_two/protos/realtime"

type Collider interface {
	GetObject() *erutan.NetObject
	OnCollisionEnter(collisionedObjectID string)
}

type ObjectBehaviour struct {
	Object erutan.NetObject
}

func (o *ObjectBehaviour) OnCollisionEnter(collisionedObjectID string) {
}

func (o *ObjectBehaviour) GetObject() *erutan.NetObject { return &o.Object }
