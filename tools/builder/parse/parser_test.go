package parse

import (
	"os"
	"testing"

	"gopkg.in/jarcoal/httpmock.v1"

	"github.com/hashicorp/consul/api"
)

var parserMock *Parser
var defaultConsulAddr = api.DefaultConfig().Address

func TestMain(m *testing.M) {
	httpmock.Activate()
	defer httpmock.Deactivate()

	consul, _ := api.NewClient(api.DefaultConfig())
	parserMock, _ = New(consul)

	os.Exit(m.Run())
}

func TestParse(t *testing.T) {
	parserMock.Parse(nil)
}
