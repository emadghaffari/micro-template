# micro-template
a template for grpc and http calls

### configs
use make config for generate basci config

#### services:
 - service tracer: jaeger
 - config server: etcd
 - cmd: cobra
 - configs: from file(dev), from config server(prod)
 - logger: zap
 - database: postgres(go-pg)
 - database ui: pgAdmin

#### config file
rename config.example.yaml to config.yaml,
for generate config file from config example use "make config" command. 
