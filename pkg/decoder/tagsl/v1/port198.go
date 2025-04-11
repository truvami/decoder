package tagsl

import (
	"time"

	"github.com/truvami/decoder/pkg/decoder"
)

type Port198Payload struct {
	Reason uint8 `json:"reason"`
}

var _ decoder.UplinkFeatureBase = &Port198Payload{}
var _ decoder.UplinkFeatureResetReason = &Port198Payload{}

func (p Port198Payload) GetTimestamp() *time.Time {
	return nil
}

func (p Port198Payload) GetResetReason() decoder.ResetReason {
	var reasonMapping = map[uint8]decoder.ResetReason{
		1: decoder.ResetReasonLrr1110FailCode,
		2: decoder.ResetReasonPowerReset,
		3: decoder.ResetReasonPinReset,
		4: decoder.ResetReasonWatchdog,
		5: decoder.ResetReasonSystemReset,
		6: decoder.ResetReasonOtherReset,
	}

	// Check if the reason is in the mapping
	// If it is, return the corresponding ResetReason
	// If not, return ResetReasonUnknown
	if reason, ok := reasonMapping[p.Reason]; ok {
		return reason
	}

	return decoder.ResetReasonUnknown
}
