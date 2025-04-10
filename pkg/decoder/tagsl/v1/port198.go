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
	switch p.Reason {
	case 1:
		return decoder.ResetReasonLrr1110FailCode
	case 2:
		return decoder.ResetReasonPowerReset
	case 3:
		return decoder.ResetReasonPinReset
	case 4:
		return decoder.ResetReasonWatchdog
	case 5:
		return decoder.ResetReasonSystemReset
	case 6:
		return decoder.ResetReasonOtherReset
	default:
		return decoder.ResetReasonUnknown
	}
}
