package task

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	tt := []struct {
		Name        string
		FilePath    string
		ExpectedOpt *Option
		ExpectError bool
	}{
		{
			Name:     "test_example",
			FilePath: "example.yaml",
			ExpectedOpt: &Option{
				Build: BuildOption{
					Lang: "go",
					Ports: AppPortOption{
						HTTP: "8008",
						GRPC: "9876",
					},
					ServiceCommunication: ServiceCommunicationOption{
						HTTP: []EgressHTTPCommunication{
							EgressHTTPCommunication{
								Name:    "ongkirapp",
								Timeout: 800,
							},
							EgressHTTPCommunication{
								Name:    "kero-addr",
								Timeout: 700,
							},
						},
						GRPC: []EgressGRPCCommunication{
							EgressGRPCCommunication{
								Name:         "ongkirapp",
								ListenerPort: "8005",
								Timeout:      1000,
							},
						},
						External: []EgressExternalCommunication{
							EgressExternalCommunication{
								Name: "jne",
								Hosts: []string{
									"tcp://1.1.1.01:80",
									"tcp://1.1.1.02:80",
								},
								Timeout:     1000,
								Port:        "10311",
								OverideHost: "https://host.jne.com",
							},
						},
						Redis: []EgressRedisCommunication{
							EgressRedisCommunication{
								Name: "redis-kero",
								Hosts: []string{
									"172.123.345.567:6379",
									"172.123.345.568:6379",
								},
							},
						},
					},
				},
			},
			ExpectError: false,
		},
		{
			Name:        "test_file_not_found",
			FilePath:    "not-found-path.yaml",
			ExpectedOpt: nil,
			ExpectError: true,
		},
		{
			Name:        "test_error_unmarshal",
			FilePath:    "example-failed-unmarshal.yaml",
			ExpectedOpt: nil,
			ExpectError: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			resopt, err := ReadOptionFromPath(tc.FilePath)
			if err != nil {
				if tc.ExpectError {
					t.Logf("[Expected Error Occured] at tc: %s with err: %s\n", tc.Name, err.Error())

				} else {
					t.Errorf("[TestReadExample] Failed ad sub-test %s: %s\n", tc.Name, err.Error())
				}
			}

			assert.Equal(t, tc.ExpectedOpt, resopt)
		})
	}
}
