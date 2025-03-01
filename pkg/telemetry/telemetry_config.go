package telemetry

type TelemetryConfig struct {
	ServiceName  string  `mapstructure:"service_name" json:"service_name" yaml:"service_name"`
	CollectorURL string  `mapstructure:"collector_url" json:"collector_url" yaml:"collector_url"`
	TraceSample  float64 `mapstructure:"trace_sample" json:"trace_sample" yaml:"trace_sample"`
}
