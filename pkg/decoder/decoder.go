package decoder

import (
	"context"
	"time"
)

type Decoder interface {
	Decode(ctx context.Context, payload string, port uint8) (*DecodedUplink, error)
}

type Feature string

const (
	FeatureTimestamp       Feature = "timestamp"
	FeatureResetReason     Feature = "resetReason"
	FeatureGNSS            Feature = "gnss"
	FeatureBuffered        Feature = "buffered"
	FeatureBattery         Feature = "battery"
	FeaturePhotovoltaic    Feature = "photovoltaic"
	FeatureTemperature     Feature = "temperature"
	FeatureHumidity        Feature = "humidity"
	FeaturePressure        Feature = "pressure"
	FeatureWiFi            Feature = "wifi"
	FeatureBle             Feature = "ble"
	FeatureButton          Feature = "button"
	FeatureConfig          Feature = "config"
	FeatureConfigChange    Feature = "configChange"
	FeatureMoving          Feature = "moving"
	FeatureDutyCycle       Feature = "dutyCycle"
	FeatureFirmwareVersion Feature = "firmwareVersion"
	FeatureHardwareVersion Feature = "hardwareVersion"
	FeatureRotationState   Feature = "rotationState"
	FeatureSequenceNumber  Feature = "sequenceNumber"
)

type DecodedUplink struct {
	features []Feature
	Data     any `json:"data"`
}

func NewDecodedUplink(features []Feature, data any) *DecodedUplink {
	return &DecodedUplink{
		features: features,
		Data:     data,
	}
}

// Is checks if the given feature is present in the DecodedUplink's features.
// It returns true if the feature is found, otherwise it returns false.
//
// Parameters:
//   - feature: The feature to check for in the DecodedUplink.
//
// Returns:
//   - bool: true if the feature is present, false otherwise.
func (d DecodedUplink) Is(feature Feature) bool {
	for _, f := range d.features {
		if f == feature {
			return true
		}
	}
	return false
}

func (d DecodedUplink) GetFeatures() []Feature {
	return d.features
}

type UplinkFeatureTimestamp interface {
	GetTimestamp() *time.Time
}

type UplinkFeatureGNSS interface {
	// GetLatitude returns the latitude of the GNSS position.
	GetLatitude() float64
	// GetLongitude returns the longitude of the GNSS position.
	GetLongitude() float64
	// GetAltitude returns the altitude of the GNSS position.
	GetAltitude() float64
	// GetAccuracy returns the accuracy of the GNSS position.
	GetAccuracy() *float64
	// GetTTF returns the time to fix of the GNSS position.
	GetTTF() *time.Duration
	// GetPDOP returns the position dilution of precision of the GNSS position.
	GetPDOP() *float64
	// GetSatellites returns the number of satellites used to calculate the GNSS position.
	GetSatellites() *uint8
}

type UplinkFeatureBuffered interface {
	IsBuffered() bool
	GetBufferLevel() *uint16
}

type UplinkFeatureBattery interface {
	// GetBatteryVoltage returns the battery voltage of the device.
	GetBatteryVoltage() float64
	// GetLowBattery returns the low battery flag of the device.
	GetLowBattery() *bool
}

type UplinkFeaturePhotovoltaic interface {
	GetPhotovoltaicVoltage() float32
}

type UplinkFeatureTemperature interface {
	GetTemperature() float32
}

type UplinkFeatureHumidity interface {
	GetHumidity() float32
}

type UplinkFeaturePressure interface {
	GetPressure() float32
}

type AccessPoint struct {
	MAC  string `json:"mac"`
	RSSI *int8  `json:"rssi"`
}

type UplinkFeatureWiFi interface {
	// GetAccessPoints returns the list of WiFi access points detected by the device.
	GetAccessPoints() []AccessPoint
}

type UplinkFeatureMoving interface {
	// IsMoving returns true if the device is moving, otherwise it returns false.
	IsMoving() bool
}

type UplinkFeatureDutyCycle interface {
	IsDutyCycle() bool
}

type UplinkFeatureConfig interface {
	GetBle() *bool
	GetGnss() *bool
	GetWifi() *bool
	GetAcceleration() *bool
	GetMovingInterval() *uint32
	GetSteadyInterval() *uint32
	GetConfigInterval() *uint32
	GetGnssTimeout() *uint16
	GetAccelerometerThreshold() *uint16
	GetAccelerometerDelay() *uint16
	GetBatteryInterval() *uint32
	GetRejoinInterval() *uint32
	GetLowLightThreshold() *uint16
	GetHighLightThreshold() *uint16
	GetLowTemperatureThreshold() *int8
	GetHighTemperatureThreshold() *int8
	GetAccessPointsThreshold() *uint8
	GetBatchSize() *uint16
	GetBufferSize() *uint16
	GetDataRate() *DataRate
}

type UplinkFeatureConfigChange interface {
	GetConfigId() *uint8
	GetConfigChange() bool
}

type UplinkFeatureFirmwareVersion interface {
	GetFirmwareHash() *string
	GetFirmwareVersion() *string
}

type UplinkFeatureHardwareVersion interface {
	// GetHardwareVersion returns the hardware version of the device.
	GetHardwareVersion() string
}

type UplinkFeatureButton interface {
	GetPressed() bool
}

type UplinkFeatureResetReason interface {
	GetResetReason() ResetReason
}

type UplinkFeatureRotationState interface {
	GetOldRotationState() RotationState
	GetNewRotationState() RotationState
	GetRotations() float64
	GetDuration() time.Duration
}

type UplinkFeatureSequenceNumber interface {
	GetSequenceNumber() uint
}
