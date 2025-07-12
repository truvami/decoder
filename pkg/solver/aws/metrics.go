package aws

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	awsPositionEstimatesTotalCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "truvami_aws_position_estimates_total",
		Help: "The total number of processed position estimate requests",
	})
	awsPositionEstimatesErrorsCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "truvami_aws_position_estimates_errors_total",
		Help: "The total number of errors encountered while processing position estimate requests",
	})
	awsPositionEstimatesDurationHistogram = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "truvami_aws_position_estimates_duration_seconds",
		Help:    "The duration of position estimate requests in seconds",
		Buckets: []float64{0.1, 0.2, 0.3, 0.5, 1, 2, 5, 10, 30, 60},
	})
	awsPositionEstimatesSuccessCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "truvami_aws_position_estimates_success_total",
		Help: "The total number of successful position estimate requests",
	})
	awsPositionEstimatesFailureCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "truvami_aws_position_estimates_failure_total",
		Help: "The total number of failed position estimate requests",
	})

	AwsLoracloudFallbackSuccess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "truvami_aws_loracloud_fallback_success_total",
		Help: "The total number of successful position estimate requests using Loracloud as a fallback",
	})
	AwsLoracloudFallbackFailure = promauto.NewCounter(prometheus.CounterOpts{
		Name: "truvami_aws_loracloud_fallback_failure_total",
		Help: "The total number of failed position estimate requests using Loracloud as a fallback",
	})
)
