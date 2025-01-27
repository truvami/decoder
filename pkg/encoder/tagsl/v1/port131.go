package tagsl

type Port131Payload struct {
	AccuracyEnhancement uint8 `json:"accuracyEnhancement" validate:"lte=59"`
}
