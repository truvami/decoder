package tagsl

import (
	"encoding/json"

	"github.com/truvami/decoder/pkg/decoder"
)

type Port198Payload struct {
	Reason   uint8   `json:"reason"`
	Line     *string `json:"line"`
	File     *string `json:"file"`
	Function *string `json:"function"`
}

func (p Port198Payload) MarshalJSON() ([]byte, error) {
	type Alias Port198Payload
	return json.Marshal(&struct {
		Reason decoder.ResetReason `json:"reason"`
		*Alias
	}{
		Reason: p.GetResetReason(),
		Alias:  (*Alias)(&p),
	})
}

var _ decoder.UplinkFeatureResetReason = &Port198Payload{}

func (p Port198Payload) GetResetReason() decoder.ResetReason {
	var reasons = map[uint8]decoder.ResetReason{
		1: decoder.ResetReasonLrr1110FailCode,
		2: decoder.ResetReasonPowerReset,
		3: decoder.ResetReasonPinReset,
		4: decoder.ResetReasonWatchdog,
		5: decoder.ResetReasonSystemReset,
		6: decoder.ResetReasonOtherReset,
	}

	if reason, ok := reasons[p.Reason]; ok {
		return reason
	}

	return decoder.ResetReasonUnknown
}
