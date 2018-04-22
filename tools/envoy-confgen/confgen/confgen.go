package confgen

import (
	"errors"
	"html/template"
	"os"
)

// IngressConfig for ingress traffic
type IngressConfig struct {
	SSL   IngressSSL    `yaml:"ssl"`
	Hosts []IngressHost `yaml:"virtual_hosts"`
}

// Validate IngressConfGen
func (ing *IngressConfig) Validate() error {
	err := ing.SSL.Validate()
	if err != nil {
		return err
	}
	return nil
}

// IngressSSL to serve ingress SSL traffic
type IngressSSL struct {
	Use            bool   `yaml:"-"`
	CertChainFile  string `yaml:"cert_chain_file"`
	PrivateKeyFile string `yaml:"private_key_file"`
}

// Validate SSLConfGen
func (ssl *IngressSSL) Validate() error {
	if (ssl.CertChainFile != "" && ssl.PrivateKeyFile == "") || (ssl.CertChainFile == "" && ssl.PrivateKeyFile != "") {
		return errors.New("[ingress][ssl] one of ssl configuration cannot be empty")
	}
	if ssl.CertChainFile != "" && ssl.PrivateKeyFile != "" {
		ssl.Use = true
	}
	return nil
}

// IngressHost to serve ingress host traffic
type IngressHost struct {
	Name   string             `yaml:"name"`
	Domain string             `yaml:"domain"`
	Routes []IngressHostRoute `yaml:"routes"`
}

// IngressHostRoute to map route
type IngressHostRoute struct {
	Name          string `yaml:"name"`
	Prefix        string `yaml:"prefix"`
	PrefixRewrite string `yaml:"prefix_rewrite"`
	RemoteAddress string `yaml:"remote_address"`
	TimeoutMs     int    `yaml:"timeout_ms"`
}

// EgressConfig for egress config
type EgressConfig struct {
	ServiceToService []EgressServiceToService `yaml:"service_to_service"`
	Grpc             []EgressGrpc             `yaml:"grpc"`
	External         []EgressExternalHost     `yaml:"external_host"`
}

// EgressServiceToService serve service to service traffic
type EgressServiceToService struct {
	Name         string   `yaml:"name"`
	TimeoutMs    int      `yaml:"timeout_ms"`
	ClusterHosts []string `yaml:"cluster_hosts"`
}

// EgressExternalHost service to external service
type EgressExternalHost struct {
	Name        string       `yaml:"name"`
	Address     string       `yaml:"address"`
	Hosts       []EgressHost `yaml:"hosts"`
	ClusterType string       `yaml:"cluster_type"`
}

// EgressHost to map host of EgressExternalHost
type EgressHost struct {
	Name           string `yaml:"name"`
	Domain         string `yaml:"domain"`
	RemoteAddress  string `yaml:"remote_address"`
	RewriteAddress string `yaml:"rewrite_address"`
	SSL            bool   `yaml:"ssl"`
}

// EgressGrpc struct
type EgressGrpc struct {
	Name         string   `yaml:"name"`
	PortListener string   `yaml:"port_listener"`
	TimeoutMs    int      `yaml:"timeout_ms"`
	Hosts        []string `yaml:"hosts"`
}

// Cluster struct
type Cluster struct {
	Name      string   `yaml:"name"`
	TimeoutMs int      `yaml:"connect_timeout_ms"`
	Type      string   `yaml:"type"`
	LbType    string   `yaml:"lb_type"`
	Features  string   `yaml:"features"`
	Hosts     []string `yaml:"hosts"`
}

// Generator to generate envoy config
type Generator struct {
	Ingress  IngressConfig `yaml:"ingress"`
	Egress   EgressConfig  `yaml:"egress"`
	Clusters []Cluster     `yaml:"cluster"`
}

// Validate ConfGen
func (gen *Generator) Validate() error {
	err := gen.Ingress.Validate()
	if err != nil {
		return err
	}
	return nil
}

// CreateClusters create cluster from ingress and egress
func (gen *Generator) CreateClusters() error {
	// ingress cluster
	for _, host := range gen.Ingress.Hosts {
		for _, r := range host.Routes {
			c := Cluster{
				Name:      r.Name,
				TimeoutMs: r.TimeoutMs,
				Type:      "strict_dns",
				LbType:    "round-robin",
				Hosts:     []string{r.RemoteAddress},
			}
			gen.Clusters = append(gen.Clusters, c)
		}
	}

	// egress cluster
	for _, sts := range gen.Egress.ServiceToService {
		c := Cluster{
			Name:      sts.Name,
			TimeoutMs: sts.TimeoutMs,
			Type:      "strict_dns",
			LbType:    "round-robin",
			Hosts:     sts.ClusterHosts,
		}
		gen.Clusters = append(gen.Clusters, c)
	}
	for _, grpc := range gen.Egress.Grpc {
		c := Cluster{
			Name:      grpc.Name,
			TimeoutMs: grpc.TimeoutMs,
			Type:      "strict_dns",
			LbType:    "round-robin",
			Hosts:     grpc.Hosts,
		}
		gen.Clusters = append(gen.Clusters, c)
	}
	// for _, external := range gen.Egress.ExternalVirtualHosts {
	// 	c := Cluster{
	// 		Name: external.Name,
	// 	}
	// }
	return nil
}

var templates = []string{ingressFilters, ingresSSLFilters, ingressRouteConfig, envoyTemplate}

// GenerateToFile for transforming config struct into a file
func GenerateToFile(gen Generator, filename string) error {
	err := gen.Validate()
	if err != nil {
		return err
	}
	err = gen.CreateClusters()
	if err != nil {
		return err
	}

	var t *template.Template
	tmpl := template.New("envoy")
	for _, val := range templates {
		t, err = tmpl.Parse(val)
		if err != nil {
			return err
		}
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return t.Execute(f, gen)
}
