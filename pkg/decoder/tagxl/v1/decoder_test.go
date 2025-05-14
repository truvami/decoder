package tagxl

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	helpers "github.com/truvami/decoder/pkg/common"
	"github.com/truvami/decoder/pkg/decoder"
	"github.com/truvami/decoder/pkg/loracloud"
)

func startMockServer(handler http.Handler) *httptest.Server {
	server := httptest.NewServer(handler)
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

	server := startMockServer(nil)
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
			port:        0,
			payload:     "00",
			devEui:      "",
			autoPadding: false,
			expected:    nil,
			expectedErr: "port 0 not supported",
		},
		{
			port:    150,
			payload: "4c07014c04681a4727",
			expected: Port150Payload{
				Timestamp: time.Date(2025, 5, 6, 17, 30, 15, 0, time.UTC),
			},
		},
		{
			port:    150,
			payload: "4c07014c04681a5127",
			expected: Port150Payload{
				Timestamp: time.Date(2025, 5, 6, 18, 12, 55, 0, time.UTC),
			},
		},
		{
			port:    151,
			payload: "4c050145020a92",
			expected: Port151Payload{
				Battery: helpers.Float32Ptr(2.706),
			},
		},
		{
			port:    151,
			payload: "4c050145020a93",
			expected: Port151Payload{
				Battery: helpers.Float32Ptr(2.707),
			},
		},
		{
			port:    151,
			payload: "4c050145020a96",
			expected: Port151Payload{
				Battery: helpers.Float32Ptr(2.710),
			},
		},
		{
			port:    151,
			payload: "4c050145020b10",
			expected: Port151Payload{
				Battery: helpers.Float32Ptr(2.832),
			},
		},
		{
			port:    151,
			payload: "4c0b0245020d914b0403de0000",
			expected: Port151Payload{
				Battery:   helpers.Float32Ptr(3.473),
				GnssScans: helpers.Uint16Ptr(990),
				WifiScans: helpers.Uint16Ptr(0),
			},
		},
		{
			port:    151,
			payload: "4c0b0245020d7b4b0404fc0000",
			expected: Port151Payload{
				Battery:   helpers.Float32Ptr(3.451),
				GnssScans: helpers.Uint16Ptr(1276),
				WifiScans: helpers.Uint16Ptr(0),
			},
		},
		{
			port:    151,
			payload: "4c0b0245020db94b0400420000",
			expected: Port151Payload{
				Battery:   helpers.Float32Ptr(3.513),
				GnssScans: helpers.Uint16Ptr(66),
				WifiScans: helpers.Uint16Ptr(0),
			},
		},
		{
			port:    151,
			payload: "4c0b0245020d8f4b0401970000",
			expected: Port151Payload{
				Battery:   helpers.Float32Ptr(3.471),
				GnssScans: helpers.Uint16Ptr(407),
				WifiScans: helpers.Uint16Ptr(0),
			},
		},
		{
			port:    152,
			payload: "020c62206822f120000d00000024",
			expected: Port152Payload{
				Version:           2,
				SequenceNumber:    98,
				NewRotationState:  0,
				OldRotationState:  2,
				Timestamp:         time.Date(2025, 5, 13, 7, 13, 36, 0, time.UTC),
				NumberOfRotations: 1.3,
				ElapsedSeconds:    36,
			},
		},
		{
			port:    152,
			payload: "020c09016823166a000000000109",
			expected: Port152Payload{
				Version:           2,
				SequenceNumber:    9,
				NewRotationState:  1,
				OldRotationState:  0,
				Timestamp:         time.Date(2025, 5, 13, 9, 52, 42, 0, time.UTC),
				NumberOfRotations: 0.0,
				ElapsedSeconds:    265,
			},
		},
		{
			port:    152,
			payload: "020cea0268230e60000000000015",
			expected: Port152Payload{
				Version:           2,
				SequenceNumber:    234,
				NewRotationState:  2,
				OldRotationState:  0,
				Timestamp:         time.Date(2025, 5, 13, 9, 18, 24, 0, time.UTC),
				NumberOfRotations: 0.0,
				ElapsedSeconds:    21,
			},
		},
		{
			port:    152,
			payload: "010b0266acbcf0000000000756",
			expected: Port152Payload{
				Version:           1,
				NewRotationState:  2,
				OldRotationState:  0,
				Timestamp:         time.Date(2024, 8, 2, 11, 3, 12, 0, time.UTC),
				NumberOfRotations: 0,
				ElapsedSeconds:    1878,
			},
		},
		{
			port:    152,
			payload: "010b1066acbe0c00a200000087",
			expected: Port152Payload{
				Version:           1,
				NewRotationState:  0,
				OldRotationState:  1,
				Timestamp:         time.Date(2024, 8, 2, 11, 7, 56, 0, time.UTC),
				NumberOfRotations: 16.2,
				ElapsedSeconds:    135,
			},
		},
		{
			port:        192,
			payload:     "87821f50490200b520fbe977844d222a3a14a89293956245cc75a9ca1bbc25ddf658542909",
			devEui:      "10CE45FFFE00C7EC",
			autoPadding: false,
			expected:    &exampleResponse,
		},
		{
			port:        192,
			payload:     "87821f50490200b520fbe977844d222a3a14a89293956245cc75a9ca1bbc25ddf658542909",
			devEui:      "10CE45FFFE00C7ED",
			autoPadding: false,
			expected:    nil,
			expectedErr: "invalid character 'j' looking for beginning of value",
		},
		{
			port:    197,
			payload: "00d63385f8ee30c2d0a0382c2601db",
			expected: Port197Payload{
				Tag:   byte(0x00),
				Rssi1: -42,
				Mac1:  "3385f8ee30c2",
				Rssi2: -48,
				Mac2:  "a0382c2601db",
			},
		},
		{
			port:    197,
			payload: "64c8b5eded55a313c0a0b8b5e86e31b894a765f3ad40",
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
			port:    197,
			payload: "aebd6fbcfdd76434bb7e7cbff22fc5b900dc0af60588b7010161302d9c",
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
			port:    197,
			payload: "fdb7218f6c166fadb359ea3bdec77daff72faac81784ab263386a455d3a73592a063900b",
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

	if err == nil || !errors.Is(err, helpers.ErrPortNotSupported) {
		t.Fatal("expected port not supported")
	}
}

