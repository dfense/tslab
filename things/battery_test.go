package things

import (
	"encoding/json"
	"testing"
)

// TestBatteryJSON verify JSON format, and even more important,
// in case it is broken in future by adding or changing fields.
func TestBatteryJSON(t *testing.T) {

	errJSONnomatch := "json test data does not match expected: %s"
	jsonTestData := `{"pack_voltage":272.683,"amp_meter":{"live_amps":-268.982,"cycle_amps_hours":856.753,"total_amp_hours":3773.437},"thermistors":[{"temperature":63.143},{"temperature":110.421}]}`

	battery := NewBatteryPack(1000)
	battery.TTLVoltage = 272.683
	battery.AmpMeter.LiveAmps = -268.982
	battery.AmpMeter.CycleAmpHrs = 856.753
	battery.AmpMeter.TTLAmpHours = 3773.437
	battery.Therms[0].Temp = 63.143
	battery.Therms[1].Temp = 110.421

	jsonData, err := json.Marshal(battery)
	if err != nil {
		t.Error(err)
	}

	if jsonTestData != string(jsonData) {
		t.Error(errJSONnomatch, string(jsonData))
	}
}
