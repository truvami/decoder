package smartlabel

// +------+------+---------------------------------------------+--------------+
// | Byte | Size | Description                                 | Format       |
// +------+------+---------------------------------------------+--------------+
// | 0    | 2    | battery voltage for 100 percent             | uint16, mV   |
// | 2    | 2    | battery voltage for 80 percent              | uint16, mV   |
// | 4    | 2    | battery voltage for 60 percent              | uint16, mV   |
// | 6    | 2    | battery voltage for 40 percent              | uint16, mV   |
// | 8    | 2    | battery voltage for 20 percent              | uint16, mV   |
// +------+------+---------------------------------------------+--------------+

type Port150Payload struct {
	Battery100Voltage float32 `json:"battery100Voltage" validate:"gte=3.6,lte=4.0"`
	Battery80Voltage  float32 `json:"battery80Voltage" validate:"gte=3.5,lte=3.7"`
	Battery60Voltage  float32 `json:"battery60Voltage" validate:"gte=3.4,lte=3.6"`
	Battery40Voltage  float32 `json:"battery40Voltage" validate:"gte=3.1,lte=3.4"`
	Battery20Voltage  float32 `json:"battery20Voltage" validate:"gte=2.7,lte=3.0"`
}
