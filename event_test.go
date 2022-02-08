package gotrue

import (
	"testing"
	"time"
)

func TestEmitter(t *testing.T) {
	eventChannel := NewEventChannel()

	var ch = make(chan struct{}, 1)

	unsubscribe := eventChannel.Subscribe(SignedInEvent, func() {
		ch <- struct{}{}
	})
	defer unsubscribe()

	go func() {
		time.Sleep(time.Second / 10)
		eventChannel.Publish(SignedInEvent)
	}()

	select {
	case <-time.After(1 * time.Second):
		t.Errorf("event timeout")
		return

	case <-ch:
		return
	}
}
