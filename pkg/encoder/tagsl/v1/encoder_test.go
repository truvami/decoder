package tagsl

import (
	"testing"
)

func TestEncodePort128(t *testing.T) {
    encoder := NewTagSLv1Encoder(
        WithAutoPadding(false),
        WithSkipValidation(true),
    )

    data := Port128Payload{
		BLE:                             1,
		GPS:                             1,
		WIFI:                            1,
		LocalizationIntervalWhileMoving: 3600,
		LocalizationIntervalWhileSteady: 7200,
		HeartbeatInterval:               86400,
		GPSTimeoutWhileWaitingForFix:    120,
		AccelerometerWakeupThreshold:    300,
		AccelerometerDelay:              1500,
		BatteryKeepAliveMessageInterval: 21600,
		BatchSize:                       10,
		BufferSize:                      4096,
	}

	expectedPayload := "01010100000e1000001c20000151800078012c05dc00005460000a1000"

    payload, _, err := encoder.Encode(data, 128, "data")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if payload != expectedPayload {
        t.Errorf("expected payload: %s, got: %s", expectedPayload, payload)
    }
}