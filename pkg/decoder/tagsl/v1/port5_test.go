package tagsl

import (
	"net"
	"strings"
	"testing"

	"github.com/truvami/decoder/pkg/common"
)

func TestGetAccessPoints(t *testing.T) {
	payload := &Port5Payload{
		Mac1:  "001122334455",
		Rssi1: -50,
		Mac2:  "66778899aabb",
		Rssi2: -60,
		Mac3:  "ccddeeff0011",
		Rssi3: -70,
		Mac4:  "223344556677",
		Rssi4: -80,
		Mac5:  "8899aabbccdd",
		Rssi5: -90,
		Mac6:  "eeff00112233",
		Rssi6: -100,
		Mac7:  "445566778899",
		Rssi7: -110,
	}

	expected := []common.WifiAccessPoint{
		{MacAddress: net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}, Rssi: -50},
		{MacAddress: net.HardwareAddr{0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB}, Rssi: -60},
		{MacAddress: net.HardwareAddr{0xCC, 0xDD, 0xEE, 0xFF, 0x00, 0x11}, Rssi: -70},
		{MacAddress: net.HardwareAddr{0x22, 0x33, 0x44, 0x55, 0x66, 0x77}, Rssi: -80},
		{MacAddress: net.HardwareAddr{0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD}, Rssi: -90},
		{MacAddress: net.HardwareAddr{0xEE, 0xFF, 0x00, 0x11, 0x22, 0x33}, Rssi: -100},
		{MacAddress: net.HardwareAddr{0x44, 0x55, 0x66, 0x77, 0x88, 0x99}, Rssi: -110},
	}

	accessPoints := payload.GetAccessPoints()
	if len(accessPoints) != len(expected) {
		t.Fatalf("expected %d access points, got %d", len(expected), len(accessPoints))
	}

	for i, ap := range accessPoints {
		if !strings.EqualFold(ap.MacAddress.String(), expected[i].MacAddress.String()) || ap.Rssi != expected[i].Rssi {
			t.Errorf("expected access point %v, got %v", expected[i], ap)
		}
	}
}
