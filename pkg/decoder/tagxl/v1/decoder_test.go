package tagxl

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/truvami/decoder/pkg/decoder"
	"github.com/truvami/decoder/pkg/loracloud"
)

func startMockServer() *httptest.Server {
	server := httptest.NewServer(nil)
	return server
}

func TestDecode(t *testing.T) {

	http.HandleFunc("/api/v1/device/send", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// check if request body contains 10CE45FFFE00C7ED
		bodyString, _ := io.ReadAll(r.Body)
		if strings.Contains(string(bodyString), "10CE45FFFE00C7ED") {
			_, _ = w.Write([]byte("{\"invalid\": json}"))
			return
		}

		// get file from testdata
		file, err := os.Open("./response.json")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		_, _ = w.Write(data)
	})

	server := startMockServer()
	middleware := loracloud.NewLoracloudMiddleware("access_token")
	middleware.BaseUrl = server.URL
	defer server.Close()

	f, _ := os.Open("./response.json")
	var exampleResponse loracloud.UplinkMsgResponse
	d, _ := io.ReadAll(f)
	err := json.Unmarshal(d, &exampleResponse)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := []struct {
		payload     string
		port        uint8
		devEui      string
		autoPadding bool
		expected    any
		expectedErr string
	}{
		{
			payload:     "87821F50490200B520FBE977844D222A3A14A89293956245CC75A9CA1BBC25DDF658542909",
			port:        192,
			devEui:      "10CE45FFFE00C7EC",
			autoPadding: false,
			expected:    &exampleResponse,
		},
		{
			payload:     "87821F50490200B520FBE977844D222A3A14A89293956245CC75A9CA1BBC25DDF658542909",
			port:        192,
			devEui:      "10CE45FFFE00C7ED",
			autoPadding: false,
			expected:    nil,
			expectedErr: "invalid character 'j' looking for beginning of value",
		},
		{
			payload:     "00",
			port:        0,
			devEui:      "",
			autoPadding: false,
			expected:    nil,
			expectedErr: "port 0 not supported",
		},
		{
			payload: "0f0078012c000a03e80f1e0e9e7393fffe0002dead10cc0953046e",
			port:    151,
			expected: Port151Payload{
				Ble:                      false,
				Gnss:                     false,
				Wifi:                     false,
				Acceleration:             false,
				Rfu:                      15,
				MovingInterval:           120,
				SteadyInterval:           300,
				AccelerationThreshold:    10,
				AccelerationDelay:        1000,
				HeartbeatInterval:        15,
				FwuAdvertisementInterval: 30,
				BatteryVoltage:           3.742,
				FirmwareHash:             "7393fffe",
				ResetCount:               2,
				ResetCause:               3735883980,
				GnssScans:                2387,
				WifiScans:                1134,
			},
		},
		{
			payload: "57012c0258006403e81e3c0dbfd6ce814d0003c0debabe0acd04b9",
			port:    151,
			expected: Port151Payload{
				Ble:                      false,
				Gnss:                     true,
				Wifi:                     false,
				Acceleration:             true,
				Rfu:                      7,
				MovingInterval:           300,
				SteadyInterval:           600,
				AccelerationThreshold:    100,
				AccelerationDelay:        1000,
				HeartbeatInterval:        30,
				FwuAdvertisementInterval: 60,
				BatteryVoltage:           3.519,
				FirmwareHash:             "d6ce814d",
				ResetCount:               3,
				ResetCause:               3235822270,
				GnssScans:                2765,
				WifiScans:                1209,
			},
		},
		{
			payload: "a3025804b0012c05dc3c780d70ca6a55150005feedface0cae09b0",
			port:    151,
			expected: Port151Payload{
				Ble:                      true,
				Gnss:                     false,
				Wifi:                     true,
				Acceleration:             false,
				Rfu:                      3,
				MovingInterval:           600,
				SteadyInterval:           1200,
				AccelerationThreshold:    300,
				AccelerationDelay:        1500,
				HeartbeatInterval:        60,
				FwuAdvertisementInterval: 120,
				BatteryVoltage:           3.440,
				FirmwareHash:             "ca6a5515",
				ResetCount:               5,
				ResetCause:               4277009102,
				GnssScans:                3246,
				WifiScans:                2480,
			},
		},
		{
			payload: "f007080e1001c207d078f00c893113870f00088badf00d10000c00",
			port:    151,
			expected: Port151Payload{
				Ble:                      true,
				Gnss:                     true,
				Wifi:                     true,
				Acceleration:             true,
				Rfu:                      0,
				MovingInterval:           1800,
				SteadyInterval:           3600,
				AccelerationThreshold:    450,
				AccelerationDelay:        2000,
				HeartbeatInterval:        120,
				FwuAdvertisementInterval: 240,
				BatteryVoltage:           3.209,
				FirmwareHash:             "3113870f",
				ResetCount:               8,
				ResetCause:               2343432205,
				GnssScans:                4096,
				WifiScans:                3072,
			},
		},
		{
			payload:     "010b0266acbcf0000000000756",
			port:        152,
			autoPadding: false,
			expected: Port152Payload{
				NewRotationState:  2,
				OldRotationState:  0,
				Timestamp:         uint32(time.Date(2024, 8, 2, 11, 3, 12, 0, time.UTC).Unix()),
				NumberOfRotations: 0,
				ElapsedSeconds:    1878,
			},
		},
		{
			payload:     "10b0266acbcf0000000000756",
			port:        152,
			autoPadding: true,
			expected: Port152Payload{
				NewRotationState:  2,
				OldRotationState:  0,
				Timestamp:         uint32(time.Date(2024, 8, 2, 11, 3, 12, 0, time.UTC).Unix()),
				NumberOfRotations: 0,
				ElapsedSeconds:    1878,
			},
		},
		{
			payload:     "010b1066acbe0c00a200000087",
			port:        152,
			autoPadding: false,
			expected: Port152Payload{
				NewRotationState:  0,
				OldRotationState:  1,
				Timestamp:         uint32(time.Date(2024, 8, 2, 11, 7, 56, 0, time.UTC).Unix()),
				NumberOfRotations: 16.2,
				ElapsedSeconds:    135,
			},
		},
		{
			payload:     "10b1066acbe0c00a200000087",
			port:        152,
			autoPadding: true,
			expected: Port152Payload{
				NewRotationState:  0,
				OldRotationState:  1,
				Timestamp:         uint32(time.Date(2024, 8, 2, 11, 7, 56, 0, time.UTC).Unix()),
				NumberOfRotations: 16.2,
				ElapsedSeconds:    135,
			},
		},
		{
			payload: "00d63385f8ee30c2d0a0382c2601db",
			port:    197,
			expected: Port197Payload{
				Tag:   byte(0x00),
				Rssi1: -42,
				Mac1:  "3385f8ee30c2",
				Rssi2: -48,
				Mac2:  "a0382c2601db",
			},
		},
		{
			payload: "64c8b5eded55a313c0a0b8b5e86e31b894a765f3ad40",
			port:    197,
			expected: Port197Payload{
				Tag:   byte(0x64),
				Rssi1: -56,
				Mac1:  "b5eded55a313",
				Rssi2: -64,
				Mac2:  "a0b8b5e86e31",
				Rssi3: -72,
				Mac3:  "94a765f3ad40",
			},
		},
		{
			payload: "aebd6fbcfdd76434bb7e7cbff22fc5b900dc0af60588b7010161302d9c",
			port:    197,
			expected: Port197Payload{
				Tag:   byte(0xae),
				Rssi1: -67,
				Mac1:  "6fbcfdd76434",
				Rssi2: -69,
				Mac2:  "7e7cbff22fc5",
				Rssi3: -71,
				Mac3:  "00dc0af60588",
				Rssi4: -73,
				Mac4:  "010161302d9c",
			},
		},
		{
			payload: "fdb7218f6c166fadb359ea3bdec77daff72faac81784ab263386a455d3a73592a063900b",
			port:    197,
			expected: Port197Payload{
				Tag:   byte(0xfd),
				Rssi1: -73,
				Mac1:  "218f6c166fad",
				Rssi2: -77,
				Mac2:  "59ea3bdec77d",
				Rssi3: -81,
				Mac3:  "f72faac81784",
				Rssi4: -85,
				Mac4:  "263386a455d3",
				Rssi5: -89,
				Mac5:  "3592a063900b",
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestPort%vWith%v", test.port, test.payload), func(t *testing.T) {
			decoder := NewTagXLv1Decoder(middleware, WithAutoPadding(test.autoPadding), WithFCount(1))
			got, err := decoder.Decode(test.payload, test.port, test.devEui)
			if err != nil && len(test.expectedErr) == 0 {
				t.Fatalf("unexpected error: %v", err)
			}

			t.Logf("got %v", got)

			if got != nil && !reflect.DeepEqual(got.Data, test.expected) && len(test.expectedErr) == 0 {
				t.Errorf("expected: %v, got: %v", test.expected, got)
			}

			if len(test.expectedErr) > 0 && err != nil && !strings.Contains(err.Error(), test.expectedErr) {
				t.Errorf("expected error: %v, got: %v", test.expectedErr, err)
			}
		})
	}
}

