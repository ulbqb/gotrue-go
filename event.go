package gotrue

import (
	"sync"
)

type EventChannel struct {
	sync.RWMutex
	listeners map[AuthChangeEvent][]func()
}

func NewEventChannel() *EventChannel {
	return &EventChannel{
		listeners: make(map[AuthChangeEvent][]func()),
	}
}

func (e *EventChannel) Subscribe(event AuthChangeEvent, fn func()) (unsub func()) {
	e.Lock()
	defer e.Unlock()

	n := len(e.listeners[event])
	e.listeners[event] = append(e.listeners[event], fn)
	return func() {
		e.listeners[event][n] = nil
	}
}

func (e *EventChannel) Publish(event AuthChangeEvent) {
	e.Lock()
	defer e.Unlock()

	if listeners, ok := e.listeners[event]; ok {
		for _, listener := range listeners {
			if listener == nil {
				continue
			}
			go listener()
		}
	}
}
