package v2

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "truvami_loracloud_v2_requests_total",
		Help: "Total number of LoRaCloud v2 solver requests",
	}, []string{"outcome"}) // outcome: success|error

	RequestDurationSeconds = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "truvami_loracloud_v2_request_duration_seconds",
		Help:    "Duration of LoRaCloud v2 solver requests in seconds",
		Buckets: prometheus.DefBuckets,
	})

	ResponseInvalidTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "truvami_loracloud_v2_response_invalid_total",
		Help: "Total number of invalid responses from LoRaCloud v2",
	})

	PositionInvalidTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "truvami_loracloud_v2_position_invalid_total",
		Help: "Total number of invalid position resolutions (missing timestamp or zero coordinates)",
	})

	BufferedDetectedTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "truvami_loracloud_v2_timestamp_buffered_detected_total",
		Help: "Total number of uplinks considered buffered due to past timestamp",
	}, []string{"threshold"})

	ErrorsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "truvami_loracloud_v2_errors_total",
		Help: "Total number of errors in LoRaCloud v2 solver",
	}, []string{"type"}) // type: build_request|request_failed|unexpected_status|decode_failed|response_invalid|position_invalid|invalid_options
)
