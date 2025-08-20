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
	loracloudPositionEstimateNoCapturedAtSetWithValidCoordinatesCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "truvami_loracloud_position_estimate_no_captured_at_set_with_valid_coordinates_total",
		Help: "The total number of position estimate responses where the captured at (UTC) timestamp is not set and the coordinates are valid",
	}, []string{"devEUI"})
	loracloudPositionEstimateValidCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "truvami_loracloud_position_estimate_valid_total",
		Help: "The total number of position estimate responses where the captured at (UTC) timestamp is set and the coordinates are valid",
	}, []string{"devEUI"})
	loracloudPositionEstimateInvalidCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "truvami_loracloud_position_estimate_invalid_total",
		Help: "The total number of position estimate responses where the position resolution is invalid",
	}, []string{"devEUI"})
)
