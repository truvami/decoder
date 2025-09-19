package tagxl

import (
	"context"
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

	"github.com/stretchr/testify/assert"
	"github.com/truvami/decoder/internal/logger"
	helpers "github.com/truvami/decoder/pkg/common"
	"github.com/truvami/decoder/pkg/decoder"
	"github.com/truvami/decoder/pkg/solver"
	"github.com/truvami/decoder/pkg/solver/aws"
	"github.com/truvami/decoder/pkg/solver/loracloud"
	"go.uber.org/zap"
)

func startMockServer(handler http.Handler) *httptest.Server {
	server := httptest.NewServer(handler)
	return server
}

func TestDecode(t *testing.T) {

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/device/send", func(w http.ResponseWriter, r *http.Request) {
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

	server := startMockServer(mux)
	middleware, err := loracloud.NewLoracloudClient(context.TODO(), "access_token", zap.NewExample())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	middleware.BaseUrl = server.URL
	defer server.Close()

	f, _ := os.Open("./response.json")
	var exampleResponse loracloud.UplinkMsgResponse
	d, _ := io.ReadAll(f)
	err = json.Unmarshal(d, &exampleResponse)
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
			port:        0,
			payload:     "00",
			devEui:      "",
			expected:    nil,
			expectedErr: "port 0 not supported",
		},
		{
			port:        150,
			payload:     "xx",
			devEui:      "",
			expected:    nil,
			expectedErr: "encoding/hex: invalid byte: U+0078 'x'",
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
			port:        151,
			payload:     "ff",
			expected:    Port151Payload{},
			expectedErr: "port not supported: port 151 tag ff",
		},
		{
			port:        151,
			payload:     "4c0501ff020000",
			expected:    Port151Payload{},
			expectedErr: "unknown tag ff",
		},
		{
			port:    151,
			payload: "4c040140010a",
			expected: Port151Payload{
				AccelerometerEnabled: helpers.BoolPtr(true),
				WifiEnabled:          helpers.BoolPtr(false),
				GnssEnabled:          helpers.BoolPtr(true),
				FirmwareUpgrade:      helpers.BoolPtr(false),
			},
		},
		{
			port:    151,
			payload: "4c040140010e",
			expected: Port151Payload{
				AccelerometerEnabled: helpers.BoolPtr(true),
				WifiEnabled:          helpers.BoolPtr(true),
				GnssEnabled:          helpers.BoolPtr(true),
				FirmwareUpgrade:      helpers.BoolPtr(false),
			},
		},
		{
			port:    151,
			payload: "4c0401400103",
			expected: Port151Payload{
				AccelerometerEnabled: helpers.BoolPtr(false),
				WifiEnabled:          helpers.BoolPtr(false),
				GnssEnabled:          helpers.BoolPtr(true),
				FirmwareUpgrade:      helpers.BoolPtr(true),
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
			port:    151,
			payload: "4c2a0940010f4104012c1c204204012c05dc43010644011e45020d4e4604f6c7d8104902000a4a0400000002",
			expected: Port151Payload{
				AccelerometerEnabled:                 helpers.BoolPtr(true),
				WifiEnabled:                          helpers.BoolPtr(true),
				GnssEnabled:                          helpers.BoolPtr(true),
				FirmwareUpgrade:                      helpers.BoolPtr(true),
				LocalizationIntervalWhileMoving:      helpers.Uint16Ptr(300),
				LocalizationIntervalWhileSteady:      helpers.Uint16Ptr(7200),
				AccelerometerWakeupThreshold:         helpers.Uint16Ptr(300),
				AccelerometerDelay:                   helpers.Uint16Ptr(1500),
				HeartbeatInterval:                    helpers.Uint8Ptr(6),
				AdvertisementFirmwareUpgradeInterval: helpers.Uint8Ptr(30),
				Battery:                              helpers.Float32Ptr(3.406),
				FirmwareHash:                         helpers.StringPtr("f6c7d810"),
				ResetCount:                           helpers.Uint16Ptr(10),
				ResetCause:                           helpers.Uint32Ptr(2),
			},
		},
		{
			port:    151,
			payload: "4c2d0a40010b410402581c204204012c05dc43010644011e45020d6c4604a25b545547010249020003",
			expected: Port151Payload{
				AccelerometerEnabled:                 helpers.BoolPtr(true),
				WifiEnabled:                          helpers.BoolPtr(false),
				GnssEnabled:                          helpers.BoolPtr(true),
				FirmwareUpgrade:                      helpers.BoolPtr(true),
				LocalizationIntervalWhileMoving:      helpers.Uint16Ptr(600),
				LocalizationIntervalWhileSteady:      helpers.Uint16Ptr(7200),
				AccelerometerWakeupThreshold:         helpers.Uint16Ptr(300),
				AccelerometerDelay:                   helpers.Uint16Ptr(1500),
				HeartbeatInterval:                    helpers.Uint8Ptr(6),
				AdvertisementFirmwareUpgradeInterval: helpers.Uint8Ptr(30),
				Battery:                              helpers.Float32Ptr(3.436),
				FirmwareHash:                         helpers.StringPtr("a25b5455"),
				RotationInvert:                       helpers.BoolPtr(false),
				RotationConfirmed:                    helpers.BoolPtr(true),
				ResetCount:                           helpers.Uint16Ptr(3),
			},
		},
		{
			port:        152,
			payload:     "ff",
			expected:    Port152Payload{},
			expectedErr: "port not supported: version 255 for port 152 not supported",
		},
		{
			port:    152,
			payload: "020c62206822f120000d00000024",
			expected: Port152Payload{
				Version:           2,
				SequenceNumber:    98,
				OldRotationState:  2,
				NewRotationState:  0,
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
				OldRotationState:  0,
				NewRotationState:  1,
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
				OldRotationState:  0,
				NewRotationState:  2,
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
				OldRotationState:  0,
				NewRotationState:  2,
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
				OldRotationState:  1,
				NewRotationState:  0,
				Timestamp:         time.Date(2024, 8, 2, 11, 7, 56, 0, time.UTC),
				NumberOfRotations: 16.2,
				ElapsedSeconds:    135,
			},
		},
		{
			port:     192,
			payload:  "87821f50490200b520fbe977844d222a3a14a89293956245cc75a9ca1bbc25ddf658542909",
			devEui:   "10CE45FFFE00C7EC",
			expected: &exampleResponse,
		},
		{
			port:        192,
			payload:     "87821f50490200b520fbe977844d222a3a14a89293956245cc75a9ca1bbc25ddf658542909",
			devEui:      "10CE45FFFE00C7ED",
			expected:    &exampleResponse,
			expectedErr: "",
		},
		{
			port:     194,
			payload:  "68b9b2318f2b157de4733aa4d27b5d3b3c6ecc9460a20a196b754655c98607",
			expected: &exampleResponse,
		},
		{
			port:     194,
			payload:  "68bad32509ab91418ae63a10b5004a0a3fef037ab2f06ce8e510820c1a0bdcecb49e1543fdd2f28f1c",
			expected: &exampleResponse,
		},
		{
			port:     194,
			payload:  "68bad32589b379e7ba0fb5006b9aaa8c8e25febf16f4e5c31d0cc8ca12a1cffdddf16c2cf82877f1edee4ecbc5ef54",
			expected: &exampleResponse,
		},
		{
			port:     195,
			payload:  "68bad3c50aabd56cb2e7ba0db5805a5ac9d4edd8de8a021b4ae2b78e8c0b8391566ab8d47d1d4c55ae794a2c2da7a637b49d32e44800",
			expected: &exampleResponse,
		},
		{
			port:     195,
			payload:  "68bad3c58aab4581b9e73a0eb580da120d7f85a75e770c6acad3dc2acdacbdcd576ab8147f5902557379b18d0f676a35fb9a6ae5ee03",
			expected: &exampleResponse,
		},
		{
			port:        197,
			payload:     "ff",
			expected:    Port197Payload{},
			expectedErr: "port not supported: version 255 for port 197 not supported",
		},
		{
			port:    197,
			payload: "003385f8ee30c2",
			expected: Port197Payload{
				Rssi1:   nil,
				Mac1:    "3385f8ee30c2",
				Moving:  false,
				Version: Port197Version1,
			},
		},
		{
			port:    197,
			payload: "003385f8ee30c2a0382c2601db",
			expected: Port197Payload{
				Rssi1:   nil,
				Mac1:    "3385f8ee30c2",
				Mac2:    helpers.StringPtr("a0382c2601db"),
				Moving:  false,
				Version: Port197Version1,
			},
		},
		{
			port:    197,
			payload: "00b5eded55a313a0b8b5e86e3194a765f3ad40",
			expected: Port197Payload{
				Rssi1:   nil,
				Mac1:    "b5eded55a313",
				Mac2:    helpers.StringPtr("a0b8b5e86e31"),
				Mac3:    helpers.StringPtr("94a765f3ad40"),
				Moving:  false,
				Version: Port197Version1,
			},
		},
		{
			port:    197,
			payload: "006fbcfdd764347e7cbff22fc500dc0af60588010161302d9c",
			expected: Port197Payload{
				Rssi1:   nil,
				Mac1:    "6fbcfdd76434",
				Mac2:    helpers.StringPtr("7e7cbff22fc5"),
				Mac3:    helpers.StringPtr("00dc0af60588"),
				Mac4:    helpers.StringPtr("010161302d9c"),
				Moving:  false,
				Version: Port197Version1,
			},
		},
		{
			port:    197,
			payload: "00218f6c166fad59ea3bdec77df72faac81784263386a455d33592a063900b",
			expected: Port197Payload{
				Rssi1:   nil,
				Mac1:    "218f6c166fad",
				Mac2:    helpers.StringPtr("59ea3bdec77d"),
				Mac3:    helpers.StringPtr("f72faac81784"),
				Mac4:    helpers.StringPtr("263386a455d3"),
				Mac5:    helpers.StringPtr("3592a063900b"),
				Moving:  false,
				Version: Port197Version1,
			},
		},
		{
			port:    197,
			payload: "01d63385f8ee30c2",
			expected: Port197Payload{
				Rssi1:   helpers.Int8Ptr(-42),
				Mac1:    "3385f8ee30c2",
				Moving:  false,
				Version: Port197Version2,
			},
		},
		{
			port:    197,
			payload: "01d63385f8ee30c2d0a0382c2601db",
			expected: Port197Payload{
				Rssi1:   helpers.Int8Ptr(-42),
				Mac1:    "3385f8ee30c2",
				Rssi2:   helpers.Int8Ptr(-48),
				Mac2:    helpers.StringPtr("a0382c2601db"),
				Moving:  false,
				Version: Port197Version2,
			},
		},
		{
			port:    197,
			payload: "01c8b5eded55a313c0a0b8b5e86e31b894a765f3ad40",
			expected: Port197Payload{
				Rssi1:   helpers.Int8Ptr(-56),
				Mac1:    "b5eded55a313",
				Rssi2:   helpers.Int8Ptr(-64),
				Mac2:    helpers.StringPtr("a0b8b5e86e31"),
				Rssi3:   helpers.Int8Ptr(-72),
				Mac3:    helpers.StringPtr("94a765f3ad40"),
				Moving:  false,
				Version: Port197Version2,
			},
		},
		{
			port:    197,
			payload: "01bd6fbcfdd76434bb7e7cbff22fc5b900dc0af60588b7010161302d9c",
			expected: Port197Payload{
				Rssi1:   helpers.Int8Ptr(-67),
				Mac1:    "6fbcfdd76434",
				Rssi2:   helpers.Int8Ptr(-69),
				Mac2:    helpers.StringPtr("7e7cbff22fc5"),
				Rssi3:   helpers.Int8Ptr(-71),
				Mac3:    helpers.StringPtr("00dc0af60588"),
				Rssi4:   helpers.Int8Ptr(-73),
				Mac4:    helpers.StringPtr("010161302d9c"),
				Moving:  false,
				Version: Port197Version2,
			},
		},
		{
			port:    197,
			payload: "01b7218f6c166fadb359ea3bdec77daff72faac81784ab263386a455d3a73592a063900b",
			expected: Port197Payload{
				Rssi1:   helpers.Int8Ptr(-73),
				Mac1:    "218f6c166fad",
				Rssi2:   helpers.Int8Ptr(-77),
				Mac2:    helpers.StringPtr("59ea3bdec77d"),
				Rssi3:   helpers.Int8Ptr(-81),
				Mac3:    helpers.StringPtr("f72faac81784"),
				Rssi4:   helpers.Int8Ptr(-85),
				Mac4:    helpers.StringPtr("263386a455d3"),
				Rssi5:   helpers.Int8Ptr(-89),
				Mac5:    helpers.StringPtr("3592a063900b"),
				Moving:  false,
				Version: Port197Version2,
			},
		},
		{
			port:    197,
			payload: "01cff0b0140c96bbcce4c32a622ea4c8e0286d8a9478b8e0286d8aabfcafa86e84e1a812",
			expected: Port197Payload{
				Mac1:    "f0b0140c96bb",
				Rssi1:   helpers.Int8Ptr(-49),
				Mac2:    helpers.StringPtr("e4c32a622ea4"),
				Rssi2:   helpers.Int8Ptr(-52),
				Mac3:    helpers.StringPtr("e0286d8a9478"),
				Rssi3:   helpers.Int8Ptr(-56),
				Mac4:    helpers.StringPtr("e0286d8aabfc"),
				Rssi4:   helpers.Int8Ptr(-72),
				Mac5:    helpers.StringPtr("a86e84e1a812"),
				Rssi5:   helpers.Int8Ptr(-81),
				Moving:  false,
				Version: Port197Version2,
			},
		},
		{
			port:        198,
			payload:     "ff",
			expected:    Port198Payload{},
			expectedErr: "port not supported: version 255 for port 198 not supported",
		},
		{
			port:    198,
			payload: "003385f8ee30c2",
			expected: Port198Payload{
				Rssi1:   nil,
				Mac1:    "3385f8ee30c2",
				Moving:  true,
				Version: Port198Version1,
			},
		},
		{
			port:    198,
			payload: "003385f8ee30c2a0382c2601db",
			expected: Port198Payload{
				Rssi1:   nil,
				Mac1:    "3385f8ee30c2",
				Mac2:    helpers.StringPtr("a0382c2601db"),
				Moving:  true,
				Version: Port198Version1,
			},
		},
		{
			port:    198,
			payload: "00b5eded55a313a0b8b5e86e3194a765f3ad40",
			expected: Port198Payload{
				Rssi1:   nil,
				Mac1:    "b5eded55a313",
				Mac2:    helpers.StringPtr("a0b8b5e86e31"),
				Mac3:    helpers.StringPtr("94a765f3ad40"),
				Moving:  true,
				Version: Port198Version1,
			},
		},
		{
			port:    198,
			payload: "006fbcfdd764347e7cbff22fc500dc0af60588010161302d9c",
			expected: Port198Payload{
				Rssi1:   nil,
				Mac1:    "6fbcfdd76434",
				Mac2:    helpers.StringPtr("7e7cbff22fc5"),
				Mac3:    helpers.StringPtr("00dc0af60588"),
				Mac4:    helpers.StringPtr("010161302d9c"),
				Moving:  true,
				Version: Port198Version1,
			},
		},
		{
			port:    198,
			payload: "00218f6c166fad59ea3bdec77df72faac81784263386a455d33592a063900b",
			expected: Port198Payload{
				Rssi1:   nil,
				Mac1:    "218f6c166fad",
				Mac2:    helpers.StringPtr("59ea3bdec77d"),
				Mac3:    helpers.StringPtr("f72faac81784"),
				Mac4:    helpers.StringPtr("263386a455d3"),
				Mac5:    helpers.StringPtr("3592a063900b"),
				Moving:  true,
				Version: Port198Version1,
			},
		},
		{
			port:    198,
			payload: "01d63385f8ee30c2",
			expected: Port198Payload{
				Rssi1:   helpers.Int8Ptr(-42),
				Mac1:    "3385f8ee30c2",
				Moving:  true,
				Version: Port198Version2,
			},
		},
		{
			port:    198,
			payload: "01d63385f8ee30c2d0a0382c2601db",
			expected: Port198Payload{
				Rssi1:   helpers.Int8Ptr(-42),
				Mac1:    "3385f8ee30c2",
				Rssi2:   helpers.Int8Ptr(-48),
				Mac2:    helpers.StringPtr("a0382c2601db"),
				Moving:  true,
				Version: Port198Version2,
			},
		},
		{
			port:    198,
			payload: "01c8b5eded55a313c0a0b8b5e86e31b894a765f3ad40",
			expected: Port198Payload{
				Rssi1:   helpers.Int8Ptr(-56),
				Mac1:    "b5eded55a313",
				Rssi2:   helpers.Int8Ptr(-64),
				Mac2:    helpers.StringPtr("a0b8b5e86e31"),
				Rssi3:   helpers.Int8Ptr(-72),
				Mac3:    helpers.StringPtr("94a765f3ad40"),
				Moving:  true,
				Version: Port198Version2,
			},
		},
		{
			port:    198,
			payload: "01bd6fbcfdd76434bb7e7cbff22fc5b900dc0af60588b7010161302d9c",
			expected: Port198Payload{
				Rssi1:   helpers.Int8Ptr(-67),
				Mac1:    "6fbcfdd76434",
				Rssi2:   helpers.Int8Ptr(-69),
				Mac2:    helpers.StringPtr("7e7cbff22fc5"),
				Rssi3:   helpers.Int8Ptr(-71),
				Mac3:    helpers.StringPtr("00dc0af60588"),
				Rssi4:   helpers.Int8Ptr(-73),
				Mac4:    helpers.StringPtr("010161302d9c"),
				Moving:  true,
				Version: Port198Version2,
			},
		},
		{
			port:    198,
			payload: "01b7218f6c166fadb359ea3bdec77daff72faac81784ab263386a455d3a73592a063900b",
			expected: Port198Payload{
				Rssi1:   helpers.Int8Ptr(-73),
				Mac1:    "218f6c166fad",
				Rssi2:   helpers.Int8Ptr(-77),
				Mac2:    helpers.StringPtr("59ea3bdec77d"),
				Rssi3:   helpers.Int8Ptr(-81),
				Mac3:    helpers.StringPtr("f72faac81784"),
				Rssi4:   helpers.Int8Ptr(-85),
				Mac4:    helpers.StringPtr("263386a455d3"),
				Rssi5:   helpers.Int8Ptr(-89),
				Mac5:    helpers.StringPtr("3592a063900b"),
				Moving:  true,
				Version: Port198Version2,
			},
		},
		{
			port:        200,
			payload:     "68b9ac21ff",
			expected:    Port200Payload{},
			expectedErr: "port not supported: version 255 for port 200 not supported",
		},
		{
			port:    200,
			payload: "68b9ac21003385f8ee30c2",
			expected: Port200Payload{
				Timestamp: time.Date(2025, 9, 4, 15, 11, 29, 0, time.UTC),
				Rssi1:     nil,
				Mac1:      "3385f8ee30c2",
				Moving:    false,
				Version:   Port200Version1,
			},
		},
		{
			port:    200,
			payload: "68b9ac21003385f8ee30c2a0382c2601db",
			expected: Port200Payload{
				Timestamp: time.Date(2025, 9, 4, 15, 11, 29, 0, time.UTC),
				Rssi1:     nil,
				Mac1:      "3385f8ee30c2",
				Mac2:      helpers.StringPtr("a0382c2601db"),
				Moving:    false,
				Version:   Port200Version1,
			},
		},
		{
			port:    200,
			payload: "68b9ac2100b5eded55a313a0b8b5e86e3194a765f3ad40",
			expected: Port200Payload{
				Timestamp: time.Date(2025, 9, 4, 15, 11, 29, 0, time.UTC),
				Rssi1:     nil,
				Mac1:      "b5eded55a313",
				Mac2:      helpers.StringPtr("a0b8b5e86e31"),
				Mac3:      helpers.StringPtr("94a765f3ad40"),
				Moving:    false,
				Version:   Port200Version1,
			},
		},
		{
			port:    200,
			payload: "68b9ac21006fbcfdd764347e7cbff22fc500dc0af60588010161302d9c",
			expected: Port200Payload{
				Timestamp: time.Date(2025, 9, 4, 15, 11, 29, 0, time.UTC),
				Rssi1:     nil,
				Mac1:      "6fbcfdd76434",
				Mac2:      helpers.StringPtr("7e7cbff22fc5"),
				Mac3:      helpers.StringPtr("00dc0af60588"),
				Mac4:      helpers.StringPtr("010161302d9c"),
				Moving:    false,
				Version:   Port200Version1,
			},
		},
		{
			port:    200,
			payload: "68b9ac2100218f6c166fad59ea3bdec77df72faac81784263386a455d33592a063900b",
			expected: Port200Payload{
				Timestamp: time.Date(2025, 9, 4, 15, 11, 29, 0, time.UTC),
				Rssi1:     nil,
				Mac1:      "218f6c166fad",
				Mac2:      helpers.StringPtr("59ea3bdec77d"),
				Mac3:      helpers.StringPtr("f72faac81784"),
				Mac4:      helpers.StringPtr("263386a455d3"),
				Mac5:      helpers.StringPtr("3592a063900b"),
				Moving:    false,
				Version:   Port200Version1,
			},
		},
		{
			port:    200,
			payload: "68b9ac2101d63385f8ee30c2",
			expected: Port200Payload{
				Timestamp: time.Date(2025, 9, 4, 15, 11, 29, 0, time.UTC),
				Rssi1:     helpers.Int8Ptr(-42),
				Mac1:      "3385f8ee30c2",
				Moving:    false,
				Version:   Port200Version2,
			},
		},
		{
			port:    200,
			payload: "68b9ac2101d63385f8ee30c2d0a0382c2601db",
			expected: Port200Payload{
				Timestamp: time.Date(2025, 9, 4, 15, 11, 29, 0, time.UTC),
				Rssi1:     helpers.Int8Ptr(-42),
				Mac1:      "3385f8ee30c2",
				Rssi2:     helpers.Int8Ptr(-48),
				Mac2:      helpers.StringPtr("a0382c2601db"),
				Moving:    false,
				Version:   Port200Version2,
			},
		},
		{
			port:    200,
			payload: "68b9ac2101c8b5eded55a313c0a0b8b5e86e31b894a765f3ad40",
			expected: Port200Payload{
				Timestamp: time.Date(2025, 9, 4, 15, 11, 29, 0, time.UTC),
				Rssi1:     helpers.Int8Ptr(-56),
				Mac1:      "b5eded55a313",
				Rssi2:     helpers.Int8Ptr(-64),
				Mac2:      helpers.StringPtr("a0b8b5e86e31"),
				Rssi3:     helpers.Int8Ptr(-72),
				Mac3:      helpers.StringPtr("94a765f3ad40"),
				Moving:    false,
				Version:   Port200Version2,
			},
		},
		{
			port:    200,
			payload: "68b9ac2101bd6fbcfdd76434bb7e7cbff22fc5b900dc0af60588b7010161302d9c",
			expected: Port200Payload{
				Timestamp: time.Date(2025, 9, 4, 15, 11, 29, 0, time.UTC),
				Rssi1:     helpers.Int8Ptr(-67),
				Mac1:      "6fbcfdd76434",
				Rssi2:     helpers.Int8Ptr(-69),
				Mac2:      helpers.StringPtr("7e7cbff22fc5"),
				Rssi3:     helpers.Int8Ptr(-71),
				Mac3:      helpers.StringPtr("00dc0af60588"),
				Rssi4:     helpers.Int8Ptr(-73),
				Mac4:      helpers.StringPtr("010161302d9c"),
				Moving:    false,
				Version:   Port200Version2,
			},
		},
		{
			port:    200,
			payload: "68b9ac2101b7218f6c166fadb359ea3bdec77daff72faac81784ab263386a455d3a73592a063900b",
			expected: Port200Payload{
				Timestamp: time.Date(2025, 9, 4, 15, 11, 29, 0, time.UTC),
				Rssi1:     helpers.Int8Ptr(-73),
				Mac1:      "218f6c166fad",
				Rssi2:     helpers.Int8Ptr(-77),
				Mac2:      helpers.StringPtr("59ea3bdec77d"),
				Rssi3:     helpers.Int8Ptr(-81),
				Mac3:      helpers.StringPtr("f72faac81784"),
				Rssi4:     helpers.Int8Ptr(-85),
				Mac4:      helpers.StringPtr("263386a455d3"),
				Rssi5:     helpers.Int8Ptr(-89),
				Mac5:      helpers.StringPtr("3592a063900b"),
				Moving:    false,
				Version:   Port200Version2,
			},
		},
		{
			port:        201,
			payload:     "68b9ac21ff",
			expected:    Port201Payload{},
			expectedErr: "port not supported: version 255 for port 201 not supported",
		},
		{
			port:    201,
			payload: "68bae3ab01d3f0b0140c96bbc7e4c32a622ea4c5e0286d8a9478b4e0286d8aabfcada86e84e1a812",
			expected: Port201Payload{
				Timestamp: time.Date(2025, 9, 5, 13, 20, 43, 0, time.UTC),
				Mac1:      "f0b0140c96bb",
				Rssi1:     helpers.Int8Ptr(-45),
				Mac2:      helpers.StringPtr("e4c32a622ea4"),
				Rssi2:     helpers.Int8Ptr(-57),
				Mac3:      helpers.StringPtr("e0286d8a9478"),
				Rssi3:     helpers.Int8Ptr(-59),
				Mac4:      helpers.StringPtr("e0286d8aabfc"),
				Rssi4:     helpers.Int8Ptr(-76),
				Mac5:      helpers.StringPtr("a86e84e1a812"),
				Rssi5:     helpers.Int8Ptr(-83),
				Moving:    true,
				Version:   Port201Version2,
			},
		},
		{
			port:    201,
			payload: "68b9ac21003385f8ee30c2",
			expected: Port201Payload{
				Timestamp: time.Date(2025, 9, 4, 15, 11, 29, 0, time.UTC),
				Rssi1:     nil,
				Mac1:      "3385f8ee30c2",
				Moving:    true,
				Version:   Port201Version1,
			},
		},
		{
			port:    201,
			payload: "68b9ac21003385f8ee30c2a0382c2601db",
			expected: Port201Payload{
				Timestamp: time.Date(2025, 9, 4, 15, 11, 29, 0, time.UTC),
				Rssi1:     nil,
				Mac1:      "3385f8ee30c2",
				Mac2:      helpers.StringPtr("a0382c2601db"),
				Moving:    true,
				Version:   Port201Version1,
			},
		},
		{
			port:    201,
			payload: "68b9ac2100b5eded55a313a0b8b5e86e3194a765f3ad40",
			expected: Port201Payload{
				Timestamp: time.Date(2025, 9, 4, 15, 11, 29, 0, time.UTC),
				Rssi1:     nil,
				Mac1:      "b5eded55a313",
				Mac2:      helpers.StringPtr("a0b8b5e86e31"),
				Mac3:      helpers.StringPtr("94a765f3ad40"),
				Moving:    true,
				Version:   Port201Version1,
			},
		},
		{
			port:    201,
			payload: "68b9ac21006fbcfdd764347e7cbff22fc500dc0af60588010161302d9c",
			expected: Port201Payload{
				Timestamp: time.Date(2025, 9, 4, 15, 11, 29, 0, time.UTC),
				Rssi1:     nil,
				Mac1:      "6fbcfdd76434",
				Mac2:      helpers.StringPtr("7e7cbff22fc5"),
				Mac3:      helpers.StringPtr("00dc0af60588"),
				Mac4:      helpers.StringPtr("010161302d9c"),
				Moving:    true,
				Version:   Port201Version1,
			},
		},
		{
			port:    201,
			payload: "68b9ac2100218f6c166fad59ea3bdec77df72faac81784263386a455d33592a063900b",
			expected: Port201Payload{
				Timestamp: time.Date(2025, 9, 4, 15, 11, 29, 0, time.UTC),
				Rssi1:     nil,
				Mac1:      "218f6c166fad",
				Mac2:      helpers.StringPtr("59ea3bdec77d"),
				Mac3:      helpers.StringPtr("f72faac81784"),
				Mac4:      helpers.StringPtr("263386a455d3"),
				Mac5:      helpers.StringPtr("3592a063900b"),
				Moving:    true,
				Version:   Port201Version1,
			},
		},
		{
			port:    201,
			payload: "68b9ac2101d63385f8ee30c2",
			expected: Port201Payload{
				Timestamp: time.Date(2025, 9, 4, 15, 11, 29, 0, time.UTC),
				Rssi1:     helpers.Int8Ptr(-42),
				Mac1:      "3385f8ee30c2",
				Moving:    true,
				Version:   Port201Version2,
			},
		},
		{
			port:    201,
			payload: "68b9ac2101d63385f8ee30c2d0a0382c2601db",
			expected: Port201Payload{
				Timestamp: time.Date(2025, 9, 4, 15, 11, 29, 0, time.UTC),
				Rssi1:     helpers.Int8Ptr(-42),
				Mac1:      "3385f8ee30c2",
				Rssi2:     helpers.Int8Ptr(-48),
				Mac2:      helpers.StringPtr("a0382c2601db"),
				Moving:    true,
				Version:   Port201Version2,
			},
		},
		{
			port:    201,
			payload: "68b9ac2101c8b5eded55a313c0a0b8b5e86e31b894a765f3ad40",
			expected: Port201Payload{
				Timestamp: time.Date(2025, 9, 4, 15, 11, 29, 0, time.UTC),
				Rssi1:     helpers.Int8Ptr(-56),
				Mac1:      "b5eded55a313",
				Rssi2:     helpers.Int8Ptr(-64),
				Mac2:      helpers.StringPtr("a0b8b5e86e31"),
				Rssi3:     helpers.Int8Ptr(-72),
				Mac3:      helpers.StringPtr("94a765f3ad40"),
				Moving:    true,
				Version:   Port201Version2,
			},
		},
		{
			port:    201,
			payload: "68b9ac2101bd6fbcfdd76434bb7e7cbff22fc5b900dc0af60588b7010161302d9c",
			expected: Port201Payload{
				Timestamp: time.Date(2025, 9, 4, 15, 11, 29, 0, time.UTC),
				Rssi1:     helpers.Int8Ptr(-67),
				Mac1:      "6fbcfdd76434",
				Rssi2:     helpers.Int8Ptr(-69),
				Mac2:      helpers.StringPtr("7e7cbff22fc5"),
				Rssi3:     helpers.Int8Ptr(-71),
				Mac3:      helpers.StringPtr("00dc0af60588"),
				Rssi4:     helpers.Int8Ptr(-73),
				Mac4:      helpers.StringPtr("010161302d9c"),
				Moving:    true,
				Version:   Port201Version2,
			},
		},
		{
			port:    201,
			payload: "68b9ac2101b7218f6c166fadb359ea3bdec77daff72faac81784ab263386a455d3a73592a063900b",
			expected: Port201Payload{
				Timestamp: time.Date(2025, 9, 4, 15, 11, 29, 0, time.UTC),
				Rssi1:     helpers.Int8Ptr(-73),
				Mac1:      "218f6c166fad",
				Rssi2:     helpers.Int8Ptr(-77),
				Mac2:      helpers.StringPtr("59ea3bdec77d"),
				Rssi3:     helpers.Int8Ptr(-81),
				Mac3:      helpers.StringPtr("f72faac81784"),
				Rssi4:     helpers.Int8Ptr(-85),
				Mac4:      helpers.StringPtr("263386a455d3"),
				Rssi5:     helpers.Int8Ptr(-89),
				Mac5:      helpers.StringPtr("3592a063900b"),
				Moving:    true,
				Version:   Port201Version2,
			},
		},
	}

	if logger.Logger == nil {
		logger.NewLogger()
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("TestPort%vWith%v", test.port, test.payload), func(t *testing.T) {
			ctx := context.WithValue(context.Background(), decoder.DEVEUI_CONTEXT_KEY, test.devEui)
			ctx = context.WithValue(ctx, decoder.FCNT_CONTEXT_KEY, 1)

			// Use SolverV2 for GNSS ports (192/193/194/195/199) so timestamped ports work without error and provide expected data.
			expectedAny := test.expected
			opts := []Option{}
			switch test.port {
			case 192, 194, 195:
				// For GNSS ports, use SolverV2 and return the same structure as port 192 expectation
				// so that tests compare against exampleResponse. For 194/195, timestamp is handled by decoder.
				v2Data := &exampleResponse
				features := []decoder.Feature{decoder.FeatureGNSS}
				if expectedAny == nil {
					expectedAny = v2Data
				}
				opts = append(opts, WithSolverV2(solver.MockSolverV2{
					Data: decoder.NewDecodedUplink(features, v2Data),
				}))
			}
			dec := NewTagXLv1Decoder(ctx, solver.MockSolverV1{}, logger.Logger, opts...)
			got, err := dec.Decode(ctx, test.payload, test.port)

			if err == nil && len(test.expectedErr) != 0 {
				t.Fatalf("expected error: %v, got %v", test.expectedErr, nil)
			}

			if err != nil && len(test.expectedErr) == 0 {
				t.Fatalf("unexpected error: %v", err)
			}

			t.Logf("got %v", got)

			if got != nil && !reflect.DeepEqual(got.Data, expectedAny) && len(test.expectedErr) == 0 {
				// marshal the expected and got values to compare
				expectedJSON, err := json.Marshal(expectedAny)
				if err != nil {
					t.Fatalf("failed to marshal expected value: %v", err)
				}
				gotJSON, err := json.Marshal(got.Data)
				if err != nil {
					t.Fatalf("failed to marshal got value: %v", err)
				}
				t.Errorf("expected: %s, got: %s", expectedJSON, gotJSON)
			}

			if len(test.expectedErr) > 0 && err != nil && !strings.Contains(err.Error(), test.expectedErr) {
				t.Errorf("expected error: %v, got: %v", test.expectedErr, err)
			}
		})
	}
}

