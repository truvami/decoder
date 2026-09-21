package solver_test

import (
	"strings"
	"testing"

	dto "github.com/prometheus/client_model/go"
	loracloud "github.com/truvami/decoder/pkg/solver/loracloud"
	loracloudv2 "github.com/truvami/decoder/pkg/solver/loracloud/v2"
)

func TestLoracloudMetricsHaveNoSensitiveLabels(t *testing.T) {
	unlabeled := []interface {
		Write(*dto.Metric) error
	}{
		loracloud.PositionEstimateNoCapturedAtSetCounter,
		loracloud.PositionEstimateZeroCoordinatesSetCounter,
		loracloud.PositionEstimateNoCapturedAtSetWithValidCoordinatesCounter,
		loracloud.PositionEstimateValidCounter,
		loracloud.PositionEstimateInvalidCounter,
		loracloudv2.ResponseInvalidTotal,
		loracloudv2.PositionInvalidTotal,
		loracloudv2.RequestDurationSeconds,
	}

	for _, collector := range unlabeled {
		metric := &dto.Metric{}
		if err := collector.Write(metric); err != nil {
			t.Fatalf("write metric: %v", err)
		}
		assertNoSensitiveLabels(t, metric.GetLabel())
	}

	labeled := []struct {
		write func(*dto.Metric) error
	}{
		{func(m *dto.Metric) error { return loracloudv2.RequestsTotal.WithLabelValues("success").Write(m) }},
		{func(m *dto.Metric) error { return loracloudv2.ErrorsTotal.WithLabelValues("request_failed").Write(m) }},
		{func(m *dto.Metric) error { return loracloudv2.BufferedDetectedTotal.WithLabelValues("5m0s").Write(m) }},
	}

	for _, tc := range labeled {
		metric := &dto.Metric{}
		if err := tc.write(metric); err != nil {
			t.Fatalf("write metric: %v", err)
		}
		assertNoSensitiveLabels(t, metric.GetLabel())
	}
}

func assertNoSensitiveLabels(t *testing.T, labels []*dto.LabelPair) {
	t.Helper()
	for _, label := range labels {
		name := strings.ToLower(label.GetName())
		switch name {
		case "deveui", "base_url", "tag", "payload", "url":
			t.Fatalf("metric must not expose sensitive label %q", label.GetName())
		}
	}
}
