// Code generated by "stringer -type=ThingType"; DO NOT EDIT.

package things

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TBatteryPack-1]
	_ = x[TInverter-2]
	_ = x[TLight-3]
}

const _ThingType_name = "TBatteryPackTInverterTLight"

var _ThingType_index = [...]uint8{0, 12, 21, 27}

func (i ThingType) String() string {
	i -= 1
	if i >= ThingType(len(_ThingType_index)-1) {
		return "ThingType(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _ThingType_name[_ThingType_index[i]:_ThingType_index[i+1]]
}
