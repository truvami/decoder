package loracloud

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	loracloudPositionEstimateNoCapturedAtSetCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "truvami_loracloud_position_estimate_no_captured_at_set_total",
		Help: "The total number of position estimate responses where the captured at (UTC) timestamp is not set",
	}, []string{"devEUI"})
	loracloudPositionEstimateZeroCoordinatesSetCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "truvami_loracloud_position_estimate_zero_coordinates_set_total",
		Help: "The total number of position estimate responses where the coordinates are set to 0",
	}, []string{"devEUI"})
)
