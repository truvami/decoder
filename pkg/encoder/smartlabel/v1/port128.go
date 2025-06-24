package smartlabel

// +------+------+---------------------------------------------+--------------+
// | Byte | Size | Description                                 | Format       |
// +------+------+---------------------------------------------+--------------+
// | 0    | 1    | flags data rate[0:2] acc[3] wifi[4] gnss[5] | byte         |
// | 1    | 2    | steady interval in seconds                  | uint16       |
// | 3    | 2    | moving interval in seconds                  | uint16       |
// | 5    | 1    | config interval in seconds                  | uint8        |
// | 6    | 2    | acceleration threshold                      | uint16, mg   |
// | 8    | 2    | acceleration delay                          | uint16, ms   |
// | 10   | 2    | temperature sensor polling interval         | uint16, s    |
// | 12   | 2    | temperature uplink hold interval            | uint16, s    |
// | 14   | 1    | temperature upper threshold                 | int8, C      |
// | 15   | 1    | temperature lower threshold                 | int8, C      |
// | 16   | 1    | minimal number of access points             | uint8        |
// +------+------+---------------------------------------------+--------------+

type Port128Payload struct {
	DataRate                   uint8  `json:"dataRate" validate:"gte=0,lte=7"`
	Acceleration               bool   `json:"acceleration"`
	Wifi                       bool   `json:"wifi"`
	Gnss                       bool   `json:"gnss"`
	SteadyInterval             uint16 `json:"steadyInterval"`
	MovingInterval             uint16 `json:"movingInterval"`
	HeartbeatInterval          uint8  `json:"heartbeatInterval"`
	AccelerationThreshold      uint16 `json:"accelerationThreshold"`
	AccelerationDelay          uint16 `json:"accelerationDelay"`
	TemperaturePollingInterval uint16 `json:"temperaturePollingInterval"`
	TemperatureUplinkInterval  uint16 `json:"temperatureUplinkInterval"`
	TemperatureUpperThreshold  int8   `json:"temperatureUpperThreshold"`
	TemperatureLowerThreshold  int8   `json:"temperatureLowerThreshold"`
	AccessPointsThreshold      uint8  `json:"accessPointsThreshold" validate:"gte=1,lte=6"`
}
