package telemetry

type TelemetryConfig struct {
	ServiceName        string  `mapstructure:"service_name" json:"service_name" yaml:"service_name"`
	ServiceVersion     string  `mapstructure:"service_version" json:"service_version" yaml:"service_version"`
	ServiceNamespace   string  `mapstructure:"service_namespace" json:"service_namespace" yaml:"service_namespace"`
	ServiceEnvironment string  `mapstructure:"service_environment" json:"service_environment" yaml:"service_environment"`
	CollectorType      string  `mapstructure:"collector_type" json:"collector_type" yaml:"collector_type"`
	CollectorURL       string  `mapstructure:"collector_url" json:"collector_url" yaml:"collector_url"`
	CollecteInterval   string  `mapstructure:"collecte_interval" json:"collecte_interval" yaml:"collecte_interval"`
	CollecteTimeout    string  `mapstructure:"collecte_timeout" json:"collecte_timeout" yaml:"collecte_timeout"`
	TraceSample        float64 `mapstructure:"trace_sample" json:"trace_sample" yaml:"trace_sample"`
}
