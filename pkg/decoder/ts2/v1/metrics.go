package ts2

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ts2DecoderSolverFailedCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "truvami_ts2_v1_decoder_solver_failed_total",
		Help: "The total number of TS2 decodes that failed to solve",
	})
	ts2DecoderSuccessfullyUsedFallbackSolverCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "truvami_ts2_v1_decoder_successfully_used_fallback_solver_total",
		Help: "The total number of TS2 decodes that successfully used a fallback solver",
	})
)
