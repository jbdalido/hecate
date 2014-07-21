package api

import (
	"fmt"
)

type Adapter interface {
	GetEndpointsRequest(string) *Request
	GetDiscoverAllRequest() *Request
	ParseResponse([]byte) ([]*Response, error)
}

func NewAdapter(adapter string, host string) (Adapter, error) {

	if adapter == "" {
		return nil, fmt.Errorf("Adapter can't be null")
	}
	// Only mesos for now
	switch adapter {

	case "mesos":
		return NewMesos(host), nil

	default:
		return nil, fmt.Errorf("Unkown adapter %s", adapter)
	}
}
