package nomadxs

import (
	"fmt"
	"testing"

	helpers "github.com/truvami/decoder/pkg/common"
	"github.com/truvami/decoder/pkg/decoder"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		payload     string
		port        int16
		autoPadding bool
		expected    interface{}
	}{
		{
			payload:     "0002c420ff005ed85a12b4180719142607240001ffbaffc2fc6f00d71d2e",
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
				TimeToFix:          36,
				AmbientLight:       1,
				AccelerometerXAxis: -70,
				AccelerometerYAxis: -62,
				AccelerometerZAxis: -913,
				Temperature:        21.5,
				Pressure:           747,
			},
		},
		{
			payload:     "0002c420ff005ed85a12b4180719142607240001ffbaffc2fc6f00d71d2e000000000000",
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
				TimeToFix:          36,
				AmbientLight:       1,
				AccelerometerXAxis: -70,
				AccelerometerYAxis: -62,
				AccelerometerZAxis: -913,
				Temperature:        21.5,
				Pressure:           747,
				GyroscopeXAxis:     0.0,
				GyroscopeYAxis:     0.0,
				GyroscopeZAxis:     0.0,
				MagnetometerXAxis:  0.0,
				MagnetometerYAxis:  0.0,
				MagnetometerZAxis:  0.0,
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
				TimeToFix:          36,
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
				TimeToFix:          36,
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
		port     int16
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
			payload:  "0002c420ff005ed85a12b4180719142607240001ffbaffc2fc6f02621d2e",
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

	if err == nil || err.Error() != "payload too short" {
		t.Fatal("expected error payload too short")
	}
}

func TestPayloadTooLong(t *testing.T) {
	decoder := NewNomadXSv1Decoder()
	_, err := decoder.Decode("deadbeef4242deadbeef4242deadbeef4242deadbeef4242deadbeef4242deadbeef4242deadbeef4242deadbeef4242", 1, "")

	if err == nil || err.Error() != "payload too long" {
		t.Fatal("expected error payload too long")
	}
}

func TestFeatures(t *testing.T) {
	tests := []struct {
		payload        string
		port           int16
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
				batteryVoltage, ok := decodedPayload.Data.(decoder.UpLinkFeatureBattery)
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
		})
	}
}
