package confgen

var envoyTemplate = `
{
	"listeners": [
		{{ if .Ingress }} {{ if .Ingress.SSL.Use }}
		{
			"address": "tcp://0.0.0.0:443",
			{{ 	template "ingress_ssl_filters" . }}
			,"ssl_context": {
				"cert_chain_file": "{{ .Ingress.SSL.CertChainFile }}",
				"private_key_file": "{{ .Ingress.SSL.PrivateKeyFile }}"
			}
		}, {{ end }}
		{
			"address": "tcp://0.0.0.0:80",
				{{ template "ingress_filters" . }}
		} {{ end }}
		{{ if .Egress }}{{ if .Egress.ServiceToService }}
		,{
			"address": "tcp://0.0.0.0:9001",
			"filters": [
				{
					"type": "read",
					"name": "http_connection_manager",
					"config": {
						"codec_type": "auto",
						"stat_prefix": "egress_http",
						"add_user_agent": true,
						"idle_timeout_s": 840,
						"use_remote_address": true,
						"forward_client_cert" : "always_forward_only",
						"route_config": {
							"virtual_hosts": [
								{{ range $stskey, $stshost := .Egress.ServiceToService }}
								{
									"name": "{{ $stshost.Name }}",
									"domains": ["{{ $stshost.Name }}"],
									"routes": [
										{
											"timeout_ms": 200,
											"prefix": "/",
											"retry_policy": {
												"retry_on": "connect-failure"
											},
											"cluster": "egress_{{ $stshost.Name }}"
										}
									]
								}
								{{ end }}
							]
						},
						"filters": [
							{   "name": "rate_limit",
								"tcp_filter_enabled" : 2,
								"config": {
									"domain": "envoy_service_to_service"
								}
							},
							{"name": "grpc_http1_bridge", "config": {}},
							{"name": "router", "config": {}}
						]
					}	
				}
			]
		}
		{{ end }} {{ if .Egress.External }}
		{{ range $extkey, $exthost := .Egress.External }}
		,{
			"address": "{{ $exthost.Address }}",
			"filters": [
				{
					"type": "read",
					"name": "http_connection_manager",
					"config": {
						"codec_type": "auto",
						"access_log": [{
							"path": "/var/log/envoy/{{ $exthost.Name }}-ext.access.log"
						}],
						"stat_prefix": "{{ $exthost.Name }}",
						"route_config": {
							"virtual_hosts": [
								{{ range $hkey, $h :=  $exthost.Hosts }}
								{
									"name": "egress_{{ $h.Name }}",
									"domains": ["{{ $h.Domain }}"],
									"routes": [
										{
											"timeout_ms": 0,
											"prefix": "/",
											"cluster": "{{ $h.Name }}",
											"retry_policy": { "retry_on": "connect-failure" }
											{{ if $h.RewriteAddress }}
											, "host_rewrite": "{{ $h.RewriteAddress }}"
											{{ end }}
										}
									]
								}
								{{ end }}
							]
						},
						"filters": [
							{"name": "router","config": {}}
						]
					}
				}
			]
		}
		{{ end }} {{ end }} {{ end }}
	],
	"admin": {
		"access_log_path": "/dev/null",
		"address": "tcp://0.0.0.0:8001"
	},
	"statsd_udp_ip_address": "0.0.0.0:9125",
	"cluster_manager": {
		"clusters": [
			{{ range $ckey, $c := .Clusters }}
			{{ if $ckey }},{{ end }}{
				"name": "{{ $c.Name }}",
				"connect_timeout_ms": {{ $c.TimeoutMs }},
				"type": "{{ $c.Type }}",
				"lb_type": "{{ $c.LbType }}",
				"hosts": [
					{{ range $k, $h := $c.Hosts }} {{ if $k }},{{end}}{ "url": "{{ $h }}" } {{ end }}
				] 
			}
			{{ end }}
		]
	}
}
`

var ingressFilters = `
{{ define "ingress_filters" }}
			"filters": [
				{
					"type": "read",
					"name": "http_connection_manager",
					"config": {
						"codec_type": "auto",
						"stat_prefix": "ingress_https",
						"access_log": [{
							"path": "/var/log/envoy/envoy.access.log"
						}],
						"use_remote_address": true,
						{{ template "ingress_route_config" . }}
						,"filters": [
							{
								"type": "decoder",
								"name": "router",
								"config": {}
							}
						]
					}
				}
			]
{{ end }}
`

var ingresSSLFilters = `
{{ define "ingress_ssl_filters" }}
			"filters": [
				{
					"type": "read",
					"name": "http_connection_manager",
					"config": {
						"codec_type": "auto",
						"stat_prefix": "ingress_https",
						"access_log": [{
							"path": "/var/log/envoy/envoy.https.access.log"
						}],
						"use_remote_address": true,
						{{ template "ingress_route_config" . }}
						,"filters": [
							{
								"type": "decoder",
								"name": "router",
								"config": {}
							}
						]
					}
				}
			]
{{ end }}
`

var ingressRouteConfig = `
{{ define "ingress_route_config" }}
						"route_config": {	
							"virtual_hosts": [
								{{ range $hostskey, $host :=  .Ingress.Hosts  }}
								{
									"name": "{{ $host.Name }}",
									"domains": ["{{ $host.Domain }}"],
									"routes": [
										{{ range $routeskey, $route := $host.Routes }}
										{
											"timeout_ms": "{{ $route.TimeoutMs }}",
											"prefix": "{{ $route.Prefix }}",
											"cluster": "{{ $route.Name}} "{{ if $route.PrefixRewrite }},
											"prefix_rewrite": "{{ $route.PrefixRewrite }}"
											{{ end }}
										} {{ end }}
									]
								} {{ end }}
							]
						}
{{ end }}
`
