package things

import (
	"encoding/json"
	"errors"
	"reflect"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	minVoltage     float64 = 227.00 // 71S * 3.2v
	maxVoltage     float64 = 300.00 // 71S * 4.2v
	minTherm       float64 = -40.0  // celcius
	maxTherm       float64 = 175.0  // celcius
	randomDelayMin int     = 1      // least time delay 100ms
	randomDelayMax int     = 1000   // most time delay 1s

	errJSONDecoding = "decoding json: %s"
)

var (
	batteryTypeCount uint64 // package counter of things by type

	// internal error objects
	errChannelIsNil = errors.New("channel can not be nil")
)

// BatteryPack very simple for demo.
// TODO can model a very accurate simulator enhancement
type BatteryPack struct {
	TTLVoltage  float64       // Total Pack Voltage
	AmpMeter    AmpMeter      // Keep all current flow information in/out of battery
	Therms      []Thermistor  // Thermistor array for Battery Pack
	id          uint64        // non serializable id
	createdTime time.Time     // time the object was created
	stopC       chan struct{} // stopC interupt
	// Cells []Cell
}

// AmpMeter represents battery pack coloumn counter
type AmpMeter struct {
	CurrentAmps     float64
	CurrentAmpHours float64
	TTLAmpHours     float64
}

// Thermistor state for thermistor temp, and can be enhanced for hi/lo, notification etc
type Thermistor struct {
	Temp float64
}

// NewBatteryPack create a battery allocating configuration
// ID = cid  code challenge id. Increment ID created by supervisor unique to all things
func NewBatteryPack(ID uint64) BatteryPack {

	// generate random data init
	return BatteryPack{id: ID, createdTime: time.Now(), stopC: make(chan struct{})}
}

// Emit implements thing interface Emit to send events over Channel
// c = channel writer for all events
// wg = waitgroup needs to know when we are done with channel
func (b BatteryPack) Emit(c chan<- ThingEvent, wg *sync.WaitGroup) {

	defer wg.Done() // tell the listener we are done
	wg.Add(1)
	randomTime := RInt(randomDelayMin, randomDelayMax)
	delay := time.Duration(randomTime) * time.Millisecond

	if c == nil {
		log.Error(errChannelIsNil)
		return
	}

	thingType := reflect.TypeOf(b)
EMIT:
	// Begin start lifecycle of thing
	for {
		select {
		case <-time.After(delay):

			// generate random data
			b.generateRandomData()

			// create new event
			eData, err := json.Marshal(b)
			if err != nil {
				log.Errorf(errJSONDecoding, err)
			}

			thingEvent := ThingEvent{
				ThingID:   b.id,
				TS:        time.Now(),
				ThingType: thingType.Name(),
				EventData: string(eData),
			}
			c <- thingEvent
		case <-b.stopC:
			break EMIT
		}

		// reset another random time, each time through loop
		randomTime := RInt(randomDelayMin, randomDelayMax)
		delay = time.Duration(randomTime) * time.Millisecond
	}
	log.Debugf("LEAVING BATTERYPACK %d", b.id)
}

// ShortD used to give brief data reprentation of this thing. implemnted from things.Thing
func (b BatteryPack) ShortD() CID {
	return CID{CidNumber: b.id, Type: reflect.TypeOf(b).Name(), CreateTime: b.createdTime}
}

// Close start shutdown sequence
func (b BatteryPack) Close() {
	b.stopC <- ZeroStruct
}

// generateRandomData just create erratic random data
// TODO model behaivor more realistic
func (b *BatteryPack) generateRandomData() {

	// TODO improve on random data to modeled behaivoral generated data.
	// b.AmpMeter.TtlAmpHours = rand

}
