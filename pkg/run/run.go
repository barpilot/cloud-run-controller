package run

import (
	"context"
	"fmt"
	"net/url"
	"unsafe"

	"github.com/barpilot/cloud-run-controller/pkg/utils"
	"github.com/spf13/pflag"
	"google.golang.org/api/option"

	runApi "google.golang.org/api/run/v1alpha1"
)

type RunManager struct {
	project string
	service *runApi.APIService
}

func NewRunManager(projectName string) (*RunManager, error) {
	if projectName == "" {
		return &RunManager{}, fmt.Errorf("empty project")
	}

	ctx := context.Background()
	rm := &RunManager{project: projectName}

	api := pflag.Lookup("apikey").Value.String()

	runService, err := runApi.NewService(ctx, option.WithAPIKey(api))
	if err != nil {
		return rm, err
	}
	rm.service = runService

	return rm, err
}

func (rm *RunManager) getAllLocations() ([]runApi.Location, error) {
	locations := []runApi.Location{}

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

func (rm *RunManager) CreateOrUpdate(parent string, service *Service) (*Service, error) {
	runApiSvc := (*runApi.Service)(unsafe.Pointer(service))
	name := utils.ServiceName(parent, service.Metadata.Name)
	result, err := rm.service.Projects.Locations.Services.Get(name).Do()
	if err != nil {
		result, err = rm.service.Projects.Locations.Services.Create(parent, runApiSvc).Do()
	} else {
		result, err = rm.service.Projects.Locations.Services.ReplaceService(name, runApiSvc).Do()
	}

	return (*Service)(unsafe.Pointer(result)), err
}

func (rm *RunManager) Delete(resource string, service Service) error {
	_, err := rm.service.Projects.Locations.Services.Delete(resource).Do()
	return err
}

func (rm *RunManager) SetIamPolicy(resource string, policy *IamPolicy) error {
	p := (*runApi.Policy)(unsafe.Pointer(policy))
	_, err := rm.service.Projects.Locations.Services.SetIamPolicy(resource, &runApi.SetIamPolicyRequest{Policy: p}).Do()
	return err
}
