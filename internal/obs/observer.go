package obs

import (
	"github.com/The-Tensox/Erutan-go/internal/ecs"
	"github.com/The-Tensox/Erutan-go/internal/mon"
	erutan "github.com/The-Tensox/Erutan-go/protobuf"
	"github.com/The-Tensox/octree"
	"github.com/The-Tensox/protometry"
	"sort"
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

	// observers implements a sortable list of `Observer`. It is indexed on
	// `Observer.Priority()`.
	observers []Observer


	// Watch implements Observable
	Watch struct {
		observers
		//sync.RWMutex
	}
)

func (o observers) Len() int {
	return len(o)
}

func (o observers) Less(i, j int) bool {
	var prio1, prio2 int

	if prior1, ok := o[i].(ecs.Prioritizer); ok {
		prio1 = prior1.Priority()
	}
	if prior2, ok := o[j].(ecs.Prioritizer); ok {
		prio2 = prior2.Priority()
	}

	return prio1 > prio2
}

func (o observers) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

func NewWatch() *Watch {
	return &Watch{}
}



func (w *Watch) Register(observer Observer) {
	//w.Lock()
	//defer w.Unlock()
	w.observers = append(w.observers, observer)
	sort.Sort(w.observers)
}

func (w *Watch) Deregister(observer Observer) {
	//w.Lock()
	//defer w.Unlock()
	for i, o := range w.observers {
		if o == observer {
			w.observers = append(w.observers[:i], w.observers[i+1:]...)
		}
	}
	sort.Sort(w.observers)
}

func (w *Watch) NotifyAll(event Event) {
	mon.ObserverEventCounter.Inc()
	//w.RLock()
	//defer w.RUnlock()
	for _, o := range w.observers { // TODO: Maybe implement priority stuff (like in ecs package)
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
		// Object is the object to apply physics on associated with a new position, later could add rotation, scale...
		Object struct{octree.Object;protometry.Vector3}
		Dt     float64
	}
	PhysicsUpdateResponse struct {
		// Objects is a slice of object and their new position to be updated, if len > 1 it is the objects collided with
		Objects []struct{octree.Object;protometry.Vector3}
		//Me  *octree.Object
		//NewPosition  protometry.Vector3
		//// Other is nil if there is no collision
		//Other  *octree.Object
		Dt float64
	}
)