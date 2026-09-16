package tagxl

import "github.com/truvami/decoder/pkg/common"

// ConfigurationDialect selects the fPort-151 command map for a device family.
type ConfigurationDialect int

const (
	ConfigurationDialectUnspecified ConfigurationDialect = iota
	ConfigurationDialectTagXL
	ConfigurationDialectSmartLabelV2
)

func (d ConfigurationDialect) String() string {
	switch d {
	case ConfigurationDialectTagXL:
		return "tagxl"
	case ConfigurationDialectSmartLabelV2:
		return "smartlabel-v2"
	default:
		return "unspecified"
	}
}

type commandSpec struct {
	setterTag   byte
	getterTag   byte
	setterLen   int
	responseLen int
	mask        byte
	prefixLen   int
	dataRate    bool
	validate    func([]byte) error
}

type dialectLimits struct {
	maxCommands      int
	maxResponseBytes int
}

type dialectSpec struct {
	limits  dialectLimits
	setters map[byte]commandSpec
	getters map[byte]commandSpec
	actions map[byte]int
}

func newDialectSpec(limits dialectLimits, commands []commandSpec, actions map[byte]int) dialectSpec {
	spec := dialectSpec{
		limits:  limits,
		setters: make(map[byte]commandSpec, len(commands)),
		getters: make(map[byte]commandSpec, len(commands)),
		actions: actions,
	}
	if spec.actions == nil {
		spec.actions = map[byte]int{}
	}
	for _, command := range commands {
		if command.setterTag != 0 {
			spec.setters[command.setterTag] = command
		}
		if command.getterTag != 0 {
			spec.getters[command.getterTag] = command
		}
	}
	return spec
}

func tagXLHeartbeatValidate(value []byte) error {
	if len(value) == 1 && value[0] > 168 {
		return common.ErrValidationFailed
	}
	return nil
}

func tagXLAdvertisementValidate(value []byte) error {
	if len(value) == 1 && value[0] == 0 {
		return common.ErrValidationFailed
	}
	return nil
}

func tagXLMovingIntervalsValidate(value []byte) error {
	if len(value) != 4 {
		return common.ErrValidationFailed
	}
	packed := common.BytesToUint32(value)
	moving := packed >> 16
	steady := packed & 0xffff
	if moving < 60 || moving > 86400 || steady < 120 || steady > 86400 {
		return common.ErrValidationFailed
	}
	return nil
}

func tagXLAccelerationValidate(value []byte) error {
	if len(value) != 4 {
		return common.ErrValidationFailed
	}
	packed := common.BytesToUint32(value)
	threshold := packed >> 16
	delay := packed & 0xffff
	if threshold < 10 || threshold > 8000 || delay < 1000 || delay > 10000 {
		return common.ErrValidationFailed
	}
	return nil
}

var tagXLDialect = newDialectSpec(dialectLimits{maxCommands: 32, maxResponseBytes: 255}, []commandSpec{
	{setterTag: setterTagDeviceFlags, getterTag: tlvTagDeviceFlags, setterLen: 1, responseLen: 1, mask: 0x0f},
	{setterTag: setterTagMovingIntervals, getterTag: tlvTagMovingIntervals, setterLen: 4, responseLen: 4, validate: tagXLMovingIntervalsValidate},
	{setterTag: setterTagAccelerationThreshold, getterTag: tlvTagAccelerationThreshold, setterLen: 4, responseLen: 4, validate: tagXLAccelerationValidate},
	{setterTag: setterTagHeartbeatInterval, getterTag: tlvTagHeartbeatInterval, setterLen: 1, responseLen: 1, validate: tagXLHeartbeatValidate},
	{setterTag: setterTagAdvertisementInterval, getterTag: tlvTagAdvertisementInterval, setterLen: 1, responseLen: 1, validate: tagXLAdvertisementValidate},
	{setterTag: setterTagRotationFlags, getterTag: tlvTagRotationFlags, setterLen: 1, responseLen: 1, mask: 0x03},
	{setterTag: 0x27, getterTag: 0x4d, setterLen: 2, responseLen: 2},
	{setterTag: setterTagDataRate, getterTag: tlvTagDataRate, setterLen: 1, responseLen: 1, dataRate: true},
	{setterTag: 0x29, getterTag: 0x4f, setterLen: 1, responseLen: 1},
	{setterTag: 0x2a, getterTag: 0x50, setterLen: 1, responseLen: 1, mask: 0x03},
	{setterTag: 0x2b, getterTag: 0x51, setterLen: 1, responseLen: 1, mask: 0x01},
	{getterTag: tlvTagBattery, responseLen: 2},
	{getterTag: tlvTagFirmwareHash, responseLen: 4},
	{getterTag: tlvTagResetCount, responseLen: 2},
	{getterTag: tlvTagResetCause, responseLen: 4},
	{getterTag: tlvTagScanCounts, responseLen: 4},
	{getterTag: 0x4c, responseLen: 4},
	{getterTag: 0x52, responseLen: 6},
}, map[byte]int{
	actionTagAlarm:        2,
	actionTagResetDevice:  0,
	actionTagScanNow:      0,
	actionTagClearStorage: 0,
	actionTagWipeAll:      0,
})

var smartLabelV2Dialect = newDialectSpec(dialectLimits{maxCommands: 30, maxResponseBytes: 51}, []commandSpec{
	{setterTag: 0x20, getterTag: 0x40, setterLen: 1, responseLen: 1, mask: 0x1f},
	{setterTag: 0x21, getterTag: 0x41, setterLen: 4, responseLen: 4},
	{setterTag: 0x22, getterTag: 0x42, setterLen: 4, responseLen: 4, prefixLen: 2},
	{setterTag: 0x23, getterTag: 0x43, setterLen: 1, responseLen: 1},
	{setterTag: 0x24, getterTag: 0x44, setterLen: 1, responseLen: 1},
	{setterTag: 0x28, getterTag: 0x4e, setterLen: 1, responseLen: 1, dataRate: true},
	{setterTag: 0x29, getterTag: 0x4c, setterLen: 2, responseLen: 2},
	{setterTag: 0x2a, getterTag: 0x4d, setterLen: 1, responseLen: 1},
	{setterTag: 0x2b, getterTag: 0x50, setterLen: 1, responseLen: 1},
	{setterTag: 0x2c, getterTag: 0x51, setterLen: 16, responseLen: 16},
	{setterTag: 0x30, getterTag: 0x60, setterLen: 7, responseLen: 7},
	{getterTag: 0x45, responseLen: 2},
	{getterTag: 0x46, responseLen: 4},
	{getterTag: 0x49, responseLen: 2},
	{getterTag: 0x4a, responseLen: 4},
	{getterTag: 0x4b, responseLen: 4},
	{getterTag: 0x61, responseLen: 1},
	{getterTag: 0x63, responseLen: 2},
}, map[byte]int{
	actionTagResetDevice:  0,
	actionTagScanNow:      0,
	actionTagClearStorage: 0,
	actionTagWipeAll:      0,
})

func dialectSpecFor(dialect ConfigurationDialect) (dialectSpec, error) {
	switch dialect {
	case ConfigurationDialectTagXL:
		return tagXLDialect, nil
	case ConfigurationDialectSmartLabelV2:
		return smartLabelV2Dialect, nil
	default:
		return dialectSpec{}, errConfigurationUnknownDialect
	}
}
