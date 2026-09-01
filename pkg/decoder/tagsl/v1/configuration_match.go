package tagsl

import (
	"fmt"
	"reflect"

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

type configurationSentPayload struct {
	Ble                    bool
	Gnss                   bool
	Wifi                   bool
	MovingInterval         uint32 `validate:"gte=5,lte=86400"`
	SteadyInterval         uint32 `validate:"gte=120,lte=86400"`
	ConfigInterval         uint32 `validate:"gte=300,lte=604800"`
	GnssTimeout            uint16 `validate:"gte=60,lte=86400"`
	AccelerometerThreshold uint16 `validate:"gte=10,lte=8000"`
	AccelerometerDelay     uint16 `validate:"gte=1000,lte=10000"`
	BatteryInterval        uint32 `validate:"gte=300,lte=604800"`
	BatchSize              uint16 `validate:"lte=50"`
	BufferSize             uint16 `validate:"gte=128,lte=8128"`
}

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

	sentPayload, err := decodeConfigurationSent(sentHex)
	if err != nil {
		return false, err
	}
	observedPayload, err := decodeConfigurationObserved(observedHex)
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

func decodeConfigurationSent(payloadHex string) (configurationSentPayload, error) {
	config := configurationSentPayloadConfig()
	decoded, err := common.Decode(&payloadHex, &config)
	if err != nil {
		return configurationSentPayload{}, err
	}
	return decoded.(configurationSentPayload), nil
}

func decodeConfigurationObserved(payloadHex string) (Port4Payload, error) {
	config := port4PayloadConfig()
	decoded, err := common.Decode(&payloadHex, &config)
	if err != nil {
		return Port4Payload{}, err
	}
	return decoded.(Port4Payload), nil
}

func configurationSentPayloadConfig() common.PayloadConfig {
	return common.PayloadConfig{
		Fields: []common.FieldConfig{
			{Name: "Ble", Start: 0, Length: 1},
			{Name: "Gnss", Start: 1, Length: 1},
			{Name: "Wifi", Start: 2, Length: 1},
			{Name: "MovingInterval", Start: 3, Length: 4},
			{Name: "SteadyInterval", Start: 7, Length: 4},
			{Name: "ConfigInterval", Start: 11, Length: 4},
			{Name: "GnssTimeout", Start: 15, Length: 2},
			{Name: "AccelerometerThreshold", Start: 17, Length: 2},
			{Name: "AccelerometerDelay", Start: 19, Length: 2},
			{Name: "BatteryInterval", Start: 21, Length: 4},
			{Name: "BatchSize", Start: 25, Length: 2},
			{Name: "BufferSize", Start: 27, Length: 2},
		},
		TargetType: reflect.TypeOf(configurationSentPayload{}),
	}
}

func configurationMatches(sent configurationSentPayload, observed Port4Payload) bool {
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
