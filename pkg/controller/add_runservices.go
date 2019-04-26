package controller

import (
	"github.com/barpilot/cloud-run-controller/pkg/controller/runservices"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, runservices.Add)
}
