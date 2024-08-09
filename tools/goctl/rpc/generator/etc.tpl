Name: {{.serviceName}}.rpc
ListenOn: 0.0.0.0:8080
Etcd:
  Hosts:
  - 127.0.0.1:2379
  Key: {{.serviceName}}.rpc
Timeout: 30000


# prometheus-exporter
Prometheus:
  Host: 0.0.0.0
  Port: 9091
  Path: /metrics

# telemetry
Telemetry:
  Name: {{.serviceName}}.rpc
  Endpoint: zincobserve.dev-tools:5080
  Sampler: 1.0
  Batcher: otlphttp
  OtlpHeaders:
    Authorization: "Basic xxxxxxxxxxxxxxxxxxxxxxxxx"
  OtlpHttpPath: /api/default/traces

# mysql
Db:
  DriverName: mysql
  Dsn: "demo:demo@tcp(127.0.0.1:3306)/demo?parseTime=True&charset=utf8mb4"
  Debug: true

# redis
BizRedis:
  Host: 127.0.0.1
  Port: 6379
  Password: xxxxxxxxxxxxxxxx
  Db: 0
  Prefix: {{.serviceName}}

# memory
BizMemory:
  DefaultExpirationMinute: 60
  CleanUpIntervalMinute: 60
  Prefix: {{.serviceName}}

# faktory
Faktory:
  # docker run --rm -it -v ./faktory:/var/lib/faktory -e "FAKTORY_PASSWORD=123456" -p 127.0.0.1:7419:7419 -p 127.0.0.1:7420:7420 contribsys/faktory:latest /faktory -b :7419 -w :7420 -e production
  Url: tcp://:xxxxxxxxxxxx@127.0.0.1:7419
  Sender:
    PoolCapacity: 1
  Worker:
    Concurrency: 20
    PullFromQueuesWithPriority:
      default: 1
      
# nsq
Nsq:
  Sender:
    NsqdAddrs: ["nsqd.queue:4150"]
  Worker:
    NsqLookupdAddrs: ["nsqlookupd.queue:4161"]
    MaxInFlight: 50
    PullFromQueuesWithPriority:
      default: 1

