package smartlabel

// +-------+------+------------------+------------------+
// | Byte  | Size | Description      | Format           |
// +-------+------+------------------+------------------+
// | 2-3   | 2    | Battery          | uint16, mV       |
// | 4-5   | 2    | Temperature      | uint16, 0.01 C   |
// | 6     | 1    | RH               | uint8, 0.5 %     |
// | 7-8   | 2    | GNSSScansCount   | uint16, -        |
// | 9-10  | 2    | WiFiScansCount   | uint16, -        |
// +-------+------+------------------+------------------+

type Port11HeartbeatPayload struct {
	Battery        float64 `json:"v_bat"`
	Temperature    float64 `json:"temp"`
	RH             float64 `json:"rh"`
	GNSSScansCount uint16  `json:"gnss_scan_count"`
	WiFiScansCount uint16  `json:"wifi_scan_count"`
}
