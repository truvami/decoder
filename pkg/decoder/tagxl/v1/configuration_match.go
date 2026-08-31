package tagxl

import (
	"encoding/hex"
	"errors"

	"github.com/truvami/decoder/pkg/common"
)

const (
	ConfigurationDownlinkPort uint8 = 151
	ConfigurationReportPort   uint8 = 151

	configurationEnvelopeMarker = 0x4c
	configurationMaxDataRate    = 7
)

const (
	configurationSetterTagDeviceFlags           = 0x20
	configurationSetterTagMovingIntervals       = 0x21
	configurationSetterTagAccelerationThreshold = 0x22
	configurationSetterTagHeartbeatInterval     = 0x23
	configurationSetterTagAdvertisementInterval = 0x24
	configurationSetterTagRotationFlags         = 0x25
	configurationSetterTagDataRate              = 0x28

	configurationReportTagDeviceFlags           = 0x40
	configurationReportTagMovingIntervals       = 0x41
	configurationReportTagAccelerationThreshold = 0x42
	configurationReportTagHeartbeatInterval     = 0x43
	configurationReportTagAdvertisementInterval = 0x44
	configurationReportTagRotationFlags         = 0x47
	configurationReportTagDataRate              = 0x4e
)

var (
	errConfigurationWrongMarker        = errors.New("tag xl configuration: wrong envelope marker")
	errConfigurationMalformedTLV       = errors.New("tag xl configuration: malformed tlv")
	errConfigurationNoSetter           = errors.New("tag xl configuration: no setter command")
	errConfigurationUnsupportedCommand = errors.New("tag xl configuration: unsupported command")
	errConfigurationDuplicateTag       = errors.New("tag xl configuration: duplicate comparable tag")
	errConfigurationInvalidDataRate    = errors.New("tag xl configuration: invalid data rate")
)

type configurationSetterSpec struct {
	reportTag byte
	valueLen  int
}

var configurationSetterSpecs = map[byte]configurationSetterSpec{
	configurationSetterTagDeviceFlags:           {reportTag: configurationReportTagDeviceFlags, valueLen: 1},
	configurationSetterTagMovingIntervals:       {reportTag: configurationReportTagMovingIntervals, valueLen: 4},
	configurationSetterTagAccelerationThreshold: {reportTag: configurationReportTagAccelerationThreshold, valueLen: 4},
	configurationSetterTagHeartbeatInterval:     {reportTag: configurationReportTagHeartbeatInterval, valueLen: 1},
	configurationSetterTagAdvertisementInterval: {reportTag: configurationReportTagAdvertisementInterval, valueLen: 1},
	configurationSetterTagRotationFlags:         {reportTag: configurationReportTagRotationFlags, valueLen: 1},
	configurationSetterTagDataRate:              {reportTag: configurationReportTagDataRate, valueLen: 1},
}

var configurationReportSpecs = func() map[byte]configurationSetterSpec {
	specs := make(map[byte]configurationSetterSpec, len(configurationSetterSpecs))
	for _, spec := range configurationSetterSpecs {
		specs[spec.reportTag] = spec
	}
	return specs
}()

// MatchConfiguration reports whether observed reflects every setter in sent.
// false with nil error means a well-formed non-match; non-nil error means malformed or unsupported.
func MatchConfiguration(sentHex, observedHex string) (bool, error) {
	sent, err := common.HexStringToBytes(sentHex)
	if err != nil {
		return false, err
	}
	observed, err := common.HexStringToBytes(observedHex)
	if err != nil {
		return false, err
	}

	reportHex, requestedSetters, err := parseConfigurationSent(sent)
	if err != nil {
		return false, err
	}
	observedReportHex, err := parseConfigurationObserved(observed)
	if err != nil {
		return false, err
	}

	sentPayload, err := decodePort151Configuration(reportHex)
	if err != nil {
		return false, err
	}
	observedPayload, err := decodePort151Configuration(observedReportHex)
	if err != nil {
		return false, err
	}

	for _, setterTag := range requestedSetters {
		if !compareConfigurationSetter(setterTag, sentPayload, observedPayload) {
			return false, nil
		}
	}

	return true, nil
}