func TestDecodeWithNoopSolver(t *testing.T) {
	tests := []struct {
		payload string
		port    uint8
		devEui  string
	}{
		{
			port:    192,
			payload: "87821f50490200b520fbe977844d222a3a14a89293956245cc75a9ca1bbc25ddf658542909",
			devEui:  "10CE45FFFE00C7EC",
		},
		{
			port:    192,
			payload: "87821f50490200b520fbe977844d222a3a14a89293956245cc75a9ca1bbc25ddf658542909",
			devEui:  "10CE45FFFE00C7ED",
		},
		{
			port:    199,
			payload: "86b5277140484a89b8f63ccf67affbfeb519b854f9d447808a50785bdfe86a77",
			devEui:  "10CE45FFFE00C7EC",
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestPort%vWith%v", test.port, test.payload), func(t *testing.T) {
			decoder := NewTagXLv1Decoder(context.TODO(), solver.NoopSolver{}, zap.NewNop())
			got, err := decoder.Decode(context.TODO(), test.payload, test.port)

			assert.NotNil(t, got)
			assert.NoError(t, err)
		})
	}
}

func TestValidationErrors(t *testing.T) {
	tests := []struct {
		payload  string
		port     uint8
		expected error
	}{}

	if logger.Logger == nil {
		logger.NewLogger()
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestPort%vValidationWith%v", test.port, test.payload), func(t *testing.T) {
			ctx := context.WithValue(context.Background(), decoder.DEVEUI_CONTEXT_KEY, "10CE45FFFE00C7ED")
			ctx = context.WithValue(ctx, decoder.FCNT_CONTEXT_KEY, 1)

			decoder := NewTagXLv1Decoder(ctx, solver.MockSolverV1{}, logger.Logger)
			got, err := decoder.Decode(ctx, test.payload, test.port)

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
	if logger.Logger == nil {
		logger.NewLogger()
	}

	decoder := NewTagXLv1Decoder(context.TODO(), solver.MockSolverV1{}, logger.Logger)
	_, err := decoder.Decode(context.TODO(), "00", 0)

	if err == nil || !errors.Is(err, helpers.ErrPortNotSupported) {
		t.Fatal("expected port not supported")
	}
}

func TestPayloadTooShort(t *testing.T) {
	if logger.Logger == nil {
		logger.NewLogger()
	}

	decoder := NewTagXLv1Decoder(context.TODO(), solver.MockSolverV1{}, logger.Logger)
	_, err := decoder.Decode(context.TODO(), "01adbeef", 152)

	if err == nil || !errors.Is(err, helpers.ErrPayloadTooShort) {
		t.Fatalf("expected error payload too short but got %v", err)
	}
}

func TestPayloadTooLong(t *testing.T) {
	if logger.Logger == nil {
		logger.NewLogger()
	}

	decoder := NewTagXLv1Decoder(context.TODO(), solver.MockSolverV1{}, logger.Logger)
	_, err := decoder.Decode(context.TODO(), "01adbeef4242deadbeef4242deadbeef4242", 152)

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
			port:    151,
			payload: "4c03014104003c0078",
		},
		{
			port:    151,
			payload: "4c03014204000a03e8",
		},
		{
			port:    151,
			payload: "4c2a0940010f4104012c1c204204012c05dc43010644011e45020d4e4604f6c7d8104902000a4a0400000002",
		},
		{
			port:    151,
			payload: "4c2d0a40010b410402581c204204012c05dc43010644011e45020d6c4604a25b545547010249020003",
		},
		{
			payload: "010b0066acbcf0000000000756",
			port:    152,
		},
		{
			payload: "010b2266acbcf0000000000756",
			port:    152,
		},
		{
			payload: "020c62116822f120000d00000024",
			port:    152,
		},
		{
			payload: "020c09336823166a000000000109",
			port:    152,
		},
		{
			payload: "87821f50490200b520fbe977844d222a3a14a89293956245cc75a9ca1bbc25ddf658542909",
			port:    192,
		},
		{
			payload: "00218f6c166fad59ea3bdec77df72faac81784263386a455d33592a063900b",
			port:    197,
		},
		{
			payload: "01b7218f6c166fadb359ea3bdec77daff72faac81784ab263386a455d3a73592a063900b",
			port:    197,
		},
		{
			payload: "00218f6c166fad59ea3bdec77df72faac81784263386a455d33592a063900b",
			port:    198,
		},
		{
			payload: "01b7218f6c166fadb359ea3bdec77daff72faac81784ab263386a455d3a73592a063900b",
			port:    198,
		},
		{
			payload:         "86b5277140484a89b8f63ccf67affbfeb519b854f9d447808a50785bdfe86a77",
			port:            199,
			allowNoFeatures: true,
		},
		{
			payload: "68b9ac2100218f6c166fad59ea3bdec77df72faac81784263386a455d33592a063900b",
			port:    200,
		},
		{
			payload: "68b9ac2101b7218f6c166fadb359ea3bdec77daff72faac81784ab263386a455d3a73592a063900b",
			port:    200,
		},
		{
			payload: "68b9ac2100218f6c166fad59ea3bdec77df72faac81784263386a455d33592a063900b",
			port:    201,
		},
		{
			payload: "68b9ac2101b7218f6c166fadb359ea3bdec77daff72faac81784ab263386a455d3a73592a063900b",
			port:    201,
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
	middleware, err := loracloud.NewLoracloudClient(context.TODO(), "access_token", zap.NewExample())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	middleware.BaseUrl = server.URL
	defer server.Close()

	if logger.Logger == nil {
		logger.NewLogger()
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestFeaturesWithPort%vAndPayload%v", test.port, test.payload), func(t *testing.T) {
			d := NewTagXLv1Decoder(context.TODO(), solver.MockSolverV1{
				Data: decoder.NewDecodedUplink([]decoder.Feature{decoder.FeatureWiFi}, Port197Payload{}),
			}, logger.Logger)
			decodedPayload, err := d.Decode(context.TODO(), test.payload, test.port)
			if err != nil {
				t.Fatalf("error %s", err)
			}

			if len(decodedPayload.GetFeatures()) == 0 && !test.allowNoFeatures {
				t.Error("expected features, got none")
			}

			if decodedPayload.Is(decoder.FeatureTimestamp) {
				timestamp, ok := decodedPayload.Data.(decoder.UplinkFeatureTimestamp)
				if !ok {
					t.Fatalf("expected UplinkFeatureTimestamp, got %T", decodedPayload)
				}
				if timestamp.GetTimestamp() == nil {
					t.Fatalf("expected non nil timestamp")
				}
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
				buffered.IsBuffered()
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
			if decodedPayload.Is(decoder.FeatureFirmwareVersion) {
				firmwareVersion, ok := decodedPayload.Data.(decoder.UplinkFeatureFirmwareVersion)
				if !ok {
					t.Fatalf("expected UplinkFeatureFirmwareVersion, got %T", decodedPayload)
				}
				firmwareVersion.GetFirmwareVersion()
				if firmwareVersion.GetFirmwareHash() == nil {
					t.Fatalf("expected non nil firmware hash")
				}
			}
			if decodedPayload.Is(decoder.FeatureRotationState) {
				rotationState, ok := decodedPayload.Data.(decoder.UplinkFeatureRotationState)
				if !ok {
					t.Fatalf("expected UplinkFeatureRotationState, got %T", decodedPayload)
				}
				// call function to check if it panics
				rotationState.GetOldRotationState()
				rotationState.GetNewRotationState()
				rotationState.GetRotations()
				rotationState.GetDuration()
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
			expected: []string{"\"oldRotationState\": \"undefined\"", "\"newRotationState\": \"pouring\"", "\"timestamp\": \"2024-08-02T11:03:12Z\"", "\"elapsedSeconds\": 1878"},
		},
	}

	if logger.Logger == nil {
		logger.NewLogger()
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestMarshalWithPort%vAndPayload%v", test.port, test.payload), func(t *testing.T) {
			decoder := NewTagXLv1Decoder(context.TODO(), solver.MockSolverV1{}, logger.Logger)

			data, _ := decoder.Decode(context.TODO(), test.payload, test.port)

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

func TestNewTagXLv1DecoderWithNilSolver(t *testing.T) {
	assert.Panics(t, func() {
		NewTagXLv1Decoder(context.TODO(), nil, zap.NewExample())
	}, "NewTagXLv1Decoder should panic when solver is nil")
}

func TestNewTagXLv1DecoderSolver(t *testing.T) {
	tests := []struct {
		port        uint8
		payload     string
		expected    any
		expectedErr string
	}{
		{
			port:        192,
			payload:     "deadbeef",
			expected:    nil,
			expectedErr: "solver failed",
		},
	}
	solver, err := aws.NewAwsPositionEstimateClient(context.TODO(), zap.NewExample())
	if err != nil {
		t.Fatalf("failed to create aws solver %v", err)
	}
	decoder := NewTagXLv1Decoder(context.TODO(), solver, zap.NewExample())
	for _, test := range tests {
		received, err := decoder.Decode(context.TODO(), test.payload, test.port)
		if received != test.expected && received != nil && test.expected != nil {
			t.Errorf("expected %v", test.expected)
			t.Errorf("received %v", received)
		}
		if test.expectedErr != "" && (err == nil || !strings.Contains(err.Error(), test.expectedErr)) {
			t.Errorf("expected %v", test.expectedErr)
			t.Errorf("received %v", err)
		}
	}
}
