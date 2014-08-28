package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GeertJohan/go.rice"
	"github.com/fsouza/go-dockerclient"
	"github.com/ogier/pflag"
)

var _bind *string = pflag.StringP("bind", "b", "0.0.0.0", "Bind address to listen on")
var _port *int = pflag.IntP("port", "p", 3000, "HTTP port to listen on")
var _docker *string = pflag.StringP("docker", "d", "unix:///var/run/docker.sock", "Docker server address")

func main() {

	// Parse the CLI arguments
	pflag.Parse()

	log.Printf("Connecting to docker at %s", *_docker)

	client, err := docker.NewClient(*_docker)
	if err != nil {
		log.Printf("Failed to connect to docker (%s)", err)
	}

	version, err := client.Version()
	if err != nil {
		log.Fatalf("Failed to connect to docker (%s)", err)
	}

	log.Printf("Successfully connected to docker v%s", version.Get("Version"))

	count, err := DeployContainers(client)
	if err != nil {
		log.Fatalf("Failed to deploy docker containers (%s)", err)
	}

	log.Printf("Successfully deployed %d docker containers", count)

	router := NewRouter()
	bind := fmt.Sprintf("%s:%d", *_bind, *_port)
	log.Printf("Starting HTTP server on %s", bind)
	http.Handle("/", http.FileServer(rice.MustFindBox("public").HTTPBox()))
	log.Fatal(http.ListenAndServe(bind, router))

}
