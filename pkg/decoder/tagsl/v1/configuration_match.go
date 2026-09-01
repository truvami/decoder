package tagsl

import (
	"fmt"

	"github.com/truvami/decoder/pkg/common"
)

const (
	ConfigurationDownlinkPort uint8 = 128
	ConfigurationReportPort   uint8 = 4

	configurationSentByteLength            = 29
	configurationReportCoreByteLength      = 28
	configurationReportWithBatchByteLength = 30
	configurationReportFullByteLength      = 32
)

// MatchConfiguration reports whether observed reflects the comparable fields in sent.
// false with nil error means a well-formed non-match; non-nil error means malformed or out of range.
func MatchConfiguration(sentHex, observedHex string) (bool, error) {
	sent, err := common.HexStringToBytes(sentHex)
	if err != nil {
		return false, err
	}
	observed, err := common.HexStringToBytes(observedHex)
	if err != nil {
		return false, err
	}

	if len(sent) != configurationSentByteLength {
		return false, fmt.Errorf("%w: sent configuration payload length %d", common.ErrInvalidPayloadLength, len(sent))
	}
	if !validConfigurationReportLength(len(observed)) {
		return false, fmt.Errorf("%w: observed configuration payload length %d", common.ErrInvalidPayloadLength, len(observed))
	}

	sentPayload, err := decodePort128Payload(sentHex)
	if err != nil {
		return false, err
	}
	observedPayload, err := decodePort4Payload(observedHex)
	if err != nil {
		return false, err
	}

	return configurationMatches(sentPayload, observedPayload), nil
}

func validConfigurationReportLength(n int) bool {
	switch n {
	case configurationReportCoreByteLength, configurationReportWithBatchByteLength, configurationReportFullByteLength:
		return true
	default:
		return false
	}
}

func decodePort128Payload(payloadHex string) (Port128Payload, error) {
	config := Port128PayloadConfig()
	decoded, err := common.Decode(&payloadHex, &config)
	if err != nil {
		return Port128Payload{}, err
	}
	return decoded.(Port128Payload), nil
}

func decodePort4Payload(payloadHex string) (Port4Payload, error) {
	config := port4PayloadConfig()
	decoded, err := common.Decode(&payloadHex, &config)
	if err != nil {
		return Port4Payload{}, err
	}
	return decoded.(Port4Payload), nil
}

func configurationMatches(sent Port128Payload, observed Port4Payload) bool {
	if sent.MovingInterval != observed.LocalizationIntervalWhileMoving {
		return false
	}
	if sent.SteadyInterval != observed.LocalizationIntervalWhileSteady {
		return false
	}
	if sent.ConfigInterval != observed.HeartbeatInterval {
		return false
	}
	if sent.GnssTimeout != observed.GPSTimeoutWhileWaitingForFix {
		return false
	}
	if sent.AccelerometerThreshold != observed.AccelerometerWakeupThreshold {
		return false
	}
	if sent.AccelerometerDelay != observed.AccelerometerDelay {
		return false
	}
	if sent.BatteryInterval != observed.BatteryKeepAliveMessageInterval {
		return false
	}
	if observed.BatchSize != nil && *observed.BatchSize != sent.BatchSize {
		return false
	}
	if observed.BufferSize != nil && *observed.BufferSize != sent.BufferSize {
		return false
	}
	return true
}