func parseConfigurationSent(payload []byte) (reportHex string, requestedSetters []byte, err error) {
	if len(payload) < 3 || payload[0] != configurationEnvelopeMarker {
		return "", nil, errConfigurationWrongMarker
	}

	report := []byte{configurationEnvelopeMarker, 0x00, 0x00}
	setters := make(map[byte]struct{})
	offset := 3
	for offset < len(payload) {
		tag, value, next, err := readConfigurationTLV(payload, offset)
		if err != nil {
			return "", nil, err
		}

		spec, ok := configurationSetterSpecs[tag]
		if !ok {
			return "", nil, errConfigurationUnsupportedCommand
		}
		if len(value) != spec.valueLen {
			return "", nil, errConfigurationMalformedTLV
		}
		if tag == configurationSetterTagDataRate && value[0] > configurationMaxDataRate {
			return "", nil, errConfigurationInvalidDataRate
		}
		if _, exists := setters[tag]; exists {
			return "", nil, errConfigurationDuplicateTag
		}

		setters[tag] = struct{}{}
		requestedSetters = append(requestedSetters, tag)
		report = append(report, spec.reportTag, byte(len(value)))
		report = append(report, value...)
		offset = next
	}

	if len(requestedSetters) == 0 {
		return "", nil, errConfigurationNoSetter
	}

	return hex.EncodeToString(report), requestedSetters, nil
}

func parseConfigurationObserved(payload []byte) (reportHex string, err error) {
	if len(payload) < 3 || payload[0] != configurationEnvelopeMarker {
		return "", errConfigurationWrongMarker
	}

	report := []byte{configurationEnvelopeMarker, 0x00, 0x00}
	reported := make(map[byte]struct{})
	offset := 3
	for offset < len(payload) {
		tag, value, next, err := readConfigurationTLV(payload, offset)
		if err != nil {
			return "", err
		}

		if spec, comparable := configurationReportSpecs[tag]; comparable {
			if len(value) != spec.valueLen {
				return "", errConfigurationMalformedTLV
			}
			if tag == configurationReportTagDataRate && value[0] > configurationMaxDataRate {
				return "", errConfigurationInvalidDataRate
			}
			if _, exists := reported[tag]; exists {
				return "", errConfigurationDuplicateTag
			}
			reported[tag] = struct{}{}
			report = append(report, tag, byte(len(value)))
			report = append(report, value...)
		}

		offset = next
	}

	return hex.EncodeToString(report), nil
}

func decodePort151Configuration(payloadHex string) (Port151Payload, error) {
	config := port151PayloadConfig()
	if err := common.ValidateLength(&payloadHex, &config); err != nil {
		return Port151Payload{}, err
	}

	decoded, err := common.Decode(&payloadHex, &config)
	if err != nil {
		return Port151Payload{}, err
	}

	return decoded.(Port151Payload), nil
}

func compareConfigurationSetter(setterTag byte, sent, observed Port151Payload) bool {
	switch setterTag {
	case configurationSetterTagDeviceFlags:
		return ptrEqual(sent.AccelerometerEnabled, observed.AccelerometerEnabled) &&
			ptrEqual(sent.WifiEnabled, observed.WifiEnabled) &&
			ptrEqual(sent.GnssEnabled, observed.GnssEnabled) &&
			ptrEqual(sent.FirmwareUpgrade, observed.FirmwareUpgrade)
	case configurationSetterTagMovingIntervals:
		return ptrEqual(sent.LocalizationIntervalWhileMoving, observed.LocalizationIntervalWhileMoving) &&
			ptrEqual(sent.LocalizationIntervalWhileSteady, observed.LocalizationIntervalWhileSteady)
	case configurationSetterTagAccelerationThreshold:
		return ptrEqual(sent.AccelerometerWakeupThreshold, observed.AccelerometerWakeupThreshold) &&
			ptrEqual(sent.AccelerometerDelay, observed.AccelerometerDelay)
	case configurationSetterTagHeartbeatInterval:
		return ptrEqual(sent.HeartbeatInterval, observed.HeartbeatInterval)
	case configurationSetterTagAdvertisementInterval:
		return ptrEqual(sent.AdvertisementFirmwareUpgradeInterval, observed.AdvertisementFirmwareUpgradeInterval)
	case configurationSetterTagRotationFlags:
		return ptrEqual(sent.RotationInvert, observed.RotationInvert) &&
			ptrEqual(sent.RotationConfirmed, observed.RotationConfirmed)
	case configurationSetterTagDataRate:
		return ptrEqual(sent.DataRate, observed.DataRate)
	default:
		return false
	}
}

func ptrEqual[T comparable](sent, observed *T) bool {
	if sent == nil || observed == nil {
		return sent == observed
	}
	return *sent == *observed
}

func readConfigurationTLV(payload []byte, offset int) (tag byte, value []byte, next int, err error) {
	if offset+2 > len(payload) {
		return 0, nil, 0, errConfigurationMalformedTLV
	}

	tag = payload[offset]
	length := int(payload[offset+1])
	valueStart := offset + 2
	valueEnd := valueStart + length
	if valueEnd > len(payload) {
		return 0, nil, 0, errConfigurationMalformedTLV
	}

	return tag, payload[valueStart:valueEnd], valueEnd, nil
}
