package obs

import (
	"github.com/The-Tensox/Erutan-go/internal/mon"
	erutan "github.com/The-Tensox/Erutan-go/protobuf"
	"github.com/The-Tensox/octree"
	"github.com/The-Tensox/protometry"
)

type (
	// Observable brings a global event dispatcher: Observer design pattern
	Observable interface {
		Register(observer Observer)
		Deregister(observer Observer)
		NotifyAll(event Event)
	}

	// Observer will listen to some events
	Observer interface {
		Handle(event Event)
	}

	// Watch implements Observable
	Watch struct {
		observers []Observer
	}
)

func NewWatch() *Watch {
	return &Watch{}
}



func (w *Watch) Register(observer Observer) {
	w.observers = append(w.observers, observer)
}

func (w *Watch) Deregister(observer Observer) {
	for i, o := range w.observers {
		if o == observer {
			w.observers = append(w.observers[:i], w.observers[i+1:]...)
		}
	}
}

func (w *Watch) NotifyAll(event Event) {
	mon.ObserverEventCounter.Inc()
	for _, o := range w.observers {
		o.Handle(event)
	}
}


type (
	Event struct {Value interface{}}

	// ClientSettingsUpdate notify of a client updating its settings
	ClientSettingsUpdate struct {
		ClientToken string
		Settings    erutan.Packet_UpdateParameters
	}
	ClientConnection struct {
		ClientToken string
		Settings    erutan.Packet_UpdateParameters
	}
	ClientDisconnection struct {
		ClientToken string
	}
	PhysicsUpdateRequest struct {
		Object octree.Object
		NewPosition  protometry.Vector3
		Dt     float64
	}
	PhysicsUpdateResponse struct {
		Me  *octree.Object
		NewPosition  protometry.Vector3
		// Other is nil if there is no collision
		Other  *octree.Object
		Dt float64
	}
)