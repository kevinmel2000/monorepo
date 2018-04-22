package httpclient

import (
	"net/http"
	"testing"
)

func TestHTTPReqValidate(t *testing.T) {
	cases := []struct {
		req     HTTPReq
		sidecar bool
		expect  error
	}{
		// missing http method
		{
			req: HTTPReq{
				URL:       "someurl",
				URLParams: map[string]string{"key1": "value1"},
			},
			sidecar: false,
			expect:  ErrMissingHTTPMethod,
		},
		// invalid sidecar address
		{
			req: HTTPReq{
				Method:      http.MethodPost,
				URL:         "something:9100",
				HostHerader: "something",
				URLParams:   map[string]string{"key1": "value1"},
			},
			sidecar: true,
			expect:  ErrInvalidSidecarAddress,
		},
		// invalid hostheader for sidecar
		{
			req: HTTPReq{
				Method:    http.MethodPost,
				URL:       "localhost:9100",
				URLParams: map[string]string{"key1": "value1"},
			},
			sidecar: true,
			expect:  ErrHostHeaderEmpty,
		},
		// valid sidecar address using localhost
		{
			req: HTTPReq{
				Method:      http.MethodPost,
				URL:         "localhost:9100",
				HostHerader: "header",
				URLParams:   map[string]string{"key1": "value1"},
			},
			sidecar: true,
			expect:  nil,
		},
		// valid sidecar address using 127.0.0.1
		{
			req: HTTPReq{
				Method:      http.MethodPost,
				URL:         "127.0.0.1:9100",
				HostHerader: "header",
				URLParams:   map[string]string{"key1": "value1"},
			},
			sidecar: true,
			expect:  nil,
		},
	}

	for _, val := range cases {
		if val.sidecar {
			val.req.useSidecar()
		}
		if err := val.req.Validate(); val.expect != err {
			t.Errorf("expecting %s but got %s", val.expect, err.Error())
		}
	}
}
