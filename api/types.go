/*
   Representing the response api to the end users
*/
package api

// Application name(id) and
// host endpoint to call
type Application struct {
	Name string // App id
	Host string // Host Adapter
}

type Request struct {
	Method string
	URI    string
	Body   string
}

// Response is the returned json type
// returned from every queries
type Response struct {
	Name      string  // AppId
	Endpoints []*Host // List of endpoints
}

// Eeach response contain at least ONE
// host with at least ONE port
type Host struct {
	Host  string   // Application host running
	Ports []string // on specific port
}
