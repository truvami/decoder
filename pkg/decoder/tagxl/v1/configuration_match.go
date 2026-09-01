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

	sentTLVHex, requestedSetters, err := parseConfigurationSent(sent)
	if err != nil {
		return false, err
	}
	observedTLVHex, err := parseConfigurationObserved(observed)
	if err != nil {
		return false, err
	}

	sentPayload, err := decodePort151Payload(sentTLVHex)
	if err != nil {
		return false, err
	}
	observedPayload, err := decodePort151Payload(observedTLVHex)
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

func parseConfigurationObserved(payload []byte) (tlvHex string, err error) {
	if len(payload) < 3 || payload[0] != configurationEnvelopeMarker {
		return "", errConfigurationWrongMarker
	}

	tlvPayload := []byte{configurationEnvelopeMarker, 0x00, 0x00}
	seen := make(map[byte]struct{})
	offset := 3
	for offset < len(payload) {
		tag, value, next, err := readTLV(payload, offset)
		if err != nil {
			return "", err
		}

		if spec, comparable := comparableTLVSpecs[tag]; comparable {
			if len(value) != spec.valueLen {
				return "", errConfigurationMalformedTLV
			}
			if tag == tlvTagDataRate && value[0] > configurationMaxDataRate {
				return "", errConfigurationInvalidDataRate
			}
			if _, exists := seen[tag]; exists {
				return "", errConfigurationDuplicateTag
			}
			seen[tag] = struct{}{}
			tlvPayload = append(tlvPayload, tag, byte(len(value)))
			tlvPayload = append(tlvPayload, value...)
		}

		offset = next
	}

	return hex.EncodeToString(tlvPayload), nil
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
