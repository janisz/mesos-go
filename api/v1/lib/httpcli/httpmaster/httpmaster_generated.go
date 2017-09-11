package httpmaster

// go generate -import github.com/mesos/mesos-go/api/v1/lib/master -import github.com/mesos/mesos-go/api/v1/lib/master/calls -type C:master.Call:master.Call{Type:master.Call_GET_METRICS}
// GENERATED CODE FOLLOWS; DO NOT EDIT.

import (
	"context"

	"github.com/mesos/mesos-go/api/v1/lib"
	"github.com/mesos/mesos-go/api/v1/lib/client"
	"github.com/mesos/mesos-go/api/v1/lib/httpcli"

	"github.com/mesos/mesos-go/api/v1/lib/master"
	"github.com/mesos/mesos-go/api/v1/lib/master/calls"
)

// ResponseClassifier determines the appropriate response class for the given call.
type ResponseClassifier func(*master.Call) (client.ResponseClass, error)

// ClientFunc sends a Request to Mesos and returns the generated Response.
type ClientFunc func(client.Request, client.ResponseClass, ...httpcli.RequestOpt) (mesos.Response, error)

// DefaultResponseClassifier is a pluggable classifier.
var DefaultResponseClassifier = ResponseClassifier(classifyResponse)

// NewSender generates a sender that uses the Mesos v1 HTTP API for encoding/decoding requests/responses.
// The ResponseClass is inferred from the first object generated by the given Request.
func NewSender(cf ClientFunc) calls.Sender {
	return calls.SenderFunc(func(ctx context.Context, r calls.Request) (mesos.Response, error) {
		var (
			obj     = r.Call()
			rc, err = DefaultResponseClassifier(obj)
		)
		if err != nil {
			return nil, err
		}

		var req client.Request

		switch r := r.(type) {
		case calls.RequestStreaming:
			req = calls.Push(r, obj)
		default:
			req = calls.NonStreaming(obj)
		}

		return cf(req, rc, httpcli.Context(ctx))
	})
}
