package decoder

import (
	"testing"
	"time"
)

func TestNewDecodedUplink_IsAndGetFeatures(t *testing.T) {
	features := []Feature{
		FeatureTimestamp,
		FeatureGNSS,
		FeatureBattery,
		FeatureConfigChange,
	}

	data := struct {
		Value string
	}{
		Value: "ok",
	}

	d := NewDecodedUplink(features, data)

	// Check GetFeatures returns an equal slice (order preserved)
	got := d.GetFeatures()
	if len(got) != len(features) {
		t.Fatalf("GetFeatures length = %d, want %d", len(got), len(features))
	}
	for i := range features {
		if got[i] != features[i] {
			t.Fatalf("GetFeatures[%d] = %v, want %v", i, got[i], features[i])
		}
	}

	// Check Is for present features
	for _, f := range features {
		if !d.Is(f) {
			t.Fatalf("Is(%q) = false, want true", f)
		}
	}

	// Check Is for absent features
	if d.Is(FeatureWiFi) {
		t.Fatalf("Is(FeatureWiFi) = true, want false")
	}

	// Just ensure the interfaces compile against this package types;
	// these are compile-time checks — not executed but ensure API coverage.
	var _ UplinkFeatureTimestamp = (*dummyTimestamp)(nil)
	var _ UplinkFeatureGNSS = (*dummyGNSS)(nil)
	var _ UplinkFeatureBuffered = (*dummyBuffered)(nil)
	var _ UplinkFeatureBattery = (*dummyBattery)(nil)
	var _ UplinkFeaturePhotovoltaic = (*dummyPhotovoltaic)(nil)
	var _ UplinkFeatureTemperature = (*dummyTemperature)(nil)
	var _ UplinkFeatureHumidity = (*dummyHumidity)(nil)
	var _ UplinkFeaturePressure = (*dummyPressure)(nil)
	var _ UplinkFeatureWiFi = (*dummyWiFi)(nil)
	var _ UplinkFeatureMoving = (*dummyMoving)(nil)
	var _ UplinkFeatureDutyCycle = (*dummyDutyCycle)(nil)
	var _ UplinkFeatureConfig = (*dummyConfig)(nil)
	var _ UplinkFeatureConfigChange = (*dummyConfigChange)(nil)
	var _ UplinkFeatureFirmwareVersion = (*dummyFirmwareVersion)(nil)
	var _ UplinkFeatureHardwareVersion = (*dummyHardwareVersion)(nil)
	var _ UplinkFeatureButton = (*dummyButton)(nil)
	var _ UplinkFeatureResetReason = (*dummyResetReason)(nil)
	var _ UplinkFeatureRotationState = (*dummyRotationState)(nil)
	var _ UplinkFeatureSequenceNumber = (*dummySequenceNumber)(nil)
}

// The following dummy types satisfy the interfaces to keep API healthy.

type dummyTimestamp struct{}

func (*dummyTimestamp) GetTimestamp() *time.Time { return nil }

type dummyGNSS struct{}

func (*dummyGNSS) GetLatitude() float64   { return 0 }
func (*dummyGNSS) GetLongitude() float64  { return 0 }
func (*dummyGNSS) GetAltitude() float64   { return 0 }
func (*dummyGNSS) GetAccuracy() *float64  { return nil }
func (*dummyGNSS) GetTTF() *time.Duration { return nil }
func (*dummyGNSS) GetPDOP() *float64      { return nil }
func (*dummyGNSS) GetSatellites() *uint8  { return nil }

type dummyBuffered struct{}

func (*dummyBuffered) IsBuffered() bool        { return false }
func (*dummyBuffered) GetBufferLevel() *uint16 { return nil }

type dummyBattery struct{}

func (*dummyBattery) GetBatteryVoltage() float64 { return 0 }
func (*dummyBattery) GetLowBattery() *bool       { return nil }

type dummyPhotovoltaic struct{}

func (*dummyPhotovoltaic) GetPhotovoltaicVoltage() float32 { return 0 }

type dummyTemperature struct{}

func (*dummyTemperature) GetTemperature() float32 { return 0 }

type dummyHumidity struct{}

func (*dummyHumidity) GetHumidity() float32 { return 0 }

type dummyPressure struct{}

func (*dummyPressure) GetPressure() float32 { return 0 }

type dummyWiFi struct{}

func (*dummyWiFi) GetAccessPoints() []AccessPoint { return nil }

type dummyMoving struct{}

func (*dummyMoving) IsMoving() bool { return false }

type dummyDutyCycle struct{}

func (*dummyDutyCycle) IsDutyCycle() bool { return false }

type dummyConfig struct{}

func (*dummyConfig) GetBle() *bool                      { return nil }
func (*dummyConfig) GetGnss() *bool                     { return nil }
func (*dummyConfig) GetWifi() *bool                     { return nil }
func (*dummyConfig) GetAcceleration() *bool             { return nil }
func (*dummyConfig) GetMovingInterval() *uint32         { return nil }
func (*dummyConfig) GetSteadyInterval() *uint32         { return nil }
func (*dummyConfig) GetConfigInterval() *uint32         { return nil }
func (*dummyConfig) GetGnssTimeout() *uint16            { return nil }
func (*dummyConfig) GetAccelerometerThreshold() *uint16 { return nil }
func (*dummyConfig) GetAccelerometerDelay() *uint16     { return nil }
func (*dummyConfig) GetBatteryInterval() *uint32        { return nil }
func (*dummyConfig) GetRejoinInterval() *uint32         { return nil }
func (*dummyConfig) GetLowLightThreshold() *uint16      { return nil }
func (*dummyConfig) GetHighLightThreshold() *uint16     { return nil }
func (*dummyConfig) GetLowTemperatureThreshold() *int8  { return nil }
func (*dummyConfig) GetHighTemperatureThreshold() *int8 { return nil }
func (*dummyConfig) GetAccessPointsThreshold() *uint8   { return nil }
func (*dummyConfig) GetBatchSize() *uint16              { return nil }
func (*dummyConfig) GetBufferSize() *uint16             { return nil }
func (*dummyConfig) GetDataRate() *DataRate             { return nil }

type dummyConfigChange struct{}

func (*dummyConfigChange) GetConfigId() *uint8   { return nil }
func (*dummyConfigChange) GetConfigChange() bool { return false }

type dummyFirmwareVersion struct{}

func (*dummyFirmwareVersion) GetFirmwareHash() *string    { return nil }
func (*dummyFirmwareVersion) GetFirmwareVersion() *string { return nil }

type dummyHardwareVersion struct{}

func (*dummyHardwareVersion) GetHardwareVersion() string { return "" }

type dummyButton struct{}

func (*dummyButton) GetPressed() bool { return false }

type dummyResetReason struct{}

func (*dummyResetReason) GetResetReason() ResetReason { return ResetReasonUnknown }

type dummyRotationState struct{}

func (*dummyRotationState) GetOldRotationState() RotationState { return RotationStateUndefined }
func (*dummyRotationState) GetNewRotationState() RotationState { return RotationStateUndefined }
func (*dummyRotationState) GetRotations() float64              { return 0 }
func (*dummyRotationState) GetDuration() time.Duration         { return 0 }

type dummySequenceNumber struct{}

func (*dummySequenceNumber) GetSequenceNumber() *uint16 { return nil }
