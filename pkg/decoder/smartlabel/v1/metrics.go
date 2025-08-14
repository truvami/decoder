package smartlabel

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	smartLabelDecoderSolverFailedCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "truvami_smartlabel_v1_decoder_solver_failed_total",
		Help: "The total number of smart label decodes that failed to solve",
	})
	smartLabelDecoderSuccessfullyUsedFallbackSolverCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "truvami_smartlabel_v1_decoder_successfully_used_fallback_solver_total",
		Help: "The total number of smart label decodes that successfully used a fallback solver",
	})
)
