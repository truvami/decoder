package nomadxs

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	helpers "github.com/truvami/decoder/pkg/common"
	"github.com/truvami/decoder/pkg/decoder"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		payload     string
		port        uint8
		autoPadding bool
		expected    any
	}{
		{
			payload:     "0102d2b47a0081f3f6115219031412361629002300170046fc19098625e3",
			port:        1,
			autoPadding: false,
			expected: Port1Payload{
				Moving:             true,
				Year:               25,
				Month:              3,
				Day:                20,
				Hour:               18,
				Minute:             54,
				Second:             22,
				Latitude:           47.363194,
				Longitude:          8.516598,
				Altitude:           443.4,
				TimeToFix:          time.Duration(41) * time.Second,
				AmbientLight:       35,
				AccelerometerXAxis: 23,
				AccelerometerYAxis: 70,
				AccelerometerZAxis: -999,
				Temperature:        24.38,
				Pressure:           969.9,
			},
		},
		{
			payload:     "0002c420ff005ed85a12b4180719142607240001ffbaffc2fc6f00d71d2e00d6ffc5ff8405310b3810b1",
			port:        1,
			autoPadding: false,
			expected: Port1Payload{
				Moving:             false,
				Year:               24,
				Month:              7,
				Day:                25,
				Hour:               20,
				Minute:             38,
				Second:             7,
				Latitude:           46.407935,
				Longitude:          6.21577,
				Altitude:           478.8,
				TimeToFix:          time.Duration(36) * time.Second,
				AmbientLight:       1,
				AccelerometerXAxis: -70,
				AccelerometerYAxis: -62,
				AccelerometerZAxis: -913,
				Temperature:        2.15,
				Pressure:           747,
				GyroscopeXAxis:     21.4,
				GyroscopeYAxis:     -5.9,
				GyroscopeZAxis:     -12.4,
				MagnetometerXAxis:  1.329,
				MagnetometerYAxis:  2.872,
				MagnetometerZAxis:  4.273,
			},
		},
		{
			payload:     "0002c420ff005ed85a12b4180719142607240001ffbaffc2fc6f",
			port:        1,
			autoPadding: false,
			expected: Port1Payload{
				Moving:             false,
				Year:               24,
				Month:              7,
				Day:                25,
				Hour:               20,
				Minute:             38,
				Second:             7,
				Latitude:           46.407935,
				Longitude:          6.21577,
				Altitude:           478.8,
				TimeToFix:          time.Duration(36) * time.Second,
				AmbientLight:       1,
				AccelerometerXAxis: -70,
				AccelerometerYAxis: -62,
				AccelerometerZAxis: -913,
			},
		},
		{
			payload:     "2c420ff005ed85a12b4180719142607240001ffbaffc2fc6f",
			port:        1,
			autoPadding: true,
			expected: Port1Payload{
				Moving:             false,
				Year:               24,
				Month:              7,
				Day:                25,
				Hour:               20,
				Minute:             38,
				Second:             7,
				Latitude:           46.407935,
				Longitude:          6.21577,
				Altitude:           478.8,
				TimeToFix:          time.Duration(36) * time.Second,
				AmbientLight:       1,
				AccelerometerXAxis: -70,
				AccelerometerYAxis: -62,
				AccelerometerZAxis: -913,
			},
		},
		{
			payload:     "0000007800000708000151800078012c05dc000100010100000258000002580500000000",
			port:        4,
			autoPadding: false,
			expected: Port4Payload{
				LocalizationIntervalWhileMoving: 120,
				LocalizationIntervalWhileSteady: 1800,
				HeartbeatInterval:               86400,
				GPSTimeoutWhileWaitingForFix:    120,
				AccelerometerWakeupThreshold:    300,
				AccelerometerDelay:              1500,
				FirmwareVersionMajor:            0,
				FirmwareVersionMinor:            1,
				FirmwareVersionPatch:            0,
				BatteryKeepAliveMessageInterval: 600,
				HardwareVersionType:             1,
				HardwareVersionRevision:         1,
				ReJoinInterval:                  600,
				AccuracyEnhancement:             5,
				LightLowerThreshold:             0,
				LightUpperThreshold:             0,
			},
		},
		{
			payload:     "7800000708000151800078012c05dc000100010100000258000002580500000000",
			port:        4,
			autoPadding: true,
			expected: Port4Payload{
				LocalizationIntervalWhileMoving: 120,
				LocalizationIntervalWhileSteady: 1800,
				HeartbeatInterval:               86400,
				GPSTimeoutWhileWaitingForFix:    120,
				AccelerometerWakeupThreshold:    300,
				AccelerometerDelay:              1500,
				FirmwareVersionMajor:            0,
				FirmwareVersionMinor:            1,
				FirmwareVersionPatch:            0,
				BatteryKeepAliveMessageInterval: 600,
				HardwareVersionType:             1,
				HardwareVersionRevision:         1,
				ReJoinInterval:                  600,
				AccuracyEnhancement:             5,
				LightLowerThreshold:             0,
				LightUpperThreshold:             0,
			},
		},
		{
			payload:     "010df6",
			port:        15,
			autoPadding: false,
			expected: Port15Payload{
				LowBattery: true,
				Battery:    3.574,
			},
		},
		{
			payload:     "10df6",
			port:        15,
			autoPadding: true,
			expected: Port15Payload{
				LowBattery: true,
				Battery:    3.574,
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestPort%vWith%v", test.port, test.payload), func(t *testing.T) {
			decoder := NewNomadXSv1Decoder(WithAutoPadding(test.autoPadding))
			got, err := decoder.Decode(test.payload, test.port, "")
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

func TestValidationErrors(t *testing.T) {
	tests := []struct {
		payload  string
		port     uint8
		expected error
	}{
		{
			payload:  "0002c420ff005ed85a12b4180719142607240001ffbaffc2fc6f00d71d2e",
			port:     1,
			expected: nil,
		},
		{
			payload:  "0005f5e100005ed85a12b4180719142607240001ffbaffc2fc6f00d71d2e",
			port:     1,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Latitude", 100),
		},
		{
			payload:  "0002c420ff0bebc20012b4180719142607240001ffbaffc2fc6f00d71d2e",
			port:     1,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Longitude", 200),
		},
		{
			payload:  "0002c420ff005ed85a12b4184919142607240001ffbaffc2fc6f00d71d2e",
			port:     1,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Month", 73),
		},
		{
			payload:  "0002c420ff005ed85a12b4180749142607240001ffbaffc2fc6f00d71d2e",
			port:     1,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Day", 73),
		},
		{
			payload:  "0002c420ff005ed85a12b4180719492607240001ffbaffc2fc6f00d71d2e",
			port:     1,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Hour", 73),
		},
		{
			payload:  "0002c420ff005ed85a12b4180719144907240001ffbaffc2fc6f00d71d2e",
			port:     1,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Minute", 73),
		},
		{
			payload:  "0002c420ff005ed85a12b4180719142649240001ffbaffc2fc6f00d71d2e",
			port:     1,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Second", 73),
		},
		{
			payload:  "0002c420ff005ed85a12b4180719142607240001ffbaffc2fc6f17d41d2e",
			port:     1,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Temperature", 61),
		},
		{
			payload:  "0002c420ff005ed85a12b4180719142607240001ffbaffc2fc6f00d72ee0",
			port:     1,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Pressure", 1200),
		},
		{
			payload:  "010df6",
			port:     15,
			expected: nil,
		},
		{
			payload:  "0101f4",
			port:     15,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Battery", 0.5),
		},
		{
			payload:  "01157c",
			port:     15,
			expected: fmt.Errorf("%s for %s %v", helpers.ErrValidationFailed, "Battery", 5.5),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestPort%vValidationWith%v", test.port, test.payload), func(t *testing.T) {
			decoder := NewNomadXSv1Decoder()
			got, err := decoder.Decode(test.payload, test.port, "")

			if err == nil && test.expected == nil {
				return
			}

			t.Logf("got %v", got)

			if err != nil && test.expected == nil || err == nil || err.Error() != test.expected.Error() {
				t.Errorf("expected: %v\ngot: %v", test.expected, err)
			}
		})
	}
}

func TestInvalidPort(t *testing.T) {
	decoder := NewNomadXSv1Decoder()
	_, err := decoder.Decode("00", 0, "")
	if err == nil || err.Error() != "port 0 not supported" {
		t.Fatal("expected port not supported")
	}
}

func TestPayloadTooShort(t *testing.T) {
	decoder := NewNomadXSv1Decoder()
	_, err := decoder.Decode("deadbeef", 1, "")

	if err == nil || !strings.Contains(err.Error(), "payload too short") {
		t.Fatal("expected error payload too short")
	}
}

func TestPayloadTooLong(t *testing.T) {
	decoder := NewNomadXSv1Decoder()
	_, err := decoder.Decode("deadbeef4242deadbeef4242deadbeef4242deadbeef4242deadbeef4242deadbeef4242deadbeef4242deadbeef4242", 1, "")

	if err == nil || !strings.Contains(err.Error(), "payload too long") {
		t.Fatal("expected error payload too long")
	}
}

func TestFeatures(t *testing.T) {
	tests := []struct {
		payload        string
		port           uint8
		skipValidation bool
	}{
		{
			payload: "0002c420ff005ed85a12b4180719142607240001ffbaffc2fc6f00d71d2e",
			port:    1,
		},
		{
			payload: "0000007800000708000151800078012c05dc000100010100000258000002580500000000",
			port:    4,
		},
		{
			payload: "010df6",
			port:    15,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestFeaturesWithPort%vAndPayload%v", test.port, test.payload), func(t *testing.T) {
			d := NewNomadXSv1Decoder(
				WithSkipValidation(test.skipValidation),
			)
			decodedPayload, _ := d.Decode(test.payload, test.port, "")

			// should be able to decode base feature
			base, ok := decodedPayload.Data.(decoder.UplinkFeatureBase)
			if !ok {
				t.Fatalf("expected UplinkFeatureBase, got %T", decodedPayload)
			}
			// check if it panics
			base.GetTimestamp()

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
			if decodedPayload.Is(decoder.FeatureConfig) {
				config, ok := decodedPayload.Data.(decoder.UplinkFeatureConfig)
				if !ok {
					t.Fatalf("expected UplinkFeatureConfig, got %T", decodedPayload)
				}
				// call functions to check if it panics
				config.GetBle()
				config.GetGnss()
				config.GetWifi()
				config.GetAcceleration()
				config.GetMovingInterval()
				config.GetSteadyInterval()
				config.GetConfigInterval()
				config.GetGnssTimeout()
				config.GetAccelerometerThreshold()
				config.GetAccelerometerDelay()
				config.GetBatteryInterval()
				config.GetRejoinInterval()
				config.GetLowLightThreshold()
				config.GetHighLightThreshold()
				config.GetLowTemperatureThreshold()
				config.GetHighTemperatureThreshold()
				config.GetAccessPointsThreshold()
				config.GetBatchSize()
				config.GetBufferSize()
			}
			if decodedPayload.Is(decoder.FeatureFirmwareVersion) {
				firmwareVersion, ok := decodedPayload.Data.(decoder.UplinkFeatureFirmwareVersion)
				if !ok {
					t.Fatalf("expected UplinkFeatureFirmwareVersion, got %T", decodedPayload)
				}
				if firmwareVersion.GetFirmwareVersion() == "" {
					t.Fatalf("expected non empty firmware version")
				}
			}

			if decodedPayload.Is(decoder.FeatureHardwareVersion) {
				hardwareVersion, ok := decodedPayload.Data.(decoder.UplinkFeatureHardwareVersion)
				if !ok {
					t.Fatalf("expected UplinkFeatureHardwareVersion, got %T", decodedPayload)
				}
				if hardwareVersion.GetHardwareVersion() == "" {
					t.Fatalf("expected non empty hardware version")
				}
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
			payload:  "0002c420ff005ed85a12b4180719142607240001ffbaffc2fc6f00d71d2e",
			port:     1,
			expected: []string{"\"altitude\": 478.8", "\"temperature\": 2.15", "\"timeToFix\": \"36s\""},
		},
		{
			payload:  "0000007800000708000151800078012c05dc000100010100000258000002580500000000",
			port:     4,
			expected: []string{"\"heartbeatInterval\": 86400", "\"reJoinInterval\": 600"},
		},
		{
			payload:  "010df6",
			port:     15,
			expected: []string{"\"lowBattery\": true", "\"battery\": 3.574"},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestMarshalWithPort%vAndPayload%v", test.port, test.payload), func(t *testing.T) {
			decoder := NewNomadXSv1Decoder()

			data, _ := decoder.Decode(test.payload, test.port, "")

			marshaled, err := json.MarshalIndent(map[string]any{
				"data":     data.Data,
				"metadata": data.Metadata,
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
