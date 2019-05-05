package controller

import (
	"github.com/barpilot/cloud-run-controller/pkg/controller/servicecreator"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, servicecreator.Add)
}
