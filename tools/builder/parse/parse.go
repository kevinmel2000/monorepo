// Package parse (tools/builder/parse) responsible for converting task.Option to Envoy.Option
package parse

import (
	"github.com/hashicorp/consul/api"
	"github.com/lab46/monorepo/tools/builder/task"
)

// Parser is struct host that responsible for the conversion.
type Parser struct {
	consul *api.Client
}

// New return new Parser with given consul api client.
func New(consulClient *api.Client) (*Parser, error) {
	return &Parser{
		consul: consulClient,
	}, nil
}

// Parse coverts task.Option to Envoy.Option
func (Parser) Parse(taskopt *task.Option) (interface{}, error) {
	return nil, nil
}