func TestValidationErrors(t *testing.T) {
	tests := []struct {
		payload  string
		port     uint8
		expected error
	}{}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestPort%vValidationWith%v", test.port, test.payload), func(t *testing.T) {
			decoder := NewTagXLv1Decoder(loracloud.NewLoracloudMiddleware("apiKey"))
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
	decoder := NewTagXLv1Decoder(loracloud.NewLoracloudMiddleware("apiKey"))
	_, err := decoder.Decode("00", 0, "")
	if err == nil || err.Error() != "port 0 not supported" {
		t.Fatal("expected port not supported")
	}
}

func TestPayloadTooShort(t *testing.T) {
	decoder := NewTagXLv1Decoder(loracloud.NewLoracloudMiddleware("apiKey"))
	_, err := decoder.Decode("deadbeef", 152, "")

	if err == nil || err.Error() != "payload too short" {
		t.Fatal("expected error payload too short")
	}
}

func TestPayloadTooLong(t *testing.T) {
	decoder := NewTagXLv1Decoder(loracloud.NewLoracloudMiddleware("apiKey"))
	_, err := decoder.Decode("deadbeef4242deadbeef4242deadbeef4242", 152, "")

	if err == nil || err.Error() != "payload too long" {
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
			payload: "57012c0258006403e81e3c0dbfd6ce814d0003c0debabe0acd04b9",
			port:    151,
		},
		{
			payload: "010b0266acbcf0000000000756",
			port:    152,
		},
		{
			payload: "fdb7218f6c166fadb359ea3bdec77daff72faac81784ab263386a455d3a73592a063900b",
			port:    197,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestFeaturesWithPort%vAndPayload%v", test.port, test.payload), func(t *testing.T) {
			d := NewTagXLv1Decoder(
				loracloud.NewLoracloudMiddleware("apiKey"),
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
				config.GetBatchSize()
				config.GetBufferSize()
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
			payload:  "010b0266acbcf0000000000756",
			port:     152,
			expected: []string{"\"timestamp\": 1722596592", "\"elapsedSeconds\": 1878"},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestMarshalWithPort%vAndPayload%v", test.port, test.payload), func(t *testing.T) {
			decoder := NewTagXLv1Decoder(loracloud.NewLoracloudMiddleware("apiKey"))

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

func TestWithFCount(t *testing.T) {
	decoder := NewTagXLv1Decoder(loracloud.NewLoracloudMiddleware("apiKey"), WithFCount(123))

	// cast to TagXLv1Decoder to access fCount
	tagXLv1Decoder := decoder.(*TagXLv1Decoder)
	if tagXLv1Decoder.fCount != 123 {
		t.Fatalf("expected fCount to be 123, got %v", tagXLv1Decoder.fCount)
	}
}
