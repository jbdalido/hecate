package main

import (
	"encoding/json"
	"fmt"
	flag "github.com/dotcloud/docker/pkg/mflag"
	"github.com/jbdalido/hecate/api"
	"log"
	"os"
)

// Command line

func main() {

	// setup command line arguments
	var (
		kubernetes bool
		mesos      bool
		scheduler  string
		host       string
		app        string
		version    bool
		discover   bool
	)

	// Override flag with docker mflag
	// TODO :
	//      - need to simplify flags, there is way too much of them
	flag.StringVar(&app, []string{"a", "-app"}, "", "App id used by your scheduler")
	flag.StringVar(&host, []string{"h", "-host"}, "", "Host to use with the scheduler")

	flag.BoolVar(&kubernetes, []string{"k", "-kubernetes"}, false, "Connect to Kubernetes endpoint")
	flag.BoolVar(&mesos, []string{"m", "-mesos"}, false, "Connect to Mesos endpoint")
	flag.StringVar(&scheduler, []string{"s", "-scheduler"}, "mesos", "Literral adapter")

	flag.BoolVar(&discover, []string{"d", "-discover"}, false, "Discover all apps and endpoints")
	flag.BoolVar(&version, []string{"v", "-version"}, false, "Show Version")

	log.SetPrefix("[WIP]")
	log.SetFlags(0)
	flag.Parse()

	if version {
		showVersion()
		os.Exit(0)
	}

	if !discover && app == "" {
		log.Fatal("Application can't be null")
	}

	if host == "" {
		log.Fatal("Host can't be null")
	}

	if kubernetes {
		scheduler = "kubernetes"
	}

	client, err := api.NewClient(scheduler, host)
	if err != nil {
		log.Fatalf("Error %s", err)
	}

	// TODO
	// - need to delegate the String() function
	//   to the client
	// - need to handle nginx conf via template
	//   or via etcd to be purge via confd
	//   or find another way
	if app != "" && !discover {
		response, err := client.Discover(app)
		if err != nil {
			log.Fatalf("%s", err)
		}
		r, _ := json.Marshal(response)
		fmt.Printf("%s\n", r)
	}

	if discover {
		response, err := client.DiscoverAll()
		if err != nil {
			log.Fatalf("%s", err)
		}
		r, _ := json.Marshal(response)
		fmt.Printf("%s\n", r)
	}

}

func showVersion() {
	fmt.Printf("Hecate v0.1.0 ©Why are u so serious ?")
}
