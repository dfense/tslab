package tslab

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/dfense/tslab/things"
	log "github.com/sirupsen/logrus"
)

var (
	eventBuffer       = 5 // buffer size of the event channel
	errJSONDecoding   = "decoding json: %s"
	errClosingWriter  = "error closing writer: %s"
	errFlushingBuffer = "error flusing buffer: %s"

	errNoTypeFound = errors.New("no thing(s) with that ThingType found")
	errIDFound     = errors.New("no thing with that CID found")
)

// Listener aggregates all events emitted from things
type Listener struct {
	waitGroup *sync.WaitGroup        // semaphore counter for all things created
	stopC     chan struct{}          // kill channel to stop the server
	writer    io.WriteCloser         // stream to persist all event data
	eventC    chan things.ThingEvent // all thing events feed into this channel:w

	thingList  []things.Thing // base thing type
	thingsLock *sync.Mutex    // lock anytime we alter table or shutdown
}

// NewListener initializes a Listener struct and creates instance
func NewListener() *Listener {

	return &Listener{
		waitGroup:  &sync.WaitGroup{},
		thingsLock: &sync.Mutex{},
		stopC:      make(chan struct{}),
		eventC:     make(chan things.ThingEvent, 5),
	}
}

// SetWriter dependency inject writer for all events
func (l *Listener) SetWriter(w io.WriteCloser) {
	l.writer = w
}

// StartListener receiver call to begin an Aggregator loop of all Events emitting from
// things it subscribes to. If a writer has NOT been set, it will create default
func (l Listener) StartListener() {

	// if writer is not set, create default
	go func() {
		streamBuffer := bufio.NewWriter(l.writer)
		for {
			select {
			case x := <-l.eventC:
				eventJSON, err := json.Marshal(x)
				if err != nil {
					log.Errorf(errJSONDecoding, err)
					continue
				}
				// writeline to event io
				fmt.Fprintln(streamBuffer, string(eventJSON))
			case <-l.stopC:
				// make sure channel is flushed
				close(l.eventC) // close channel
				log.Debug("Turning all the lights out, closing the doors")
				err := streamBuffer.Flush()
				if err != nil {
					log.Errorf(errFlushingBuffer, err)
				}
				err = l.writer.Close()
				if err != nil {
					log.Errorf(errClosingWriter, err)
				}
				// signal to Stop() we are all finished here
				l.waitGroup.Done()
				return
			}
		}
	}()
}

// SubscribeToThing listen for all Events published by a Thing
func (l *Listener) SubscribeToThing(t things.Thing) {

	// lock the list
	defer l.thingsLock.Unlock()
	l.thingsLock.Lock()
	l.thingList = append(l.thingList, t) // add thing to list
	go t.Emit(l.eventC, l.waitGroup)     // start emitting
}

// GetThingsShortD return a short description things.CID of all things
// registered in listener.
func (l Listener) GetThingsShortD() []things.CID {

	cids := make([]things.CID, 0) // create empty list

	defer l.thingsLock.Unlock()
	l.thingsLock.Lock()
	for _, t := range l.thingList {
		cids = append(cids, t.ShortD())
	}

	return cids
}

// Stop kills the listener, only after calling stop on all the things and waiting
// for them to gracefully shutdown
func (l *Listener) Stop(exit bool) {

	// lock list
	defer l.thingsLock.Unlock()
	l.thingsLock.Lock()

	for i, t := range l.thingList {
		log.Debugf("removeSlice[%d] id[%d]\n", i, t.ShortD().CidNumber)
		// pop item from list
		l.thingList = l.thingList[1:]
		t.Stop()
	}
	l.waitGroup.Wait() // waiting for semaphore to hit zero
	log.Debug("WaitGroup returned")

	if exit {
		// wait for eventC to be flushed by our io writer
		for len(l.eventC) > 0 {
			time.Sleep(time.Millisecond * 10)
		}

		// send interrupt to the running listener loop
		l.waitGroup.Add(1)
		l.stopC <- things.ZeroStruct
		l.waitGroup.Wait()
	}
}

// StopByType stop thing by ThingsType
// return errNoTypeFound
func (l *Listener) StopByType(tt things.ThingType) error {

	defer l.thingsLock.Unlock()
	l.thingsLock.Lock()

	// additional lookup
	// convert the string to an enum
	var tmpTType things.ThingType // temp TType from thing

	i := 0                // current position
	s := len(l.thingList) // size of slice
	d := 0                // items deleted

	// remove items as we iterate through the list if they match ThingType
	for i < s {
		t := l.thingList[i]
		switch t.ShortD().Type {
		case "BatteryPack":
			tmpTType = things.TBatteryPack
		case "Inverter":
			tmpTType = things.TInverter
		case "Light":
			tmpTType = things.TLight
		}

		if tmpTType == tt { // thingtype == commandType
			l.thingList[i].Stop()
			l.thingList = append(l.thingList[:i], l.thingList[i+1:]...)
			s--
			d++
		} else {
			i++
		}
	}

	if d == 0 {
		return errNoTypeFound
	}
	return nil
}

// StopByCID stop thing by CID
// returns errNoIDFound
func (l *Listener) StopByCID(cid uint64) error {

	defer l.thingsLock.Unlock()
	l.thingsLock.Lock()
	for i, t := range l.thingList {
		if t.ShortD().CidNumber == cid {
			t.Stop()
			if len(l.thingList) > i {
				l.thingList = append(l.thingList[:i], l.thingList[i+1:]...)
			} else { // last one in list
				l.thingList = l.thingList[:i]
			}
		}
	}
	return nil

}

// createDefaultWrite creates a default file based io writer
func createDefaultWriter() (io.WriteCloser, error) {
	// Create a file for writing
	f, err := os.Create("eventLog.txt")

	return f, err
}
