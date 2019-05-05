package service

import (
	"context"
	"reflect"
	"time"

	cloudruncontrollerv1alpha1 "github.com/barpilot/cloud-run-controller/pkg/apis/cloudruncontroller/v1alpha1"
	"github.com/barpilot/cloud-run-controller/pkg/run"
	"github.com/barpilot/cloud-run-controller/pkg/utils"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_service")

const requeueAfter = time.Minute * 5

const (
	serviceFinalizer   = "finalizer.service.cloud-run-controler.barpilot.io/v1alpha1"
	annotationDeletion = "removeOnDelete"
)

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Service Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileService{
		client:    mgr.GetClient(),
		scheme:    mgr.GetScheme(),
		finalizer: utils.NewFinalizer(serviceFinalizer),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("service-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Service
	err = c.Watch(&source.Kind{Type: &cloudruncontrollerv1alpha1.Service{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}
	return nil
}

var _ reconcile.Reconciler = &ReconcileService{}

// ReconcileService reconciles a Service object
type ReconcileService struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client    client.Client
	scheme    *runtime.Scheme
	finalizer *utils.Finalizer
}

// Reconcile reads that state of the cluster for a Service object and makes changes based on the state read
// and what is in the Service.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileService) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Service")

	// Fetch the Service instance
	currentInstance := &cloudruncontrollerv1alpha1.Service{}
	err := r.client.Get(context.TODO(), request.NamespacedName, currentInstance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	instance := currentInstance.DeepCopy()

	// Be sure namespace is correctly set
	if instance.Spec.Service.Metadata.Namespace == "" {
		instance.Spec.Service.Metadata.Namespace = instance.Spec.Project
	}

	rm, err := run.NewRunManager(instance.Spec.Project)
	if err != nil {
		return reconcile.Result{}, err
	}

	parent := utils.Parent(instance.Spec.Project, instance.Spec.Location)
	resource := utils.ServiceName(parent, instance.Spec.Service.Metadata.Name)

	if r.finalizer.IsDeletionCandidate(instance) {
		if value, exists := instance.GetAnnotations()[annotationDeletion]; exists && value == "true" {
			log.Info("Delete Cloud Run Service", "service", instance.Spec.Service)
			if err := rm.Delete(resource, instance.Spec.Service); err != nil {
				return reconcile.Result{}, err
			}
		}
		r.finalizer.Remove(instance)
		return reconcile.Result{}, r.client.Update(context.TODO(), instance)
	}
	r.finalizer.Add(instance)

	if err := rm.SetIamPolicy(resource, &instance.Spec.IamPolicy); err != nil {
		return reconcile.Result{}, err
	}

	if createdService, err := rm.CreateOrUpdate(parent, &instance.Spec.Service); err != nil {
		return reconcile.Result{}, err
	} else {
		instance.Status.Url = createdService.Status.Address.Hostname
	}

	// Update Status
	if !reflect.DeepEqual(currentInstance.Status, instance.Status) {
		if err := r.client.Status().Update(context.TODO(), instance); err != nil {
			reqLogger.Error(err, "Failed to update instance status")
			return reconcile.Result{}, err
		}
	}

	reqLogger.Info("End of work", "requeue", requeueAfter)
	return reconcile.Result{RequeueAfter: requeueAfter}, nil
}
