/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	installv1alpha1 "github.com/audacioustux/ZeroSDP/platform-operator/api/v1alpha1"
)

const zerosdpFinalizer = "zerosdp.audacioustux.com/finalizer"

// Definitions to manage status conditions
const (
	// typeAvailableMemcached represents the status of the Deployment reconciliation
	typeAvailableMemcached = "Available"
	// typeDegradedMemcached represents the status used when the custom resource is deleted and the finalizer operations are must to occur.
	typeDegradedMemcached = "Degraded"
)

// ZeroSDPReconciler reconciles a ZeroSDP object
type ZeroSDPReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=install.zerosdp.audacioustux.com,resources=zerosdps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=install.zerosdp.audacioustux.com,resources=zerosdps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=install.zerosdp.audacioustux.com,resources=zerosdps/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ZeroSDP object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *ZeroSDPReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the ZeroSDP instance
	zerosdp := &installv1alpha1.ZeroSDP{}
	err := r.Get(ctx, req.NamespacedName, zerosdp)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("ZeroSDP resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get ZeroSDP")
		return ctrl.Result{}, err
	}

	if zerosdp.Status.Conditions == nil || len(zerosdp.Status.Conditions) == 0 {
		meta.SetStatusCondition(&zerosdp.Status.Conditions, metav1.Condition{
			Type:    typeAvailableMemcached,
			Status:  metav1.ConditionUnknown,
			Reason:  "Reconciling",
			Message: "Starting reconciliation",
		})
		if err := r.Status().Update(ctx, zerosdp); err != nil {
			log.Error(err, "Failed to update ZeroSDP status")
			return ctrl.Result{}, err
		}

		if err := r.Get(ctx, req.NamespacedName, zerosdp); err != nil {
			log.Error(err, "Failed to re-fetch ZeroSDP")
			return ctrl.Result{}, err
		}
	}

	if !controllerutil.ContainsFinalizer(zerosdp, zerosdpFinalizer) {
		log.Info("Adding Finalizer for the ZeroSDP")
		if ok := controllerutil.AddFinalizer(zerosdp, zerosdpFinalizer); !ok {
			log.Error(err, "Failed to add finalizer into the custom resource")
			return ctrl.Result{Requeue: true}, err
		}
		if err := r.Update(ctx, zerosdp); err != nil {
			log.Error(err, "Failed to update ZeroSDP with finalizer")
			return ctrl.Result{}, err
		}
	}

	isZeroSDPMarkedToBeDeleted := zerosdp.GetDeletionTimestamp() != nil
	if isZeroSDPMarkedToBeDeleted {
		if controllerutil.ContainsFinalizer(zerosdp, zerosdpFinalizer) {
			log.Info("Performing Finalization for the ZeroSDP before delete CR")

			meta.SetStatusCondition(&zerosdp.Status.Conditions, metav1.Condition{
				Type:    typeDegradedMemcached,
				Status:  metav1.ConditionUnknown,
				Reason:  "Finalizing",
				Message: fmt.Sprintf("Performing finalizer operations for the custome resource %s", zerosdp.Name),
			})
			if err := r.Status().Update(ctx, zerosdp); err != nil {
				log.Error(err, "Failed to update ZeroSDP status")
				return ctrl.Result{}, err
			}

			r.doFinalizerOperationsForZeroSDP(zerosdp)

			if err := r.Get(ctx, req.NamespacedName, zerosdp); err != nil {
				log.Error(err, "Failed to re-fetch ZeroSDP")
				return ctrl.Result{}, err
			}

			meta.SetStatusCondition(&zerosdp.Status.Conditions, metav1.Condition{
				Type:    typeDegradedMemcached,
				Status:  metav1.ConditionTrue,
				Reason:  "Finalizing",
				Message: fmt.Sprintf("Finalizer operations for the custome resource %s completed", zerosdp.Name),
			})

			if err := r.Status().Update(ctx, zerosdp); err != nil {
				log.Error(err, "Failed to update ZeroSDP status")
				return ctrl.Result{}, err
			}

			log.Info("Removing Finalizer for the ZeroSDP after successfully perform the operations")
			if ok := controllerutil.RemoveFinalizer(zerosdp, zerosdpFinalizer); !ok {
				log.Error(err, "Failed to remove finalizer from the custom resource")
				return ctrl.Result{}, err
			}
			if err := r.Update(ctx, zerosdp); err != nil {
				log.Error(err, "Failed to update ZeroSDP with finalizer")
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// Check if the deployment already exists, if not create a new one
	found := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: zerosdp.Name, Namespace: zerosdp.Namespace}, found)
	if err != nil && apierrors.IsNotFound(err) {
		// Define a new deployment
		dep, err := r.deploymentForZeroSDP(zerosdp)
		if err != nil {
			log.Error(err, "Failed to define new Deployment resource for ZeroSDP")

			// The following implementation will update the status
			meta.SetStatusCondition(&zerosdp.Status.Conditions, metav1.Condition{Type: typeAvailableMemcached,
				Status: metav1.ConditionFalse, Reason: "Reconciling",
				Message: fmt.Sprintf("Failed to create Deployment for the custom resource (%s): (%s)", zerosdp.Name, err)})

			if err := r.Status().Update(ctx, zerosdp); err != nil {
				log.Error(err, "Failed to update ZeroSDP status")
				return ctrl.Result{}, err
			}

			return ctrl.Result{}, err
		}

		log.Info("Creating a new Deployment",
			"Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		if err = r.Create(ctx, dep); err != nil {
			log.Error(err, "Failed to create new Deployment",
				"Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return ctrl.Result{}, err
		}

		// Deployment created successfully
		// We will requeue the reconciliation so that we can ensure the state
		// and move forward for the next operations
		return ctrl.Result{RequeueAfter: time.Minute}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		// Let's return the error for the reconciliation be re-trigged again
		return ctrl.Result{}, err
	}

	size := zerosdp.Spec.Size
	if *found.Spec.Replicas != size {
		found.Spec.Replicas = &size
		if err = r.Update(ctx, found); err != nil {
			log.Error(err, "Failed to update Deployment",
				"Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)

			if err := r.Get(ctx, req.NamespacedName, zerosdp); err != nil {
				log.Error(err, "Failed to re-fetch zerosdp")
				return ctrl.Result{}, err
			}

			// The following implementation will update the status
			meta.SetStatusCondition(&zerosdp.Status.Conditions, metav1.Condition{Type: typeAvailableMemcached,
				Status: metav1.ConditionFalse, Reason: "Resizing",
				Message: fmt.Sprintf("Failed to update the size for the custom resource (%s): (%s)", zerosdp.Name, err)})

			if err := r.Status().Update(ctx, zerosdp); err != nil {
				log.Error(err, "Failed to update ZeroSDP status")
				return ctrl.Result{}, err
			}

			return ctrl.Result{}, err
		}

		// Now, that we update the size we want to requeue the reconciliation
		// so that we can ensure that we have the latest state of the resource before
		// update. Also, it will help ensure the desired state on the cluster
		return ctrl.Result{Requeue: true}, nil
	}

	// The following implementation will update the status
	meta.SetStatusCondition(&zerosdp.Status.Conditions, metav1.Condition{Type: typeAvailableMemcached,
		Status: metav1.ConditionTrue, Reason: "Reconciling",
		Message: fmt.Sprintf("Deployment for custom resource (%s) with %d replicas created successfully", zerosdp.Name, size)})

	if err := r.Status().Update(ctx, zerosdp); err != nil {
		log.Error(err, "Failed to update ZeroSDP status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *ZeroSDPReconciler) doFinalizerOperationsForZeroSDP(cr *installv1alpha1.ZeroSDP) {
	r.Recorder.Event(cr, "Warning", "Deleting",
		fmt.Sprintf("Custom Resource %s is being deleted from the namespace %s", cr.Name, cr.Namespace),
	)
}

func (r *ZeroSDPReconciler) deploymentForZeroSDP(
	zerosdp *installv1alpha1.ZeroSDP) (*appsv1.Deployment, error) {

	// nginx
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      zerosdp.Name,
			Namespace: zerosdp.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &zerosdp.Spec.Size,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": zerosdp.Name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": zerosdp.Name},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  zerosdp.Name,
							Image: "nginx:1.14.2",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: zerosdp.Spec.ContainerPort,
									Name:          "http",
								},
							},
						},
					},
				},
			},
		},
	}

	if err := ctrl.SetControllerReference(zerosdp, dep, r.Scheme); err != nil {
		return nil, err
	}
	return dep, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ZeroSDPReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&installv1alpha1.ZeroSDP{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
