package api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Client struct {
	Host       string       // Hostname of the system
	HttpClient *http.Client // http client to do query
	Adapter    Adapter      // Adapter to query properly
}

func NewClient(adapter string, host string) (*Client, error) {
	// Check if the host is valid
	u, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	/*
		TODO :
		  - ssl powered request might not be workign at all
	*/
	httpClient := newHTTPClient(u, nil)

	// Build the correct adapter
	a, err := NewAdapter(adapter, host)
	if err != nil {
		return nil, err
	}

	return &Client{
		Host:       host,
		HttpClient: httpClient,
		Adapter:    a,
	}, nil
}

// Get endpoints from the defined adapter
func (c *Client) Discover(app string) ([]*Response, error) {

	// First gets the proper uri from the right adapter
	// ex: "/v1/endpoints/ID for mesos
	request := c.Adapter.GetEndpointsRequest(app)

	datas, err := c.doRequest(request)
	if err != nil {
		return nil, err
	}

	response, err := c.Adapter.ParseResponse(datas)
	if err != nil {
		return nil, err
	}

	return response, nil

}

// Discover all enpoints in a very ugly way to fix fast
func (c *Client) DiscoverAll() ([]*Response, error) {

	// Might change this it's ugly but it's late
	request := c.Adapter.GetDiscoverAllRequest()

	datas, err := c.doRequest(request)
	if err != nil {
		return nil, err
	}

	responses, err := c.Adapter.ParseResponse(datas)
	if err != nil {
		return nil, err
	}

	if len(responses) > 0 {
		for _, application := range responses {
			response, err := c.Discover(application.Name)
			if err != nil {
				return nil, err
			}
			// Uglyness
			application.Endpoints = response[0].Endpoints
		}
	}

	return responses, nil

}

// Do the http request and returning a []byte
// Need to work on that
func (c *Client) doRequest(request *Request) ([]byte, error) {
	// Bufferize the body event if it's empty
	b := bytes.NewBuffer([]byte(request.Body))
	req, err := http.NewRequest(request.Method, c.Host+request.URI, b)
	if err != nil {
		return nil, err
	}

	// Set content and accept header
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("Not found")
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("%s: %s", resp.Status, data)
	}
	return data, nil
}
