package tagsl

// +------+------+----------------+-----------+
// | Byte | Size | Description    | Format    |
// +------+------+----------------+-----------+
// | 0    | 2    | Scan pointer   | uint16    |
// | 2    | 1    | Total messages | uint8     |
// | 3    | 1    | #Message       | uint8     |
// | 4-9  | 6    | MAC1           | 6 x uint8 |
// | 10   | 1    | RSSI1          | int8      |
// | …    |      |                |           |
// |      | 6    | MACN           | 6 x uint8 |
// |      | 1    | RSSIN          | int8      |
// +------+------+----------------+-----------+

type Port3Payload struct {
	ScanPointer    uint16 `json:"scanPointer"`
	TotalMessages  uint8  `json:"totalMessages"`
	CurrentMessage uint8  `json:"currentMessage"`
	Mac1           string `json:"mac1"`
	Rssi1          int8   `json:"rssi1"`
	Mac2           string `json:"mac2"`
	Rssi2          int8   `json:"rssi2"`
	Mac3           string `json:"mac3"`
	Rssi3          int8   `json:"rssi3"`
	Mac4           string `json:"mac4"`
	Rssi4          int8   `json:"rssi4"`
	Mac5           string `json:"mac5"`
	Rssi5          int8   `json:"rssi5"`
	Mac6           string `json:"mac6"`
	Rssi6          int8   `json:"rssi6"`
}
