package main

import (
	"sync"
)

type (
	// Observable brings a global event dispatcher: Observer design pattern
	Observable interface {
		Add(observer Observer)
		Notify(event Event)
		Remove(event Event)
	}

	// Observer will listen to some events
	Observer interface {
		NotifyCallback(event Event)
	}

	// Watch implements Observable
	Watch struct {
		observer sync.Map
	}

	EventID int

	Event struct {
		eventID EventID
		value   interface{}
	}
)

const (
	FoodMoved EventID = iota
	AnimalReproduced
	AnimalDied
	// ...
)

func (w *Watch) Add(observer Observer) {
	w.observer.Store(observer, struct{}{})
}

func (w *Watch) Remove(observer Observer) {
	w.observer.Delete(observer)
}

func (w *Watch) Notify(event Event) {
	w.observer.Range(func(key, value interface{}) bool {
		if key == nil {
			return false
		}

		key.(Observer).NotifyCallback(event)
		return true
	})
}

/*
Usage:
func (s Soldier) NotifyCallback(event interface{}) {
	if event.(string) == s.zone {
		fmt.Printf("Soldier %d, seen an enemy on zone %s\n", s.id, event)
	}
}
*/
