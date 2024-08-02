package tagxl

type Port152Payload struct {
	NewRotationState  uint8   `json:"newRotationState"`
	OldRotationState  uint8   `json:"oldRotationState"`
	Timestamp         uint32  `json:"timestamp"`
	NumberOfRotations float64 `json:"numberOfRotations"`
	ElapsedSeconds    uint32  `json:"elapsedSeconds"`
}
