package run

import (
	"context"
	"fmt"
	"net/url"

	"github.com/spf13/pflag"
	"google.golang.org/api/option"
	run "google.golang.org/api/run/v1alpha1"
)

type RunManager struct {
	project string
	service *run.APIService
}

func NewRunManager(projectName string) (*RunManager, error) {
	if projectName == "" {
		return &RunManager{}, fmt.Errorf("empty project")
	}

	ctx := context.Background()
	rm := &RunManager{project: projectName}

	api := pflag.Lookup("apikey").Value.String()

	runService, err := run.NewService(ctx, option.WithAPIKey(api))
	if err != nil {
		return rm, err
	}
	rm.service = runService

	return rm, err
}

func (rm *RunManager) getAllLocations() ([]run.Location, error) {
	locations := []run.Location{}

	list, err := rm.service.Projects.Locations.List(fmt.Sprintf("projects/%s", rm.project)).Do()

	if err != nil {
		return locations, err
	}

	for _, location := range list.Locations {
		locations = append(locations, *location)
	}

	return locations, nil
}

type RunService struct {
	Name     string
	Hostname string
}

func (rm *RunManager) GetAllServices() ([]RunService, error) {
	runServices := []RunService{}

	locations, err := rm.getAllLocations()
	if err != nil {
		return runServices, err
	}
	for _, location := range locations {
		services, err := rm.service.Projects.Locations.Services.List(location.Name).Do()
		if err != nil {
			return runServices, err
		}

		for _, item := range services.Items {
			u, err := url.Parse(item.Status.Address.Hostname)
			if err != nil {
				return runServices, err
			}

			runServices = append(runServices, RunService{Name: item.Metadata.Name, Hostname: u.Host})
		}
	}
	return runServices, nil
}
