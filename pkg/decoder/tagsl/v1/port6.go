package tagsl

// +------+------+----------------------------------------+--------+
// | Byte | Size | Description                            | Format |
// +------+------+----------------------------------------+--------+
// | 0    | 1    | In case of a button-press 0x01 is sent | uint8  |
// +------+------+----------------------------------------+--------+

type Port6Payload struct {
	ButtonPressed bool `json:"button_pressed"`
}
