package tagxl

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	tagXlDecoderSolverFailedCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "truvami_tagxl_v1_decoder_solver_failed_total",
		Help: "The total number of tag XL decodes that failed to solve",
	})
	tagXlDecoderSuccessfullyUsedFallbackSolverCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "truvami_tagxl_v1_decoder_successfully_used_fallback_solver_total",
		Help: "The total number of tag XL decodes that successfully used a fallback solver",
	})
)
