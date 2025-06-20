package nomadxl

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	helpers "github.com/truvami/decoder/pkg/common"
	"github.com/truvami/decoder/pkg/decoder"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		payload  string
		port     uint8
		expected any
	}{
		{
			port:    101,
			payload: "00000001fdd5c693000079300001b45d000000000000000000d700000000000000000b3fd724",
			expected: Port101Payload{
				SystemTime:         8553612947,
				UTCDate:            31024,
				UTCTime:            111709,
				Temperature:        21.5,
				Pressure:           0,
				TimeToFix:          time.Duration(36) * time.Second,
				AccelerometerXAxis: 0,
				AccelerometerYAxis: 0,
				AccelerometerZAxis: 0,
				Battery:            2.879,
				BatteryLorawan:     215,
			},
		},
		{
			port:    103,
			payload: "0000793000020152004b6076000c838c00003994",
			expected: Port103Payload{
				UTCDate:   31024,
				UTCTime:   131410,
				Latitude:  49.39894,
				Longitude: 8.20108,
				Altitude:  147.4,
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestPort%vWith%v", test.port, test.payload), func(t *testing.T) {
			decoder := NewNomadXLv1Decoder()
			got, err := decoder.Decode(test.payload, test.port)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			t.Logf("got %v", got)

			if got.Data != test.expected {
				t.Errorf("expected: %v, got: %v", test.expected, got)
			}
		})
	}
}

func TestInvalidPort(t *testing.T) {
	decoder := NewNomadXLv1Decoder()
	_, err := decoder.Decode("00", 0)
	if err == nil || !errors.Is(err, helpers.ErrPortNotSupported) {
		t.Fatal("expected port not supported")
	}
}

func TestPayloadTooShort(t *testing.T) {
	decoder := NewNomadXLv1Decoder()
	_, err := decoder.Decode("deadbeef", 101)

	if err == nil || !errors.Is(err, helpers.ErrPayloadTooShort) {
		t.Fatal("expected error payload too short")
	}
}

func TestPayloadTooLong(t *testing.T) {
	decoder := NewNomadXLv1Decoder()
	_, err := decoder.Decode("deadbeef4242deadbeef4242deadbeef4242deadbeef4242deadbeef4242deadbeef4242deadbeef4242", 101)

	if err == nil || !errors.Is(err, helpers.ErrPayloadTooLong) {
		t.Fatal("expected error payload too long")
	}
}

