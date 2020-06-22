package tslab

import (
	"errors"
	"os"
	"sync"

	"github.com/dfense/tslab/things"
	log "github.com/sirupsen/logrus"
)

// show example of a single instance Supervisor in tslab package
var (
	slock       sync.Mutex //lock to change supervisor variables
	initialized uint64     // use atomic reader to verify initialized
	listener    *Listener
	nextID      uint64 // the next ID to be assigned to a new thing

	errNoThingType            = errors.New("no thing type by that name")
	errAlreadyInitialized     = errors.New("listener already initialized")
	ErrInvalidAutoStartOption = errors.New("invalid autostart option")
)

// ConfigData read in on CLI
type ConfigData struct {
	Autostart string
}

// Initialize process command line parameters and initialize the start of app
func Initialize(c ConfigData) error {

	switch c.Autostart {
	case "true":
		CreateThing(things.TBatteryPack, 1)
		CreateThing(things.TInverter, 1)
		CreateThing(things.TLight, 1)
	case "false":

	default:
		return ErrInvalidAutoStartOption
	}

	return nil
}

// SetListener inject a listener into supervisor
// listener can only be set at startup.
func SetListener(l *Listener) error {
	defer slock.Unlock()
	slock.Lock()
	if initialized == 1 {
		return errAlreadyInitialized
	}
	initialized = 1 // set flag so other func that need a listener error out
	listener = l

	return nil
}

// Stop calls Stop() on Listener which stops and deletes all things
// exit = if true stop listener, close channel and io.writer
//        if false just stop all the things
func Stop(exit bool) {
	// dispatch stop to listener
	listener.Stop(exit)

	if exit {
		log.Debugf("Elvis is leaving the building!")
		os.Exit(0)
	}
}

// StopThingsByType shut all agents down by type
func StopThingsByType(tt things.ThingType) error {
	return listener.StopByType(tt)
}

// StopThingsByCID shut down all things by CID
func StopThingsByCID(cid uint64) error {
	return listener.StopByCID(cid)
}

// CreateThing create new thing.
// type = the thing type to start
// qty = number of thing agents to start
func CreateThing(thingtype things.ThingType, qty int) error {

	for i := 0; i < qty; i++ {
		switch thingtype {

		case things.TBatteryPack:
			battery := things.NewBatteryPack(getNextID())

			// add to listener
			listener.SubscribeToThing(&battery)
			log.Printf("batterypack")

		case things.TInverter:
			inverter := things.NewInverter(getNextID())
			log.Printf("inverter")
			listener.SubscribeToThing(&inverter)

		case things.TLight:
			light := things.NewLight(getNextID())
			log.Printf("light")
			listener.SubscribeToThing(&light)

		default:
			return errNoThingType

		}
	}
	return nil
}

// GetThingsList get a list of all running things Short Description
func GetThingsList() []things.CID {
	return listener.GetThingsShortD()
}

// ConfigureWriter take configuration for writer
// CLI options forwarded here (or config file upgrade)
func ConfigureWriter() {

	// talk to listener and inject writer
}

func getNextID() uint64 {
	slock.Lock()
	nextID++
	slock.Unlock()
	return nextID
}
