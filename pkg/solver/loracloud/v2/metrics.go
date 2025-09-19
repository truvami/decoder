package v2

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	loracloudV2RequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "truvami_loracloud_v2_requests_total",
		Help: "Total number of LoRaCloud v2 solver requests",
	}, []string{"base_url", "outcome"}) // outcome: success|error

	loracloudV2RequestDurationSeconds = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "truvami_loracloud_v2_request_duration_seconds",
		Help:    "Duration of LoRaCloud v2 solver requests in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"base_url"})

	loracloudV2ResponseInvalidTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "truvami_loracloud_v2_response_invalid_total",
		Help: "Total number of invalid responses from LoRaCloud v2",
	}, []string{"base_url"})

	loracloudV2PositionInvalidTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "truvami_loracloud_v2_position_invalid_total",
		Help: "Total number of invalid position resolutions (missing timestamp or zero coordinates)",
	}, []string{"devEui"})

	loracloudV2BufferedDetectedTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "truvami_loracloud_v2_timestamp_buffered_detected_total",
		Help: "Total number of uplinks considered buffered due to past timestamp",
	}, []string{"devEui", "threshold"})

	loracloudV2ErrorsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "truvami_loracloud_v2_errors_total",
		Help: "Total number of errors in LoRaCloud v2 solver",
	}, []string{"base_url", "type"}) // type: build_request|request_failed|unexpected_status|decode_failed|response_invalid|position_invalid|invalid_options
)
