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
	lRandomDelayMin int = 1000 // least time delay 1ms
	lRandomDelayMax int = 3000 // most time delay 3s
	minCCT          int = 2000 // min ColorSpectrum
	maxCCT          int = 6000 // max ColorSpectrum
	minLL           int = 0    // min LightLevel
	maxLL           int = 100  // max LightLevel
)

// Light defines a luminaire
type Light struct {
	LightLevel    byte  `json:"light_level"`    // live watt reading (or time buffer)
	ColorSpectrum int16 `json:"color_spectrum"` // color spectrum 2000-6000 CCT
	State         bool  `json:"state"`          // state = [on, off] (very simple state)

	// TODO looks like common to Composition candidate
	id          uint64        // non serializable id
	createdTime time.Time     // time the object was created
	stopC       chan struct{} // internal stopC interupt
	evtCount    uint64        // number of events generated
}

// NewLight create a battery allocating configuration
// ID = cid  code challenge id. Increment ID created by supervisor unique to all things
func NewLight(ID uint64) Light {

	l := Light{id: ID, createdTime: time.Now(), stopC: make(chan struct{})}
	l.generateRandomData()
	return l
}

// Emit implements thing interface to send events over Channel
// c = channel writer for all events
// wg = waitgroup needs to know when we are done with channel
// CANDIDATE for Composition
func (l *Light) Emit(c chan<- ThingEvent, wg *sync.WaitGroup) {

	defer wg.Done() // tell the listener we are done
	wg.Add(1)

	if c == nil {
		log.Error(errChannelIsNil)
		return
	}

	thingType := reflect.TypeOf(*l)

	randomTime := RInt(lRandomDelayMin, lRandomDelayMax)
	delay := time.Duration(randomTime) * time.Millisecond

EMIT:
	// Begin start lifecycle of thing
	for {
		select {
		case <-time.After(delay):

			// generate random data
			l.generateRandomData()

			// create new event
			thingEvent := ThingEvent{
				ThingID:   l.id,
				TS:        time.Now(),
				ThingType: thingType.Name(),
				EventData: *l,
			}
			c <- thingEvent
			atomic.AddUint64(&l.evtCount, 1)
		case <-l.stopC:
			break EMIT
		}

		// reset another random time, each time through loop
		randomTime := RInt(lRandomDelayMin, lRandomDelayMax)
		delay = time.Duration(randomTime) * time.Millisecond
	}
	log.Debugf("exiting battery pack: %d", l.id)
}

// ShortD used to give brief data reprentation of this thing. implemnted from things.Thing
func (l Light) ShortD() CID {
	return CID{CidNumber: l.id, Type: reflect.TypeOf(l).Name(), CreateTime: l.createdTime, TTLEvents: atomic.LoadUint64(&l.evtCount)}
}

// Stop break Emit loop
func (l Light) Stop() {
	l.stopC <- ZeroStruct
}

// generateRandomData just create erratic random data
// TODO model behaivor more realistic
func (l *Light) generateRandomData() {

	l.LightLevel = byte(RInt(minLL, maxLL))
	l.ColorSpectrum = int16(RInt(minCCT, maxCCT))
	l.State = bool(!(RInt(0, 2) == 0)) // not == 0 then true

}
