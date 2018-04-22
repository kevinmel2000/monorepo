// Package task (tools/builder/task) responsible for unmarshaling tkp-task.yaml to task.Option struct.
package task

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Option represent configuration available in tkp-task.yaml file
type Option struct {
	Build BuildOption `yaml:"build"`
}

// BuildOption represent options for Building the app.
// This option will be read to generate app configurations (such as envoy configuration)
type BuildOption struct {
	Lang                 string                     `yaml:"lang"`
	Ports                AppPortOption              `yaml:"ports"`
	ServiceCommunication ServiceCommunicationOption `yaml:"service_communication"`
}

// AppPortOption stores HTTP and gRPC port that being exposed by the app.
type AppPortOption struct {
	HTTP string `yaml:"http"`
	GRPC string `yaml:"grpc"`
}

// ServiceCommunicationOption stores possible communications made by the apps.
// Currently it supports HTTP, gRPC, External and Redis communication.
// This config would be used later on generating egress configurations.
type ServiceCommunicationOption struct {
	HTTP     []EgressHTTPCommunication     `yaml:"http"`
	GRPC     []EgressGRPCCommunication     `yaml:"grpc"`
	External []EgressExternalCommunication `yaml:"external"`
	Redis    []EgressRedisCommunication    `yaml:"redis"`
}

// EgressHTTPCommunication stores possible via HTTP communication to internal service.
type EgressHTTPCommunication struct {
	// Name should be service name from service discovery.
	// This will be used to lookup hosts in service discovery
	Name    string `yaml:"name"`
	Timeout int64  `yaml:"timeout_ms"`
}

// EgressGRPCCommunication stores possible via gRPC communication to internal service.
type EgressGRPCCommunication struct {
	Name         string `yaml:"name"`
	ListenerPort string `yaml:"listener_port"`
	Timeout      int64  `yaml:"timeout_ms"`
}

// EgressExternalCommunication stores possible via HTTP communication to external APIs.
// 3rd party APIs such as Logistic's or Google's should be stored in this.
type EgressExternalCommunication struct {
	Name        string   `yaml:"name"`
	Hosts       []string `yaml:"hosts"`
	Timeout     int64    `yaml:"timeout_ms"`
	Port        string   `yaml:"port"`
	OverideHost string   `yaml:"overide_host"`
}

// EgressRedisCommunication stores redis connection(s) made by the app.
// Load balancer configuration would be made from given Hosts.
type EgressRedisCommunication struct {
	Name  string   `yaml:"name"`
	Hosts []string `yaml:"hosts"`
}

// ReadOptionFromPath read files in given path and unmarshal it using readBytesConfig.
func ReadOptionFromPath(path string) (*Option, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return readBytesOption(f)
}

func readBytesOption(b []byte) (*Option, error) {
	var opt Option
	err := yaml.Unmarshal(b, &opt)
	if err != nil {
		return nil, err
	}
	return &opt, nil
}
