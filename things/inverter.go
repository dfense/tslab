package things

import (
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
)

// some general sane ranges. not of great value, more than a placeholder
// erratic jumps with Random generator are absurd or maybe comical? but not the point of this challenge
const (
	invRandomDelayMin int     = 50    // least time delay 50ms
	invRandomDelayMax int     = 2000  // most time delay 2s// Inverter minimal sample
	minWatts          float64 = 10    // min watts
	maxWatts          float64 = 10000 // max watts
)

type Inverter struct {
	Watts float64 `json:"watts"` // live watt reading (or time buffer)
	Volts float64 `json:"volts"` // live voltage source
	State bool    `json:"state"` // state = [on, off] (very simple state)

	// TODO looks like common to Composition candidate
	id          uint64        // non serializable id
	createdTime time.Time     // time the object was created
	stopC       chan struct{} // internal stopC interupt
	evtCount    uint64        // number of events generated
}

// NewInverter create a battery allocating configuration
// ID = cid  code challenge id. Increment ID created by supervisor unique to all things
func NewInverter(ID uint64) Inverter {

	i := Inverter{id: ID, createdTime: time.Now(), stopC: make(chan struct{})}
	i.generateRandomData()
	return i
}

// Emit implements thing interface to send events over Channel
// c = channel writer for all events
// wg = waitgroup needs to know when we are done with channel
// CANDIDATE for Composition -
// had to use pointer reference since incrementing evtCount
func (i *Inverter) Emit(c chan<- ThingEvent, wg *sync.WaitGroup) {

	defer wg.Done() // tell the listener we are done
	wg.Add(1)
	randomTime := RInt(invRandomDelayMin, invRandomDelayMax)
	delay := time.Duration(randomTime) * time.Millisecond

	if c == nil {
		log.Error(errChannelIsNil)
		return
	}

	thingType := reflect.TypeOf(*i)
EMIT:
	// Begin start lifecycle of thing
	for {
		select {

		// simulate non-deterministic timing
		case <-time.After(delay):

			// generate random data
			i.generateRandomData()

			// create new event
			thingEvent := ThingEvent{
				ThingID:   (*i).id,
				TS:        time.Now(),
				ThingType: thingType.Name(),
				EventData: *i,
			}
			c <- thingEvent
			atomic.AddUint64(&i.evtCount, 1)

		case <-i.stopC:
			break EMIT
		}

		// reset another random time, each time through loop
		randomTime := RInt(invRandomDelayMin, invRandomDelayMax)
		delay = time.Duration(randomTime) * time.Millisecond
	}
	log.Debugf("exiting battery pack: %d", i.id)
}

// ShortD used to give brief data reprentation of this thing. implemnted from things.Thing
func (i Inverter) ShortD() CID {
	return CID{CidNumber: i.id, Type: reflect.TypeOf(i).Name(), CreateTime: i.createdTime, TTLEvents: atomic.LoadUint64(&i.evtCount)}
}

// Stop start shutdown sequence
func (i Inverter) Stop() {
	i.stopC <- ZeroStruct
}

// generateRandomData just create erratic random data
// TODO model behaivor more realistic
func (i *Inverter) generateRandomData() {

	i.Watts = RFloat(minWatts, maxWatts)
	i.Volts = RFloat(minVolts, maxVolts)
	i.State = bool(!(RInt(0, 2) == 0)) // not == 0 then true

}
