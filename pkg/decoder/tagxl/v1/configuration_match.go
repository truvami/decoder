package tagxl

import (
	"bytes"
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
	errConfigurationTooLarge           = errors.New("tag xl configuration: response exceeds dialect limit")
	errConfigurationTooManyCommands    = errors.New("tag xl configuration: command count exceeds dialect limit")
	errConfigurationUnknownDialect     = errors.New("tag xl configuration: unknown dialect")
)

// ConfigurationComparison is the applicability-aware result of comparing a
// sent settings downlink with a port-151 settings uplink.
type ConfigurationComparison int

const (
	// ConfigurationIncomplete means the uplink omitted at least one required tag.
	ConfigurationIncomplete ConfigurationComparison = iota
	// ConfigurationMismatch means every required tag was present and a setter differed.
	ConfigurationMismatch
	// ConfigurationMatch means every required tag was present and every setter matched.
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

type configurationRequirement struct {
	spec      commandSpec
	wantValue []byte
}

// MatchConfiguration reports whether observed reflects every setter in sent
// using the Tag XL dialect.
func MatchConfiguration(sentHex, observedHex string) (bool, error) {
	result, err := CompareConfiguration(sentHex, observedHex)
	return result == ConfigurationMatch, err
}

// CompareConfiguration classifies observed against the Tag XL dialect.
func CompareConfiguration(sentHex, observedHex string) (ConfigurationComparison, error) {
	return CompareConfigurationFor(ConfigurationDialectTagXL, sentHex, observedHex)
}

// CompareConfigurationFor classifies observed against sent for one dialect.
func CompareConfigurationFor(dialect ConfigurationDialect, sentHex, observedHex string) (ConfigurationComparison, error) {
	spec, err := dialectSpecFor(dialect)
	if err != nil {
		return 0, err
	}
	sent, err := common.HexStringToBytes(sentHex)
	if err != nil {
		return 0, err
	}
	observed, err := common.HexStringToBytes(observedHex)
	if err != nil {
		return 0, err
	}

	requirements, err := parseConfigurationSent(spec, sent)
	if err != nil {
		return 0, err
	}
	seen, err := parseConfigurationObserved(requirements, observed)
	if err != nil {
		return 0, err
	}

	for tag := range requirements {
		if _, ok := seen[tag]; !ok {
			return ConfigurationIncomplete, nil
		}
	}
	for tag, requirement := range requirements {
		if requirement.wantValue == nil {
			continue
		}
		for _, value := range seen[tag] {
			if !compareRequirementValue(requirement, value) {
				return ConfigurationMismatch, nil
			}
		}
	}
	return ConfigurationMatch, nil
}

// ValidateConfiguration reports whether sent is a verifiable payload for dialect.
func ValidateConfiguration(dialect ConfigurationDialect, sentHex string) error {
	spec, err := dialectSpecFor(dialect)
	if err != nil {
		return err
	}
	sent, err := common.HexStringToBytes(sentHex)
	if err != nil {
		return err
	}
	_, err = parseConfigurationSent(spec, sent)
	return err
}

// BuildCurrentConfigurationRequest derives a getter-only readback from sent.
func BuildCurrentConfigurationRequest(dialect ConfigurationDialect, sentHex string) ([]byte, error) {
	spec, err := dialectSpecFor(dialect)
	if err != nil {
		return nil, err
	}
	sent, err := common.HexStringToBytes(sentHex)
	if err != nil {
		return nil, err
	}
	requirements, err := parseConfigurationSent(spec, sent)
	if err != nil {
		return nil, err
	}

	payload := []byte{configurationEnvelopeMarker, 0x00, 0x00}
	for _, requirement := range orderedRequirements(sent, spec, requirements) {
		payload = append(payload, requirement.spec.getterTag, 0x00)
	}
	payload[2] = byte((len(payload) - 3) / 2)
	payload[1] = byte(len(payload) - 2)
	return payload, nil
}

func parseConfigurationSent(spec dialectSpec, payload []byte) (map[byte]configurationRequirement, error) {
	if len(payload) < 3 || payload[0] != configurationEnvelopeMarker {
		return nil, errConfigurationWrongMarker
	}

	requirements := make(map[byte]configurationRequirement)
	commandCount := 0
	offset := 3
	for offset < len(payload) {
		tag, value, next, err := readTLV(payload, offset)
		if err != nil {
			return nil, err
		}
		commandCount++
		if commandCount > spec.limits.maxCommands {
			return nil, errConfigurationTooManyCommands
		}

		if setter, ok := spec.setters[tag]; ok {
			if len(value) != setter.setterLen {
				return nil, errConfigurationMalformedTLV
			}
			if err := validateCommandValue(setter, value); err != nil {
				return nil, err
			}
			if existing, exists := requirements[setter.getterTag]; exists && existing.wantValue != nil {
				return nil, errConfigurationDuplicateTag
			}
			copied := append([]byte(nil), value...)
			requirements[setter.getterTag] = configurationRequirement{spec: setter, wantValue: copied}
			offset = next
			continue
		}

		getter, ok := spec.getters[tag]
		if !ok || len(value) != 0 {
			return nil, errConfigurationUnsupportedCommand
		}
		if _, exists := requirements[getter.getterTag]; !exists {
			requirements[getter.getterTag] = configurationRequirement{spec: getter}
		}
		offset = next
	}

	if len(requirements) == 0 {
		return nil, errConfigurationNoSetter
	}
	if expectedResponseSize(requirements) > spec.limits.maxResponseBytes {
		return nil, errConfigurationTooLarge
	}
	return requirements, nil
}

func parseConfigurationObserved(requirements map[byte]configurationRequirement, payload []byte) (map[byte][][]byte, error) {
	if len(payload) < 3 || payload[0] != configurationEnvelopeMarker {
		return nil, errConfigurationWrongMarker
	}

	seen := make(map[byte][][]byte)
	offset := 3
	for offset < len(payload) {
		tag, value, next, err := readTLV(payload, offset)
		if err != nil {
			return nil, err
		}
		requirement, required := requirements[tag]
		if required {
			if len(value) != requirement.spec.responseLen {
				return nil, errConfigurationMalformedTLV
			}
			if requirement.spec.dataRate && len(value) > 0 && value[0] > configurationMaxDataRate {
				return nil, errConfigurationInvalidDataRate
			}
			seen[tag] = append(seen[tag], append([]byte(nil), value...))
		}
		offset = next
	}
	return seen, nil
}

func validateCommandValue(spec commandSpec, value []byte) error {
	if spec.dataRate && len(value) > 0 && value[0] > configurationMaxDataRate {
		return errConfigurationInvalidDataRate
	}
	if spec.validate != nil {
		return spec.validate(value)
	}
	return nil
}

func compareRequirementValue(requirement configurationRequirement, observed []byte) bool {
	sent := requirement.wantValue
	if requirement.spec.prefixLen > 0 {
		n := requirement.spec.prefixLen
		if len(sent) < n || len(observed) < n {
			return false
		}
		return bytes.Equal(sent[:n], observed[:n])
	}
	if requirement.spec.mask != 0 {
		if len(sent) == 0 || len(observed) == 0 {
			return false
		}
		return (sent[0] & requirement.spec.mask) == (observed[0] & requirement.spec.mask)
	}
	return bytes.Equal(sent, observed)
}

func expectedResponseSize(requirements map[byte]configurationRequirement) int {
	size := 3
	for _, requirement := range requirements {
		size += 2 + requirement.spec.responseLen
	}
	return size
}

func orderedRequirements(payload []byte, spec dialectSpec, requirements map[byte]configurationRequirement) []configurationRequirement {
	ordered := make([]configurationRequirement, 0, len(requirements))
	seen := make(map[byte]struct{}, len(requirements))
	offset := 3
	for offset < len(payload) {
		tag, _, next, err := readTLV(payload, offset)
		if err != nil {
			break
		}
		var getter byte
		if setter, ok := spec.setters[tag]; ok {
			getter = setter.getterTag
		} else if getterSpec, ok := spec.getters[tag]; ok {
			getter = getterSpec.getterTag
		}
		if getter != 0 {
			if _, exists := seen[getter]; !exists {
				if requirement, ok := requirements[getter]; ok {
					ordered = append(ordered, requirement)
					seen[getter] = struct{}{}
				}
			}
		}
		offset = next
	}
	return ordered
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
