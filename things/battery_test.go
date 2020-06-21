package things

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestBatteryJSON(t *testing.T) {

	battery := NewBatteryPack(1000)
	battery.generateRandomData()

	jsonData, err := json.Marshal(battery)
	if err != nil {
		t.Error(err)
	}

	// spot check some string compares (simple)
	fmt.Println(string(jsonData))
}
