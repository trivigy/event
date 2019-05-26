// Package event implements python's threading.Event api using golang primitives.
package event

import (
	"context"
)

// Event is a communication primitive allowing for multiple goroutines to wait
// on an event to be set.
type Event struct {
	next chan chan struct{}
}

// New creates a new event instance.
func New() *Event {
	next := make(chan chan struct{}, 1)
	next <- make(chan struct{})
	return &Event{next: next}
}

// IsSet indicates whether an event is set or not.
func (r *Event) IsSet() bool {
	event := <-r.next
	r.next <- event
	return event == nil
}

// Set changes the internal state of an event to true.
func (r *Event) Set() {
	event := <-r.next
	if event != nil {
		close(event)
		event = nil
	}
	r.next <- event
}

// Clear resets the internal event state back to false.
func (r *Event) Clear() {
	event := <-r.next
	if event == nil {
		event = make(chan struct{})
	}
	r.next <- event
}

// Wait blocks until the event is set to true. If the event is already set,
// returns immediately. Otherwise blocks until another goroutine sets the event.
func (r *Event) Wait(ctx context.Context) bool {
	if ctx == nil {
		panic("context is nil")
	}

	event := <-r.next
	r.next <- event
	if event != nil {
		select {
		case <-ctx.Done():
			return false
		case <-event:
		}
	}
	return true
}
