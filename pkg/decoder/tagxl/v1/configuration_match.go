package tagxl

import (
	"encoding/hex"
	"errors"

	"github.com/truvami/decoder/pkg/common"
)

var (
	errConfigurationWrongMarker        = errors.New("tag xl configuration: wrong envelope marker")
	errConfigurationMalformedTLV       = errors.New("tag xl configuration: malformed tlv")
	errConfigurationNoSetter           = errors.New("tag xl configuration: no setter command")
	errConfigurationUnsupportedCommand = errors.New("tag xl configuration: unsupported command")
	errConfigurationDuplicateTag       = errors.New("tag xl configuration: duplicate comparable tag")
	errConfigurationInvalidDataRate    = errors.New("tag xl configuration: invalid data rate")
)

// ConfigurationComparison is the applicability-aware result of comparing a
// sent Tag XL settings downlink with a port-151 settings uplink.
type ConfigurationComparison int

const (
	// ConfigurationIncomplete means the uplink omitted at least one requested setter.
	ConfigurationIncomplete ConfigurationComparison = iota
	// ConfigurationMismatch means every requested setter was present and at least one differed.
	ConfigurationMismatch
	// ConfigurationMatch means every requested setter was present and equal.
	ConfigurationMatch
)

func (c ConfigurationComparison) String() string {
	switch c {
	case ConfigurationIncomplete:
		return "incomplete"
	case ConfigurationMismatch:
		return "mismatch"
	case ConfigurationMatch:
		return "match"
	default:
		return "unknown"
	}
}

// MatchConfiguration reports whether observed reflects every setter in sent.
// false with nil error means a well-formed non-match or incomplete report;
// non-nil error means malformed or unsupported.
func MatchConfiguration(sentHex, observedHex string) (bool, error) {
	result, err := CompareConfiguration(sentHex, observedHex)
	return result == ConfigurationMatch, err
}

// CompareConfiguration classifies observed against every setter in sent.
// Incomplete reports are distinguished from complete contradictions so callers
// can ignore partial telemetry without treating it as a mismatch.
func CompareConfiguration(sentHex, observedHex string) (ConfigurationComparison, error) {
	sent, err := common.HexStringToBytes(sentHex)
	if err != nil {
		return 0, err
	}
	observed, err := common.HexStringToBytes(observedHex)
	if err != nil {
		return 0, err
	}

	sentTLVHex, requestedSetters, err := parseConfigurationSent(sent)
	if err != nil {
		return 0, err
	}
	observedTLVHex, seen, err := parseConfigurationObserved(observed)
	if err != nil {
		return 0, err
	}

	for _, setterTag := range requestedSetters {
		spec, ok := setterSpecs[setterTag]
		if !ok {
			return 0, errConfigurationUnsupportedCommand
		}
		if _, ok := seen[spec.tlvTag]; !ok {
			return ConfigurationIncomplete, nil
		}
	}

	sentPayload, err := decodePort151Payload(sentTLVHex)
	if err != nil {
		return 0, err
	}
	observedPayload, err := decodePort151Payload(observedTLVHex)
	if err != nil {
		return 0, err
	}

	for _, setterTag := range requestedSetters {
		if !compareConfigurationSetter(setterTag, sentPayload, observedPayload) {
			return ConfigurationMismatch, nil
		}
	}

	return ConfigurationMatch, nil
}

func parseConfigurationSent(payload []byte) (tlvHex string, requestedSetters []byte, err error) {
	if len(payload) < 3 || payload[0] != configurationEnvelopeMarker {
		return "", nil, errConfigurationWrongMarker
	}

	tlvPayload := []byte{configurationEnvelopeMarker, 0x00, 0x00}
	setters := make(map[byte]struct{})
	offset := 3
	for offset < len(payload) {
		tag, value, next, err := readTLV(payload, offset)
		if err != nil {
			return "", nil, err
		}

		spec, ok := setterSpecs[tag]
		if !ok {
			return "", nil, errConfigurationUnsupportedCommand
		}
		if len(value) != spec.valueLen {
			return "", nil, errConfigurationMalformedTLV
		}
		if tag == setterTagDataRate && value[0] > configurationMaxDataRate {
			return "", nil, errConfigurationInvalidDataRate
		}
		if _, exists := setters[tag]; exists {
			return "", nil, errConfigurationDuplicateTag
		}

		setters[tag] = struct{}{}
		requestedSetters = append(requestedSetters, tag)
		tlvPayload = append(tlvPayload, spec.tlvTag, byte(len(value)))
		tlvPayload = append(tlvPayload, value...)
		offset = next
	}

	if len(requestedSetters) == 0 {
		return "", nil, errConfigurationNoSetter
	}

	return hex.EncodeToString(tlvPayload), requestedSetters, nil
}

func parseConfigurationObserved(payload []byte) (tlvHex string, seen map[byte]struct{}, err error) {
	if len(payload) < 3 || payload[0] != configurationEnvelopeMarker {
		return "", nil, errConfigurationWrongMarker
	}

	tlvPayload := []byte{configurationEnvelopeMarker, 0x00, 0x00}
	seen = make(map[byte]struct{})
	offset := 3
	for offset < len(payload) {
		tag, value, next, err := readTLV(payload, offset)
		if err != nil {
			return "", nil, err
		}

		if spec, comparable := comparableTLVSpecs[tag]; comparable {
			if len(value) != spec.valueLen {
				return "", nil, errConfigurationMalformedTLV
			}
			if tag == tlvTagDataRate && value[0] > configurationMaxDataRate {
				return "", nil, errConfigurationInvalidDataRate
			}
			if _, exists := seen[tag]; exists {
				return "", nil, errConfigurationDuplicateTag
			}
			seen[tag] = struct{}{}
			tlvPayload = append(tlvPayload, tag, byte(len(value)))
			tlvPayload = append(tlvPayload, value...)
		}

		offset = next
	}

	return hex.EncodeToString(tlvPayload), seen, nil
}

func decodePort151Payload(payloadHex string) (Port151Payload, error) {
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
	case setterTagDeviceFlags:
		return ptrEqual(sent.AccelerometerEnabled, observed.AccelerometerEnabled) &&
			ptrEqual(sent.WifiEnabled, observed.WifiEnabled) &&
			ptrEqual(sent.GnssEnabled, observed.GnssEnabled) &&
			ptrEqual(sent.FirmwareUpgrade, observed.FirmwareUpgrade)
	case setterTagMovingIntervals:
		return ptrEqual(sent.LocalizationIntervalWhileMoving, observed.LocalizationIntervalWhileMoving) &&
			ptrEqual(sent.LocalizationIntervalWhileSteady, observed.LocalizationIntervalWhileSteady)
	case setterTagAccelerationThreshold:
		return ptrEqual(sent.AccelerometerWakeupThreshold, observed.AccelerometerWakeupThreshold) &&
			ptrEqual(sent.AccelerometerDelay, observed.AccelerometerDelay)
	case setterTagHeartbeatInterval:
		return ptrEqual(sent.HeartbeatInterval, observed.HeartbeatInterval)
	case setterTagAdvertisementInterval:
		return ptrEqual(sent.AdvertisementFirmwareUpgradeInterval, observed.AdvertisementFirmwareUpgradeInterval)
	case setterTagRotationFlags:
		return ptrEqual(sent.RotationInvert, observed.RotationInvert) &&
			ptrEqual(sent.RotationConfirmed, observed.RotationConfirmed)
	case setterTagDataRate:
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

func readTLV(payload []byte, offset int) (tag byte, value []byte, next int, err error) {
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
