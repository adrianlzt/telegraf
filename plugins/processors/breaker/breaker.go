package printer

import (
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/processors"
)

type Breaker struct {
	active       bool
	Name         string
	Tag          tag
	Field        string
	ValueEnable  interface{}
	ValueDisable interface{}
}

type tag []struct {
	Key   string
	Value string
}

var sampleConfig = `
`

func (b *Breaker) SampleConfig() string {
	return sampleConfig
}

func (b *Breaker) Description() string {
	return "Print all metrics that pass through this filter."
}

func (b *Breaker) Apply(in ...telegraf.Metric) []telegraf.Metric {
	acceptedMetrics := []telegraf.Metric{}
L1:
	for _, metric := range in {
		// Capture metric used as the breaker setter
		if metric.Name() == b.Name {
			// Check if defined tags are set in the metric
			for _, tag := range b.Tag {
				tagValue, exists := metric.GetTag(tag.Key)
				if !exists || tagValue != tag.Value {
					// this metric is not the one we are looking for
					continue L1
				}
			}

			value, exists := metric.GetField(b.Field)
			if !exists {
				// this metric is not the one we are looking for
				continue L1
			}

			if value == b.ValueEnable {
				b.active = true
			} else if value == b.ValueDisable {
				b.active = false
			}
		}

		// Ignore metrics if breaker is active
		if !b.active {
			acceptedMetrics = append(acceptedMetrics, metric)
		}
	}
	return acceptedMetrics
}

func init() {
	processors.Add("breaker", func() telegraf.Processor {
		return &Breaker{}
	})
}