func TestPayloadTooShort(t *testing.T) {
	decoder := NewTagXLv1Decoder(loracloud.NewLoracloudMiddleware("apiKey"))
	_, err := decoder.Decode("01adbeef", 152, "")

	if err == nil || !errors.Is(err, helpers.ErrPayloadTooShort) {
		t.Fatalf("expected error payload too short but got %v", err)
	}
}

func TestPayloadTooLong(t *testing.T) {
	decoder := NewTagXLv1Decoder(loracloud.NewLoracloudMiddleware("apiKey"))
	_, err := decoder.Decode("01adbeef4242deadbeef4242deadbeef4242", 152, "")

	if err == nil || !errors.Is(err, helpers.ErrPayloadTooLong) {
		t.Fatal("expected error payload too long")
	}
}

func TestFeatures(t *testing.T) {
	tests := []struct {
		payload         string
		port            uint8
		allowNoFeatures bool
	}{
		{
			payload:         "4c07014c04681a5127",
			port:            150,
			allowNoFeatures: true,
		},
		{
			port:    151,
			payload: "4c050145020b10",
		},
		{
			payload: "010b0266acbcf0000000000756",
			port:    152,
		},
		{
			payload: "020c62206822f120000d00000024",
			port:    152,
		},
		{
			payload: "020c09016823166a000000000109",
			port:    152,
		},
		{
			payload: "87821f50490200b520fbe977844d222a3a14a89293956245cc75a9ca1bbc25ddf658542909",
			port:    192,
		},
		{
			payload: "fdb7218f6c166fadb359ea3bdec77daff72faac81784ab263386a455d3a73592a063900b",
			port:    197,
		},
		{
			payload:         "86b5277140484a89b8f63ccf67affbfeb519b854f9d447808a50785bdfe86a77",
			port:            199,
			allowNoFeatures: true,
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/device/send", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
				"result": {
					"deveui": "927da4b72110927d",
					"position_solution": {
							"llh": [51.49278, 0.0212, 83.93],
							"accuracy": 20.7,
							"gdop": 2.48,
							"capture_time_utc": 1722433373.18046
					},
					"operation": "gnss"
				}
			}`))
	})

	server := startMockServer(mux)
	middleware := loracloud.NewLoracloudMiddleware("access_token")
	middleware.BaseUrl = server.URL
	defer server.Close()

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestFeaturesWithPort%vAndPayload%v", test.port, test.payload), func(t *testing.T) {
			d := NewTagXLv1Decoder(middleware, WithFCount(42))
			decodedPayload, err := d.Decode(test.payload, test.port, "927da4b72110927d")
			if err != nil {
				t.Fatalf("error %s", err)
			}

			// should be able to decode base feature
			base, ok := decodedPayload.Data.(decoder.UplinkFeatureBase)
			if !ok {
				t.Fatalf("expected UplinkFeatureBase, got %T", decodedPayload)
			}
			// check if it panics
			base.GetTimestamp()

			if len(decodedPayload.GetFeatures()) == 0 && !test.allowNoFeatures {
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
				config.GetDataRate()
			}
			if decodedPayload.Is(decoder.FeatureRotationState) {
				rotationState, ok := decodedPayload.Data.(decoder.UplinkFeatureRotationState)
				if !ok {
					t.Fatalf("expected UplinkFeatureRotationState, got %T", decodedPayload)
				}
				// call function to check if it panics
				rotationState.GetRotationState()
			}
			if decodedPayload.Is(decoder.FeatureSequenceNumber) {
				sequenceNumber, ok := decodedPayload.Data.(decoder.UplinkFeatureSequenceNumber)
				if !ok {
					t.Fatalf("expected UplinkFeatureSequenceNumber, got %T", decodedPayload)
				}
				// call function to check if it panics
				sequenceNumber.GetSequenceNumber()
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
			expected: []string{"\"timestamp\": \"2024-08-02T11:03:12Z\"", "\"elapsedSeconds\": 1878"},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestMarshalWithPort%vAndPayload%v", test.port, test.payload), func(t *testing.T) {
			decoder := NewTagXLv1Decoder(loracloud.NewLoracloudMiddleware("apiKey"))

			data, _ := decoder.Decode(test.payload, test.port, "")

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

func TestWithFCount(t *testing.T) {
	decoder := NewTagXLv1Decoder(loracloud.NewLoracloudMiddleware("apiKey"), WithFCount(123))

	// cast to TagXLv1Decoder to access fCount
	tagXLv1Decoder := decoder.(*TagXLv1Decoder)
	if tagXLv1Decoder.fCount != 123 {
		t.Fatalf("expected fCount to be 123, got %v", tagXLv1Decoder.fCount)
	}
}
