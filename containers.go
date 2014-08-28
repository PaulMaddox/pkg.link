package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/GeertJohan/go.rice"
	"github.com/fsouza/go-dockerclient"
)

func DeployContainers(client *docker.Client) (int, error) {

	// Iterate the containers/* directory
	containers := rice.MustFindBox("containers")

	containers.Walk("", func(path string, info os.FileInfo, err error) error {

		if err != nil {
			log.Fatal(err)
		}

		if !info.IsDir() {
			// Skip files
			return nil
		}

		if path == "" {
			// Skip root directory
			return nil
		}

		log.Printf("Deploying docker container: %s", path)

		archive, err := TarGz(path, info, containers)
		if err != nil {
			log.Fatalf("Unable to create docker container %s (%s)", info.Name(), err)
		}

		if err := client.BuildImage(docker.BuildImageOptions{
			Name:         info.Name(),
			InputStream:  archive,
			OutputStream: os.Stdout,
		}); err != nil {
			log.Fatalf("Failed to build docker image for %s (%s)", info.Name(), err)
		}

		return filepath.SkipDir

	})

	return 0, nil

}
