package metric

import (
	logger "github.com/gotech-hub/go-logging"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	meter "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
)

var (
	meterGlobal = otel.Meter("instrumentation")
)

func init() {
	log := logger.GetLogger()
	exporter, err := prometheus.New()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize prometheus exporter")
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(exporter),
		metric.WithView(metric.NewView(
			metric.Instrument{Name: "http.server.duration"},
			metric.Stream{Aggregation: metric.AggregationExplicitBucketHistogram{
				Boundaries: []float64{0.01, 0.05, 0.1, 0.5, 1, 2, 5},
			}},
		)),
	)

	meterGlobal = meterProvider.Meter("instrumentation")
	otel.SetMeterProvider(meterProvider)
}

func GetGlobalMeter() meter.Meter {
	return meterGlobal
}
