package tagsl

// +------+------+---------------------------------------------+------------+
// | Byte | Size | Description                                 | Format     |
// +------+------+---------------------------------------------+------------+
// | 0    | 1    | Status[6:2] + Low battery flag[0] (low = 1) | uint8      |
// | 1-2  | 2    | Battery voltage                             | uint16, mV |
// +------+------+---------------------------------------------+------------+

type Port15Payload struct {
	LowBattery     bool    `json:"low_battery"`
	BatteryVoltage float64 `json:"battery_voltage"`
}
