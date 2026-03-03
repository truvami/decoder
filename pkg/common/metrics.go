package common

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	unknownTLVTagsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "truvami_common_unknown_tlv_tags_total",
		Help: "The total number of unknown TLV tags encountered during decoding",
	}, []string{"tag"})
)
