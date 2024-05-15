Name: {{.serviceName}}
Host: {{.host}}
Port: {{.port}}
Mode: dev
Timeout: 30000

# prometheus-exporter
Prometheus:
  Host: 0.0.0.0
  Port: 9091
  Path: /metrics
  
# telemetry
Telemetry:
  Name: {{.serviceName}}
  Endpoint: zincobserve.dev-tools:5080
  Sampler: 1.0
  Batcher: otlphttp
  OtlpHeaders:
    Authorization: "Basic xxxxxxxxxxxxxxxxxxxxxxxxx"
  OtlpHttpPath: /api/default/traces

