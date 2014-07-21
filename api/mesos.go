package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Mesos struct {
	Host   string
	Client *http.Client
}

type MesosResponse struct {
	App  *MesosApp `json:"app"`
	Apps []*MesosApp
}

type MesosApp struct {
	Id        string            `json:"id"`
	Instances []*MesosEndpoints `json:"tasks"`
	//Tasks     []*MesosEndpoints `json:"tasks"`
}

type MesosEndpoints struct {
	DockerId string `json:"id"`
	Host     string `json:"host"`
	Ports    []int  `json:"ports"`
}

func NewMesos(host string) *Mesos {
	return &Mesos{
		Host: host,
	}
}

// Endpoints discovery from mesos api,
func (m *Mesos) GetEndpointsRequest(app string) *Request {

	// Return the correct request struct
	return &Request{
		URI:    fmt.Sprint("/v2/apps/" + app),
		Method: "GET",
		Body:   "",
	}
}

// Endpoints discovery from mesos api,
func (m *Mesos) GetDiscoverAllRequest() *Request {

	// Return the correct request struct
	return &Request{
		URI:    fmt.Sprint("/v2/apps/"),
		Method: "GET",
		Body:   "",
	}
}

func (m *Mesos) ParseResponse(datas []byte) ([]*Response, error) {
	// First unmarshal the json to the MesosResponse
	mesosReponse := MesosResponse{}
	err := json.Unmarshal(datas, &mesosReponse)

	if err != nil {
		return nil, fmt.Errorf("Cant parse response received from Mesos endpoint %s\n\n", err)
	}

	responses := []*Response{}

	if mesosReponse.App != nil {
		response, err := m.ParseApplication(mesosReponse.App)
		if err != nil {
			return nil, err
		}
		responses = append(responses, response)
	}

	if mesosReponse.Apps != nil {
		for _, application := range mesosReponse.Apps {
			response, err := m.ParseApplication(application)
			if err != nil {
				return nil, err
			}
			responses = append(responses, response)
		}
	}

	return responses, nil
}

func (m *Mesos) ParseApplication(application *MesosApp) (*Response, error) {

	if application.Id == "" {
		return nil, fmt.Errorf("Empty response")
	}
	// And convert it to an Hecate Response
	response := &Response{
		Name: application.Id,
	}

	for _, endpoint := range application.Instances {

		host := &Host{
			Host: endpoint.Host,
		}
		// Grab the ports allocate by mesos
		// BUG only one port get from mesos api
		if len(endpoint.Ports) > 0 {
			for _, port := range endpoint.Ports {
				host.Ports = append(host.Ports, strconv.Itoa(port))

			}
		}
		response.Endpoints = append(response.Endpoints, host)
	}
	return response, nil
}
