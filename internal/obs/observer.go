package obs

import (
	erutan "github.com/The-Tensox/erutan/protobuf"
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
		//sync.RWMutex
	}
)

func NewWatch() *Watch {
	return &Watch{}
}



func (w *Watch) Register(observer Observer) {
	//w.Lock()
	//defer w.Unlock()
	//utils.DebugLogf("waiting")

	w.observers = append(w.observers, observer)
}

func (w *Watch) Deregister(observer Observer) {
	//w.Lock()
	//defer w.Unlock()
	for i, o := range w.observers {
		if o == observer {
			w.observers = append(w.observers[:i], w.observers[i+1:]...)
		}
	}
}

func (w *Watch) NotifyAll(event Event) {
	//utils.DebugLogf("lock")
	//w.RLock()
	//defer w.RUnlock()
	//utils.DebugLogf("zaz %T %v", event.Value, event.Value)

	for _, o := range w.observers {
		o.Handle(event)
	}
}


// Events !!
type (
	Event struct {Value interface{}}

	OnClientSettingsUpdate struct {
		ClientToken string
		Settings    erutan.Packet_UpdateParameters
	}
	OnClientConnection struct {
		ClientToken string
		Settings    erutan.Packet_UpdateParameters
	}
	OnPhysicsUpdateRequest struct {
		Object octree.Object
		NewPosition  protometry.Vector3
		Dt     float64
	}
	OnPhysicsUpdateResponse struct {
		Me  *octree.Object
		NewPosition  protometry.Vector3
		// Other is nil if there is no collision
		Other  *octree.Object
		Dt float64
	}
)