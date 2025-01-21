package tagsl

import (
    "testing"
    "time"

    "github.com/truvami/decoder/pkg/decoder/helpers"
)

func TestEncodePort128(t *testing.T) {
    encoder := NewTagSLv1Encoder(
        WithAutoPadding(true),
        WithSkipValidation(false),
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
        DeviceState:                     2,
        FirmwareVersionMajor:            2,
        FirmwareVersionMinor:            1,
        FirmwareVersionPatch:            0,
        HardwareVersionType:             1,
        HardwareVersionRevision:         2,
        BatteryKeepAliveMessageInterval: 21600,
        BatchSize:                       10,
        BufferSize:                      4096,
    }

    expectedPayload := "01010100000e1000001c20000151800078012c05dc02020100010200005460000a1000"

    payload, err := encoder.Encode(data, 128)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if payload != expectedPayload {
        t.Errorf("expected payload: %s, got: %s", expectedPayload, payload)
    }
}