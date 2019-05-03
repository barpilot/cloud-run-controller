package utils

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// finalizable is an kubernetes resource api object that supports finalizers
type finalizable interface {
	runtime.Object
	GetFinalizers() []string
	SetFinalizers(finalizers []string)
	GetDeletionTimestamp() *metav1.Time
}

// Finalizer manages the finalizers for resources in kubernetes
type Finalizer struct {
	value string
}

// NewFinalizer constructs a new finalizer manager
func NewFinalizer(finalizerValue string) *Finalizer {
	return &Finalizer{
		value: finalizerValue,
	}
}

// Add adds a finalizer to an object
func (c *Finalizer) Add(resource finalizable) {
	finalizers := append(resource.GetFinalizers(), c.value)
	resource.SetFinalizers(finalizers)
}

// Remove removes a finalizer from an object
func (c *Finalizer) Remove(resource finalizable) {
	finalizers := resource.GetFinalizers()
	for idx, finalizer := range finalizers {
		if finalizer == c.value {
			finalizers = append(finalizers[:idx], finalizers[idx+1:]...)
			break
		}
	}
	resource.SetFinalizers(finalizers)
}

// IsDeletionCandidate checks if the resource is a candidate for deletion
func (c *Finalizer) IsDeletionCandidate(resource finalizable) bool {
	return resource.GetDeletionTimestamp() != nil && c.getIndex(resource) != -1
}

// NeedToAdd checks if the resource should have but does not have the finalizer
func (c *Finalizer) NeedToAdd(resource finalizable) bool {
	return resource.GetDeletionTimestamp() == nil && c.getIndex(resource) == -1
}

func (c *Finalizer) getIndex(resource finalizable) int {
	for i, v := range resource.GetFinalizers() {
		if v == c.value {
			return i
		}
	}
	return -1
}
