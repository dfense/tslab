package things

import (
	"math"
	"math/rand"
	"sync"
	"time"
)

const (
	precision = 1000 // round to 4 digits on all floats
)

var (
	// note: crypto rand not required here
	seed = rand.NewSource(time.Now().UnixNano()) // seed for random generators
	rn   = rand.New(seed)                        // init the random object

	// ZeroStruct empty struct to use as trigger
	ZeroStruct = struct{}{}
)

//go:generate stringer -type=ThingType
type ThingType uint8

// possible types to instantiate for thing interface
const (
	TBatteryPack ThingType = iota + 1
	TInverter
	TLight
)

// ThingEvent event that holds things published data
// It enforces certain fields will be implemented by all things
type ThingEvent struct {
	TS        time.Time   `json:"ts"`               //  time event was created
	ThingID   uint64      `json:"event_type_count"` // count of each event type that was created. This is unique to each type
	ThingType string      `json:"thing_type"`       // type of thing that emitted the event
	EventData interface{} `json:"event_data"`       // json serialized struct of each event type
}

// CID short description used to display running CIDs (CodeChallenge ID / things)
type CID struct {
	CidNumber  uint64    // CodeChannenge ID
	CreateTime time.Time // time thing was started
	Type       string    // name of thing type
	TTLEvents  uint64    // total events published
}

// Thing this interface is implemented by all things
type Thing interface {
	Emit(chan<- ThingEvent, *sync.WaitGroup) // subscribe to thing events
	ShortD() CID                             // short discription of thing data
	Stop()                                   // stop sending events, and exit emit()
}

//---------------------------------------------------------
// convenience utility funcs below
//---------------------------------------------------------

// RFloat generic generates a Float64 between min<->max range
// things use it to generate erratic, but bounded data
// round everything to four precision digits for simpler output
func RFloat(min, max float64) float64 {
	r := min + rn.Float64()*(max-min)
	return math.Round(r*precision) / precision
}

// RInt generates an int between min<->max range
func RInt(min, max int) int {
	return rand.Intn(max-min) + min
}
