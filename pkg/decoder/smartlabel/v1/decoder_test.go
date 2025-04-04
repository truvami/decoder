package smartlabel

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
		expected    any
		expectedErr string
	}{
		{
			payload: "0ca90dbd",
			port:    1,
			expected: Port1Payload{
				BatteryVoltage:      3.241,
				PhotovoltaicVoltage: 3.517,
			},
		},
		{
			payload: "0dfa0e1a",
			port:    1,
			expected: Port1Payload{
				BatteryVoltage:      3.578,
				PhotovoltaicVoltage: 3.610,
			},
		},
		{
			payload: "0e860ef7",
			port:    1,
			expected: Port1Payload{
				BatteryVoltage:      3.718,
				PhotovoltaicVoltage: 3.831,
			},
		},
		{
			payload: "0f501079",
			port:    1,
			expected: Port1Payload{
				BatteryVoltage:      3.920,
				PhotovoltaicVoltage: 4.217,
			},
		},
		{
			payload: "07fa69",
			port:    2,
			expected: Port2Payload{
				Temperature: 20.42,
				Humidity:    52.5,
			},
		},
		{
			payload: "074070",
			port:    2,
			expected: Port2Payload{
				Temperature: 18.56,
				Humidity:    56.0,
			},
		},
		{
			payload: "06947d",
			port:    2,
			expected: Port2Payload{
				Temperature: 16.84,
				Humidity:    62.5,
			},
		},
		{
			payload: "04da8d",
			port:    2,
			expected: Port2Payload{
				Temperature: 12.42,
				Humidity:    70.5,
			},
		},
		{
			payload: "00012c00780f000a03e8003c0078f81c01000402",
			port:    4,
			expected: Port4Payload{
				DataRate:                   0,
				Acceleration:               false,
				Wifi:                       false,
				Gnss:                       false,
				SteadyInterval:             300,
				MovingInterval:             120,
				HeartbeatInterval:          15,
				AccelerationThreshold:      10,
				AccelerationDelay:          1000,
				TemperaturePollingInterval: 60,
				TemperatureUplinkInterval:  120,
				TemperatureLowerThreshold:  -8,
				TemperatureUpperThreshold:  +28,
				AccessPointsThreshold:      1,
				FirmwareVersionMajor:       0,
				FirmwareVersionMinor:       4,
				FirmwareVersionPatch:       2,
			},
		},
		{
			payload: "2a0258012c1e006403e8003c012cf61e02010307",
			port:    4,
			expected: Port4Payload{
				DataRate:                   2,
				Acceleration:               true,
				Wifi:                       false,
				Gnss:                       true,
				SteadyInterval:             600,
				MovingInterval:             300,
				HeartbeatInterval:          30,
				AccelerationThreshold:      100,
				AccelerationDelay:          1000,
				TemperaturePollingInterval: 60,
				TemperatureUplinkInterval:  300,
				TemperatureLowerThreshold:  -10,
				TemperatureUpperThreshold:  +30,
				AccessPointsThreshold:      2,
				FirmwareVersionMajor:       1,
				FirmwareVersionMinor:       3,
				FirmwareVersionPatch:       7,
			},
		},
		{
			payload: "1b04b002583c012c05dc003c0258f12304020400",
			port:    4,
			expected: Port4Payload{
				DataRate:                   3,
				Acceleration:               true,
				Wifi:                       true,
				Gnss:                       false,
				SteadyInterval:             1200,
				MovingInterval:             600,
				HeartbeatInterval:          60,
				AccelerationThreshold:      300,
				AccelerationDelay:          1500,
				TemperaturePollingInterval: 60,
				TemperatureUplinkInterval:  600,
				TemperatureLowerThreshold:  -15,
				TemperatureUpperThreshold:  +35,
				AccessPointsThreshold:      4,
				FirmwareVersionMajor:       2,
				FirmwareVersionMinor:       4,
				FirmwareVersionPatch:       0,
			},
		},
		{
			payload: "3f0e1007087801c207d0003c04b0ec280603020c",
			port:    4,
			expected: Port4Payload{
				DataRate:                   7,
				Acceleration:               true,
				Wifi:                       true,
				Gnss:                       true,
				SteadyInterval:             3600,
				MovingInterval:             1800,
				HeartbeatInterval:          120,
				AccelerationThreshold:      450,
				AccelerationDelay:          2000,
				TemperaturePollingInterval: 60,
				TemperatureUplinkInterval:  1200,
				TemperatureLowerThreshold:  -20,
				TemperatureUpperThreshold:  +40,
				AccessPointsThreshold:      6,
				FirmwareVersionMajor:       3,
				FirmwareVersionMinor:       2,
				FirmwareVersionPatch:       12,
			},
		},
		{
			payload: "0ca90dbd07fa69",
			port:    11,
			expected: Port11Payload{
				BatteryVoltage:      3.241,
				PhotovoltaicVoltage: 3.517,
				Temperature:         20.42,
				Humidity:            52.5,
			},
		},
		{
			payload: "0dfa0e1a074070",
			port:    11,
			expected: Port11Payload{
				BatteryVoltage:      3.578,
				PhotovoltaicVoltage: 3.610,
				Temperature:         18.56,
				Humidity:            56.0,
			},
		},
		{
			payload: "0e860ef706947d",
			port:    11,
			expected: Port11Payload{
				BatteryVoltage:      3.718,
				PhotovoltaicVoltage: 3.831,
				Temperature:         16.84,
				Humidity:            62.5,
			},
		},
		{
			payload: "0f50107904da8d",
			port:    11,
			expected: Port11Payload{
				BatteryVoltage:      3.920,
				PhotovoltaicVoltage: 4.217,
				Temperature:         12.42,
				Humidity:            70.5,
			},
		},
		{
			payload:  "87821F50490200B520FBE977844D222A3A14A89293956245CC75A9CA1BBC25DDF658542909",
			port:     192,
			devEui:   "10CE45FFFE00C7EC",
			expected: &exampleResponse,
		},
		{
			payload:     "87821F50490200B520FBE977844D222A3A14A89293956245CC75A9CA1BBC25DDF658542909",
			port:        192,
			devEui:      "10CE45FFFE00C7ED",
			expected:    nil,
			expectedErr: "invalid character 'j' looking for beginning of value",
		},
		{
			payload:     "00",
			port:        0,
			devEui:      "",
			expected:    nil,
			expectedErr: "port 0 not supported",
		},
		{
			payload: "0e1a0db60d520c260a96",
			port:    150,
			expected: Port150Payload{
				Battery100Voltage: 3.610,
				Battery80Voltage:  3.510,
				Battery60Voltage:  3.410,
				Battery40Voltage:  3.110,
				Battery20Voltage:  2.710,
			},
		},
		{
			payload: "0ed80e420dac0cb20b22",
			port:    150,
			expected: Port150Payload{
				Battery100Voltage: 3.800,
				Battery80Voltage:  3.650,
				Battery60Voltage:  3.500,
				Battery40Voltage:  3.250,
				Battery20Voltage:  2.850,
			},
		},
		{
			payload: "0f960e6a0e060d3e0bae",
			port:    150,
			expected: Port150Payload{
				Battery100Voltage: 3.990,
				Battery80Voltage:  3.690,
				Battery60Voltage:  3.590,
				Battery40Voltage:  3.390,
				Battery20Voltage:  2.990,
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
			payload: "aebd6fbcfdd76434bb7e7cbff22fc5b900dc0af60588b7010161302d9cb51bf1f8d1a97b",
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
				Rssi5: -75,
				Mac5:  "1bf1f8d1a97b",
			},
		},
		{
			payload: "fdb7218f6c166fadb359ea3bdec77daff72faac81784ab263386a455d3a73592a063900ba262b95a6ffc86",
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
				Rssi6: -94,
				Mac6:  "62b95a6ffc86",
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestPort%vWith%v", test.port, test.payload), func(t *testing.T) {
			decoder := NewSmartLabelv1Decoder(middleware, WithFCount(1))
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

func TestInvalidPort(t *testing.T) {
	decoder := NewSmartLabelv1Decoder(loracloud.NewLoracloudMiddleware("appEui"))
	_, err := decoder.Decode("00", 0, "")
	if err == nil || err.Error() != "port 0 not supported" {
		t.Fatal("expected port not supported")
	}
}

func TestPayloadTooShort(t *testing.T) {
	decoder := NewSmartLabelv1Decoder(loracloud.NewLoracloudMiddleware("appEui"))
	_, err := decoder.Decode("0ff0", 1, "")

	if err == nil || err.Error() != "payload too short" {
		t.Fatal("expected error payload too short")
	}
}

func TestPayloadTooLong(t *testing.T) {
	decoder := NewSmartLabelv1Decoder(loracloud.NewLoracloudMiddleware("appEui"))
	_, err := decoder.Decode("0ff00ff00ff0", 1, "")

	if err == nil || err.Error() != "payload too long" {
		t.Fatal("expected error payload too long")
	}
}

func TestFeatures(t *testing.T) {
	tests := []struct {
		payload string
		port    uint8
	}{
		{
			payload: "0f501079",
			port:    1,
		},
		{
			payload: "04da8d",
			port:    2,
		},
		{
			payload: "3f0e1007087801c207d0003c04b0ec280603020c",
			port:    4,
		},
		{
			payload: "0f50107904da8d",
			port:    11,
		},
		{
			payload: "0ed80e420dac0cb20b22",
			port:    150,
		},
		{
			payload: "87821f50490200b520fbe977844d222a3a14a89293956245cc75a9ca1bbc25ddf658542909",
			port:    192,
		},
		{
			payload: "fdb7218f6c166fadb359ea3bdec77daff72faac81784ab263386a455d3a73592a063900ba262b95a6ffc86",
			port:    197,
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
			d := NewSmartLabelv1Decoder(
				middleware,
				WithFCount(42),
			)
			data, err := d.Decode(test.payload, test.port, "927da4b72110927d")
			if err != nil {
				t.Fatalf("error %s", err)
			}

			// should be able to decode base feature
			base, ok := data.Data.(decoder.UplinkFeatureBase)
			if !ok {
				t.Fatalf("expected UplinkFeatureBase, got %T", data)
			}
			// check if it panics
			base.GetTimestamp()

			if data.Is(decoder.FeatureGNSS) {
				gnss, ok := data.Data.(decoder.UplinkFeatureGNSS)
				if !ok {
					t.Fatalf("expected UplinkFeatureGNSS, got %T", data)
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
			if data.Is(decoder.FeatureBattery) {
				batteryVoltage, ok := data.Data.(decoder.UplinkFeatureBattery)
				if !ok {
					t.Fatalf("expected UplinkFeatureBattery, got %T", data)
				}
				if batteryVoltage.GetBatteryVoltage() == 0 {
					t.Fatalf("expected non zero battery voltage")
				}
			}
			if data.Is(decoder.FeatureWiFi) {
				wifi, ok := data.Data.(decoder.UplinkFeatureWiFi)
				if !ok {
					t.Fatalf("expected UplinkFeatureWiFi, got %T", data)
				}
				if wifi.GetAccessPoints() == nil {
					t.Fatalf("expected non nil access points")
				}
			}
			if data.Is(decoder.FeaturePhotovoltaic) {
				photovoltaicVoltage, ok := data.Data.(decoder.UplinkFeaturePhotovoltaic)
				if !ok {
					t.Fatalf("expected UplinkFeaturePhotovoltaic, got %T", data)
				}
				if photovoltaicVoltage.GetPhotovoltaicVoltage() == 0 {
					t.Fatalf("expected non zero photovoltaic voltage")
				}
			}
			if data.Is(decoder.FeatureTemperature) {
				temperature, ok := data.Data.(decoder.UplinkFeatureTemperature)
				if !ok {
					t.Fatalf("expected UplinkFeatureTemperature, got %T", data)
				}
				if temperature.GetTemperature() == 0 {
					t.Fatalf("expected non zero temperature")
				}
			}
			if data.Is(decoder.FeatureHumidity) {
				humidity, ok := data.Data.(decoder.UplinkFeatureHumidity)
				if !ok {
					t.Fatalf("expected UplinkFeatureHumidity, got %T", data)
				}
				if humidity.GetHumidity() == 0 {
					t.Fatalf("expected non zero humidity")
				}
			}
			if data.Is(decoder.FeatureConfig) {
				config, ok := data.Data.(decoder.UplinkFeatureConfig)
				if !ok {
					t.Fatalf("expected UplinkFeatureConfig, got %T", data)
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
			if data.Is(decoder.FeatureFirmwareVersion) {
				firmwareVersion, ok := data.Data.(decoder.UplinkFeatureFirmwareVersion)
				if !ok {
					t.Fatalf("expected UplinkFeatureFirmwareVersion, got %T", data)
				}
				if firmwareVersion.GetFirmwareVersion() == "" {
					t.Fatalf("expected non empty firmware version")
				}
			}
		})
	}
}

func TestWithAutoPadding(t *testing.T) {
	middleware := loracloud.NewLoracloudMiddleware("access_token")

	decoder := NewSmartLabelv1Decoder(
		middleware,
		WithAutoPadding(true),
	)

	// Type assert to access the internal field
	if d, ok := decoder.(*SmartLabelv1Decoder); ok {
		if !d.autoPadding {
			t.Error("expected autoPadding to be true")
		}
	} else {
		t.Error("failed to type assert decoder")
	}
}

func TestGetPort11PayloadType(t *testing.T) {
	tests := []struct {
		name        string
		data        string
		want        string
		expectedErr string
	}{
		{
			name:        "Empty data",
			data:        "",
			want:        "",
			expectedErr: "data length is less than 2",
		},
		{
			name:        "Single byte",
			data:        "0",
			want:        "",
			expectedErr: "data length is less than 2",
		},
		{
			name:        "Configuration payload 0E",
			data:        "0E",
			want:        "configuration",
			expectedErr: "",
		},
		{
			name:        "Configuration payload 11",
			data:        "11",
			want:        "configuration",
			expectedErr: "",
		},
		{
			name:        "Heartbeat payload",
			data:        "0A",
			want:        "heartbeat",
			expectedErr: "",
		},
		{
			name:        "Invalid payload",
			data:        "FF",
			want:        "",
			expectedErr: "invalid payload for port 11",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getPort11PayloadType(tt.data)
			if err != nil && tt.expectedErr == "" {
				t.Errorf("unexpected error: %v", err)
			}
			if err == nil && tt.expectedErr != "" {
				t.Errorf("expected error containing %q, got nil", tt.expectedErr)
			}
			if err != nil && tt.expectedErr != "" && !strings.Contains(err.Error(), tt.expectedErr) {
				t.Errorf("expected error containing %q, got %q", tt.expectedErr, err.Error())
			}
			if got != tt.want {
				t.Errorf("getPort11PayloadType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithFCount(t *testing.T) {
	decoder := NewSmartLabelv1Decoder(loracloud.NewLoracloudMiddleware("apiKey"), WithFCount(123))

	// cast to SmartLabelv1Decoder to access fCount
	tagXLv1Decoder := decoder.(*SmartLabelv1Decoder)
	if tagXLv1Decoder.fCount != 123 {
		t.Fatalf("expected fCount to be 123, got %v", tagXLv1Decoder.fCount)
	}
}
