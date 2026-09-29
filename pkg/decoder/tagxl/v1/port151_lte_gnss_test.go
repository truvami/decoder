package tagxl

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/truvami/decoder/internal/logger"
	"github.com/truvami/decoder/pkg/common"
	"github.com/truvami/decoder/pkg/decoder"
	"github.com/truvami/decoder/pkg/solver"
)

func decodePort151(t *testing.T, payload string) (*decoder.DecodedUplink, Port151Payload) {
	t.Helper()
	dec := NewTagXLv1Decoder(context.TODO(), solver.MockSolverV1{}, logger.Logger)
	decoded, err := dec.Decode(context.TODO(), payload, 151)
	require.NoError(t, err)
	data, ok := decoded.Data.(Port151Payload)
	require.True(t, ok, "unexpected payload type %T", decoded.Data)
	return decoded, data
}

func TestPort151_GpsSettings(t *testing.T) {
	t.Parallel()

	// settings report after a C2 + C1 downlink: GPS+Galileo, low power; scan 120 s, cool-off 10 s, PDOP <= 3.0
	decoded, data := decodePort151(t, "4c1002d2020201d10a0078000a00000000001e")

	assert.True(t, decoded.Is(decoder.FeatureConfig))
	assert.Equal(t, common.Uint8Ptr(2), data.GpsConstellation)
	assert.Equal(t, common.Uint8Ptr(1), data.GpsPowerMode)
	assert.Equal(t, common.Uint16Ptr(120), data.GpsScanDuration)
	assert.Equal(t, common.Uint16Ptr(10), data.GpsCooloffDuration)
	assert.Equal(t, common.Float32Ptr(0), data.GpsHdopThreshold)
	assert.Equal(t, common.Float32Ptr(0), data.GpsVdopThreshold)
	assert.Equal(t, common.Float32Ptr(3), data.GpsPdopThreshold)
}

func TestPort151_LteSettings(t *testing.T) {
	t.Parallel()

	// APN "iot.1nce.net", API key "abc", upload 02:00 UTC, 100 records, 180 s timeout, 30 min back-off, 3 retries
	decoded, data := decodePort151(t, "4c2203b10c696f742e316e63652e6e6574b203616263b3090078006400b4001e03")

	assert.True(t, decoded.Is(decoder.FeatureConfig))
	assert.Equal(t, common.StringPtr("iot.1nce.net"), data.LteApn)
	assert.Equal(t, common.BoolPtr(true), data.LteApiKeyConfigured)
	assert.Equal(t, common.Uint16Ptr(120), data.LteUploadTimeOfDay)
	assert.Equal(t, common.Uint16Ptr(100), data.LteUploadChunkSize)
	assert.Equal(t, common.Uint16Ptr(180), data.LteConnectionTimeout)
	assert.Equal(t, common.Uint16Ptr(30), data.LteUploadBackoff)
	assert.Equal(t, common.Uint8Ptr(3), data.LteUploadRetryCount)

	// the api key itself must never be exposed
	encoded, err := json.Marshal(data)
	require.NoError(t, err)
	assert.False(t, strings.Contains(string(encoded), "abc"))
}

func TestPort151_EmptyLteApiKey(t *testing.T) {
	t.Parallel()

	_, data := decodePort151(t, "4c0301b200")
	assert.Equal(t, common.BoolPtr(false), data.LteApiKeyConfigured)
}

func TestPort151_TruncatedLteGpsValues(t *testing.T) {
	t.Parallel()

	// truncated 0xD1 and 0xB3 only fill the fields that are fully present
	_, data := decodePort151(t, "4c0802d1020078b30300780a")

	assert.Equal(t, common.Uint16Ptr(120), data.GpsScanDuration)
	assert.Nil(t, data.GpsCooloffDuration)
	assert.Nil(t, data.GpsPdopThreshold)
	assert.Equal(t, common.Uint16Ptr(120), data.LteUploadTimeOfDay)
	assert.Nil(t, data.LteUploadChunkSize)
	assert.Nil(t, data.LteUploadRetryCount)
}

func TestPort151_LteGpsNullWhenNotReported(t *testing.T) {
	t.Parallel()

	// a Tag XL report without the new tags: fields stay nil and are encoded as null,
	// and 0x46 is still the firmware hash
	_, data := decodePort151(t, "4c0d0245020a924604aabbccdd")

	assert.Equal(t, common.StringPtr("aabbccdd"), data.FirmwareHash)
	assert.Nil(t, data.GetFirmwareVersion())
	assert.Nil(t, data.LteApn)
	assert.Nil(t, data.GpsConstellation)

	encoded, err := json.Marshal(data)
	require.NoError(t, err)
	assert.Contains(t, string(encoded), `"gpsConstellation":null`)
	assert.Contains(t, string(encoded), `"lteApn":null`)
}
