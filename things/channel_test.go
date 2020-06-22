package things

import (
	"sync"
	"testing"
	"time"
)

var things []Thing // base thing type
// Test to confirm that all things created, reliable shut down with a increment,
// and decrement to waitGroup. Expected failure in this case will be a test timeout
//
// TODO i don't like relying on system timeouts for tests, but this does
// exercise the workflow, and a simple demonstration.
func TestThingShutdown(t *testing.T) {

	maxAllowedTime := int64(2500)                       // milliseconds
	errExceededTimeout := "exceeded maxAllowedTime: %s" // allowed runtime

	startTime := time.Now()

	// create 1 thing of each type to assure everyone behaives with emit rule
	battery := NewBatteryPack(1)
	inverter := NewInverter(2)
	light := NewLight(3)

	// passing reference since things require pointer receiver for emit()
	things := []Thing{&battery, &inverter, &light}
	eventC := make(chan ThingEvent)
	quitC := make(chan struct{})
	wg := &sync.WaitGroup{}

	go func() {
		for {
			select {
			case <-eventC:
				// consume events
			case <-quitC:
				// leave for
			}
		}
	}()

	// subscribe to all things
	for _, t := range things {
		go t.Emit(eventC, wg)
	}
	time.Sleep(time.Millisecond * 2000)

	// stop all things
	for _, t := range things {
		t.Stop() // have things send interrupt and exit, resetting waitgroup semaphore
	}

	// test timeout will occur here if any of the things hang
	wg.Wait()

	now := time.Now()
	elapsed := now.Sub(startTime).Milliseconds()

	if elapsed > maxAllowedTime {
		t.Errorf(errExceededTimeout, elapsed)
	}
}
