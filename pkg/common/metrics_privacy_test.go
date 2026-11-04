package common

import (
	"strings"
	"testing"

	dto "github.com/prometheus/client_model/go"
)

func TestUnknownTLVMetricHasNoSensitiveLabels(t *testing.T) {
	metric := &dto.Metric{}
	if err := unknownTLVTagsTotal.Write(metric); err != nil {
		t.Fatalf("write metric: %v", err)
	}

	for _, label := range metric.GetLabel() {
		name := label.GetName()
		if strings.Contains(strings.ToLower(name), "tag") ||
			strings.Contains(strings.ToLower(name), "payload") ||
			strings.Contains(strings.ToLower(name), "deveui") {
			t.Fatalf("metric must not expose sensitive label %q", name)
		}
	}
}
