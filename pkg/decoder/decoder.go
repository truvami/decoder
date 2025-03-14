package decoder

import (
	"time"
)

type Decoder interface {
	Decode(payload string, port int16, devEui string) (*DecodedUplink, error)
}

type Feature string

const (
	FeatureGNSS            Feature = "gnss"
	FeatureBuffered        Feature = "buffered"
	FeatureBattery         Feature = "battery"
	FeatureTemperature     Feature = "temperature"
	FeatureWiFi            Feature = "wifi"
	FeatureBle             Feature = "ble"
	FeatureButton          Feature = "button"
	FeatureConfig          Feature = "config"
	FeatureMoving          Feature = "moving"
	FeatureDutyCycle       Feature = "dutyCycle"
	FeatureFirmwareVersion Feature = "firmwareVersion"
	FeatureHardwareVersion Feature = "hardwareVersion"
)

type DecodedUplink struct {
	features []Feature
	Data     interface{} `json:"data"`
	Metadata interface{} `json:"metadata"`
}

func NewDecodedUplink(features []Feature, data interface{}, metadata interface{}) *DecodedUplink {
	return &DecodedUplink{
		features: features,
		Data:     data,
		Metadata: metadata,
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

type UplinkFeatureBase interface {
	// GetTimestamp returns the timestamp of the uplink message.
	// Not all uplink messages have a timestamp, so this method returns a pointer to a time.Time.
	// If the uplink message does not have a timestamp, the method returns nil.
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
	GetTTF() *float64
	// GetPDOP returns the position dilution of precision of the GNSS position.
	GetPDOP() *float64
	// GetSatellites returns the number of satellites used to calculate the GNSS position.
	GetSatellites() *uint8
}

type UplinkFeatureBuffered interface {
	// GetBufferLevel returns the buffer level of the device.
	GetBufferLevel() uint16
}

type UpLinkFeatureBattery interface {
	// GetBatteryVoltage returns the battery voltage of the device.
	GetBatteryVoltage() float64
}

type AccessPoint struct {
	MAC  string `json:"mac"`
	RSSI int8   `json:"rssi"`
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

type UplinkFeatureFirmwareVersion interface {
	// GetFirmwareVersion returns the firmware version of the device.
	GetFirmwareVersion() string
}

type UplinkFeatureHardwareVersion interface {
	// GetHardwareVersion returns the hardware version of the device.
	GetHardwareVersion() string
}
