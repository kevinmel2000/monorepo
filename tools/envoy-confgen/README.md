# Envoy Config Generator

Envoy config generator is a configuration generator for `envoy-proxy`

## Generating Configuration

To generate a configuration, `envoyconf.yaml` is needed.

Run: `envoy-confgen -gen-file=envoyconf.yaml generate` to generate `envoy_config.json` 

Example of the `yaml` configuration:

```yaml
ingress:
  ssl:
    cert_chain_file: "path_to_cert"
    private_key_file: "path_to_private"
  use_remote_address: true
    - name: "local_service"
      domain: "*"
      routes:
        - name: "kero-maps"
          prefix: "/"
          prefix_rewrite: ""
          remote_address: "tcp://127.0.0.1:9000"
          timeout_ms: 0
          headers:
            - name: "content-type"
              value: "application/json"
egress:
  service_to_service:
    - name: "kero"
      timeout_ms: 500
      cluster_hosts: ['tcp://kero.aliyun.consul']
  grpc:
    - name: "svc"
      port_listener: ":9132"
      timeout_ms: 250
      hosts:
        - "tcp://svc2-1:9211"
  external_virtual_hosts:
    - name: "something"
      address: "tcp://127.0.0.1:9204"
      hosts:
          - name: "something"
            domain: "*"
            remote_address: "something.id:80"
            rewrite_address: "something.id"
            ssl: false
      cluster_type: "logical_dns" 
```