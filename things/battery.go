package things

import (
	"errors"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
)

// some general sane ranges. not of great value, more than a placeholder
// erratic jumps with Random generator are absurd or maybe comical? but not the point of this challenge
const (
	battRandomDelayMin int     = 1       // least time delay 1ms
	battRandomDelayMax int     = 1000    // most time delay 1s
	minVolts           float64 = 227.00  // 71S * 3.2v
	maxVolts           float64 = 300.00  // 71S * 4.2v
	minTherm           float64 = -40.0   // celcius
	maxTherm           float64 = 175.0   // celcius
	minLiveAmps        float64 = -1000.0 // Amps
	maxLiveAmps        float64 = 1000.0  // Amps
	minCycleAh         float64 = 90.0    // Ah
	maxCycleAh         float64 = 1000.0  // Ah
	minTTLAh           float64 = 0       // kAh
	maxTTLAh           float64 = 30000   // kAh

	errJSONDecoding = "decoding json: %s"
)

var (
	batteryTypeCount uint64 // package counter of things by type

	// internal error objects
	errChannelIsNil = errors.New("channel can not be nil")
)

// BatteryPack very simple for demo.
// TODO should model a more accurate simulator
type BatteryPack struct {
	TTLVoltage  float64       `json:"pack_voltage"` // Total Pack Voltage
	AmpMeter    AmpMeter      `json:"amp_meter"`    // Keep all current flow information in/out of battery
	Therms      []Thermistor  `json:"thermistors"`  // Thermistor array for Battery Pack
	id          uint64        // non serializable id
	createdTime time.Time     // time the object was created
	stopC       chan struct{} // internal stopC interupt
	evtCount    uint64        // number of events generated
	// Cells []Cell
}

// AmpMeter represents battery pack coloumn counter
type AmpMeter struct {
	LiveAmps    float64 `json:"live_amps"`        // realtime or last (n) time buffer avg
	CycleAmpHrs float64 `json:"cycle_amps_hours"` // existing duty cycle
	TTLAmpHours float64 `json:"total_amp_hours"`  // battery odometer
}

// Thermistor state for thermistor temp, and can be enhanced for hi/lo, notification etc
type Thermistor struct {
	Temp float64 `json:"temperature"`
}

// NewBatteryPack create a battery allocating configuration
// ID = cid  code challenge id. Increment ID created by supervisor unique to all things
func NewBatteryPack(ID uint64) BatteryPack {

	// generate random data init
	battery := BatteryPack{id: ID, createdTime: time.Now(), Therms: make([]Thermistor, 2), stopC: make(chan struct{})}
	battery.generateRandomData()
	return battery
}

// Emit implements thing interface to send events over Channel
// c = channel writer for all events
// wg = waitgroup needs to know when we are done with channel
// TODO far too redundant with other Things. Refactor/Reuse
// Note: Not concurrent safe if multiple Supervisors would ever be required!
func (b *BatteryPack) Emit(c chan<- ThingEvent, wg *sync.WaitGroup) {

	defer wg.Done() // tell the listener we are done
	wg.Add(1)
	randomTime := RInt(battRandomDelayMin, battRandomDelayMax)
	delay := time.Duration(randomTime) * time.Millisecond

	if c == nil {
		log.Error(errChannelIsNil)
		return
	}

	thingType := reflect.TypeOf(*b)
EMIT:
	// Begin start lifecycle of thing
	for {
		select {
		case <-time.After(delay):

			// generate random data
			// lock here down if multiple supervisors required
			b.generateRandomData()

			// create new event
			thingEvent := ThingEvent{
				ThingID:   b.id,
				TS:        time.Now(),
				ThingType: thingType.Name(),
				EventData: *b,
			}
			c <- thingEvent
			atomic.AddUint64(&b.evtCount, 1)
		case <-b.stopC:
			break EMIT
		}

		// reset another random time, each time through loop
		randomTime := RInt(battRandomDelayMin, battRandomDelayMax)
		delay = time.Duration(randomTime) * time.Millisecond
	}
	log.Debugf("exiting battery pack: %d", b.id)
}

// ShortD used to give brief data reprentation of this thing. implemnted from things.Thing
func (b BatteryPack) ShortD() CID {
	return CID{CidNumber: b.id, Type: reflect.TypeOf(b).Name(), CreateTime: b.createdTime, TTLEvents: atomic.LoadUint64(&b.evtCount)}
}

// Stop start shutdown sequence
func (b BatteryPack) Stop() {
	b.stopC <- ZeroStruct
}

// generateRandomData just create erratic random data
// TODO model behaivor more realistic
func (b *BatteryPack) generateRandomData() {

	// TODO improve on random data to modeled behaivoral generated data.
	// b.AmpMeter.TtlAmpHours = rand
	b.TTLVoltage = RFloat(minVolts, maxVolts)
	b.AmpMeter.LiveAmps = RFloat(minLiveAmps, maxLiveAmps)
	b.AmpMeter.CycleAmpHrs = RFloat(minCycleAh, maxCycleAh)
	b.AmpMeter.TTLAmpHours = RFloat(minTTLAh, maxTTLAh)
	b.Therms[0].Temp = RFloat(minTherm, maxTherm)
	b.Therms[1].Temp = RFloat(minTherm, maxTherm)

}