func TestFeatures(t *testing.T) {
	tests := []struct {
		payload string
		port    uint8
	}{
		{
			payload: "00000001fdd5c693000079300001b45d000000000000000000d71ce60000000000000b3fd724",
			port:    101,
		},
		{
			payload: "0000793000020152004b6076000c838c00003994",
			port:    103,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestFeaturesWithPort%vAndPayload%v", test.port, test.payload), func(t *testing.T) {
			d := NewNomadXLv1Decoder()
			decodedPayload, _ := d.Decode(test.payload, test.port)

			// should be able to decode base feature
			base, ok := decodedPayload.Data.(decoder.UplinkFeatureBase)
			if !ok {
				t.Fatalf("expected UplinkFeatureBase, got %T", decodedPayload)
			}
			// check if it panics
			base.GetTimestamp()

			if len(decodedPayload.GetFeatures()) == 0 {
				t.Error("expected features, got none")
			}

			if decodedPayload.Is(decoder.FeatureGNSS) {
				gnss, ok := decodedPayload.Data.(decoder.UplinkFeatureGNSS)
				if !ok {
					t.Fatalf("expected UplinkFeatureGNSS, got %T", decodedPayload)
				}
				if gnss.GetLatitude() == 0 {
					t.Fatalf("expected non zero latitude")
				}
				if gnss.GetLongitude() == 0 {
					t.Fatalf("expected non zero longitude")
				}
				if gnss.GetAltitude() == 0 {
					t.Fatalf("expected non zero altitude")
				}
				// call function to check if it panics
				gnss.GetAltitude()
				gnss.GetPDOP()
				gnss.GetSatellites()
				gnss.GetTTF()
				gnss.GetAccuracy()
			}
			if decodedPayload.Is(decoder.FeatureTemperature) {
				temperature, ok := decodedPayload.Data.(decoder.UplinkFeatureTemperature)
				if !ok {
					t.Fatalf("expected UplinkFeatureTemperature, got %T", decodedPayload)
				}
				if temperature.GetTemperature() == 0 {
					t.Fatalf("expected non zero temperature")
				}
			}
			if decodedPayload.Is(decoder.FeaturePressure) {
				temperature, ok := decodedPayload.Data.(decoder.UplinkFeaturePressure)
				if !ok {
					t.Fatalf("expected UplinkFeaturePressure, got %T", decodedPayload)
				}
				if temperature.GetPressure() == 0 {
					t.Fatalf("expected non zero pressure")
				}
			}
			if decodedPayload.Is(decoder.FeatureBuffered) {
				buffered, ok := decodedPayload.Data.(decoder.UplinkFeatureBuffered)
				if !ok {
					t.Fatalf("expected UplinkFeatureBuffered, got %T", decodedPayload)
				}
				// call function to check if it panics
				buffered.GetBufferLevel()
			}
			if decodedPayload.Is(decoder.FeatureBattery) {
				batteryVoltage, ok := decodedPayload.Data.(decoder.UplinkFeatureBattery)
				if !ok {
					t.Fatalf("expected UplinkFeatureBattery, got %T", decodedPayload)
				}
				if batteryVoltage.GetBatteryVoltage() == 0 {
					t.Fatalf("expected non zero battery voltage")
				}
				// call function to check if it panics
				batteryVoltage.GetLowBattery()
			}
			if decodedPayload.Is(decoder.FeatureWiFi) {
				wifi, ok := decodedPayload.Data.(decoder.UplinkFeatureWiFi)
				if !ok {
					t.Fatalf("expected UplinkFeatureWiFi, got %T", decodedPayload)
				}
				if wifi.GetAccessPoints() == nil {
					t.Fatalf("expected non nil access points")
				}
			}
			if decodedPayload.Is(decoder.FeatureMoving) {
				moving, ok := decodedPayload.Data.(decoder.UplinkFeatureMoving)
				if !ok {
					t.Fatalf("expected UplinkFeatureMoving, got %T", decodedPayload)
				}
				// call function to check if it panics
				moving.IsMoving()
			}
		})
	}
}

func TestMarshal(t *testing.T) {
	tests := []struct {
		payload  string
		port     uint8
		expected []string
	}{
		{
			payload:  "00000001fdd5c693000079300001b45d000000000000000000d700000000000000000b3fd724",
			port:     101,
			expected: []string{"\"temperature\": 21.5", "\"battery\": 2.879", "\"timeToFix\": \"36s\""},
		},
		{
			payload:  "0000793000020152004b6076000c838c00003994",
			port:     103,
			expected: []string{"\"latitude\": 49.39894", "\"longitude\": 8.20108", "\"altitude\": 147.4"},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestMarshalWithPort%vAndPayload%v", test.port, test.payload), func(t *testing.T) {
			decoder := NewNomadXLv1Decoder()

			data, _ := decoder.Decode(test.payload, test.port)

			marshaled, err := json.MarshalIndent(map[string]any{
				"data": data.Data,
			}, "", "   ")

			if err != nil {
				t.Fatalf("marshalling json failed because %s", err)
			}

			t.Logf("%s\n", marshaled)

			for _, value := range test.expected {
				fmt.Printf("value:%s\n", value)
				if !strings.Contains(string(marshaled), value) {
					t.Fatalf("expected to find %s\n", value)
				}
			}
		})
	}
}
