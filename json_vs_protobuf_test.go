package main

// Result:
// BenchmarkJSONMarshalling              50          25473968 ns/op
// BenchmarkProtobufMarshalling         100          18804006 ns/op

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/golang/protobuf/proto"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/model"
)

func generateSamples(count int) *model.Samples {
	samples := model.Samples{}
	now := model.Now()
	for i := 1; i <= count; i++ {
		samples = append(samples,
			&model.Sample{
				Metric: model.Metric{
					"__name__": model.LabelValue(fmt.Sprintf("my_metric_name_%d", i)),
					"foo":      "bar",
					"instance": "myinstance"},
				Timestamp: now,
				Value:     model.SampleValue(1.003452 + float64(i)),
			})
	}
	return &samples
}

func MarshalJSON(samples *model.Samples) {
	type Sample struct {
		Labels    model.Metric
		Timestamp int64
		Value     float64
	}
	type Samples []*Sample
	outSamples := make(Samples, 0, len(*samples))
	for _, s := range *samples {
		outSamples = append(outSamples,
			&Sample{
				Labels:    s.Metric,
				Timestamp: int64(s.Timestamp),
				Value:     float64(s.Value),
			})
	}
	_, err := json.Marshal(outSamples)
	if err != nil {
		panic(err)
	}
}

func BenchmarkJSONMarshalling(b *testing.B) {
	samples := generateSamples(10000)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		MarshalJSON(samples)
	}
}

func labelsFromMetric(m model.Metric) []*dto.LabelPair {
	labels := make([]*dto.LabelPair, 0, len(m)-1)

	for l, v := range m {
		if l == model.MetricNameLabel {
			continue
		}

		labels = append(labels, &dto.LabelPair{
			Name:  proto.String(string(l)),
			Value: proto.String(string(v)),
		})
	}

	return labels
}

func MarshalProtobuf(samples *model.Samples) {
	for _, s := range *samples {
		labels := labelsFromMetric(s.Metric)
		_, err := proto.Marshal(&dto.Metric{
			Label: labels,
			Untyped: &dto.Untyped{
				Value: proto.Float64(float64(s.Value)),
			},
			TimestampMs: proto.Int64(int64(s.Timestamp)),
		})
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkProtobufMarshalling(b *testing.B) {
	samples := generateSamples(10000)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		MarshalProtobuf(samples)
	}
}
