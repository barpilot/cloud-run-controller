package utils

import (
	"path"
)

func Parent(project, location string) string {
	return path.Join("projects", project, "locations", location)
}

func ServiceName(parent, service string) string {
	return path.Join(parent, "services", service)
}
