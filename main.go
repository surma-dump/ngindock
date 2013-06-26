package main

import (
	"encoding/json"
	"github.com/dotcloud/docker"
	"github.com/voxelbrain/goptions"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"text/template"
)

const (
	VERSION = "1.0.0"
)

var (
	options = struct {
		Docker     string        `goptions:"-H, --docker, description='Address of docker daemon'"`
		Template   *os.File      `goptions:"-t, --template, description='Template to render', rdonly"`
		Output     string        `goptions:"-o, --output, description='File to render to'"`
		DontReload bool          `goptions:"--dont-reload, description='Dont make nginx reload its configuration'"`
		Help       goptions.Help `goptions:"-h, --help, description='Show this help'"`
	}{
		Docker: "localhost:4243",
		Output: "/etc/nginx/conf.d/docker.conf",
	}
)

func main() {
	goptions.ParseAndFail(&options)
	output, err := os.Create(options.Output)
	if err != nil {
		log.Fatalf("Could not open output file %s: %s", options.Output, err)
	}
	defer output.Close()

	data := []byte(DefaultTemplate)
	if options.Template != nil {
		data, err = ioutil.ReadAll(options.Template)
		if err != nil {
			log.Fatalf("Could not read template file %s: %s", options.Template, err)
		}
	}
	tpl, err := template.New("nginx").Parse(string(data))
	if err != nil {
		log.Fatalf("Could not parse template: %s", err)
	}

	containers, err := allContainers(options.Docker)
	if err != nil {
		log.Fatalf("Could not gather info about running containers: %s", err)
	}
	containers = filterContainers(containers)
	err = tpl.Execute(output, containers)
	if err != nil {
		log.Fatalf("Could not render template: %s", err)
	}

	if !options.DontReload {
		err := exec.Command("nginx", "-s", "reload").Run()
		if err != nil {
			log.Fatalf("Could not make nginx reload its configuration: %s", err)
		}
	}

}

func allContainers(addr string) ([]docker.Container, error) {
	resp, err := http.Get("http://" + addr + "/containers/json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var containers []docker.APIContainers
	err = json.NewDecoder(resp.Body).Decode(&containers)
	if err != nil {
		return nil, err
	}

	allContainerDetails := make([]docker.Container, 0, len(containers))
	for _, container := range containers {
		details, err := containerDetails(addr, container.ID)
		if err != nil {
			return nil, err
		}
		allContainerDetails = append(allContainerDetails, details)
	}
	return allContainerDetails, nil
}

func containerDetails(addr string, id string) (docker.Container, error) {
	resp, err := http.Get("http://" + addr + "/containers/" + id + "/json")
	if err != nil {
		return docker.Container{}, err
	}
	defer resp.Body.Close()

	var container docker.Container
	return container, json.NewDecoder(resp.Body).Decode(&container)
}

func filterContainers(containers []docker.Container) []docker.Container {
	r := make([]docker.Container, 0, len(containers))
	for _, container := range containers {
		if !container.State.Running {
			continue
		}
		if _, ok := container.NetworkSettings.PortMapping["80"]; !ok {
			continue
		}
		r = append(r, container)
	}
	return r
}
