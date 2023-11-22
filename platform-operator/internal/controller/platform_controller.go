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

package controller

import (
	"context"

	// "helm.sh/helm/v3/pkg/repo"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	zerosdpv1 "github.com/audacioustux/zerosdp/platform-operator/api/v1"
	_ctrl "github.com/audacioustux/zerosdp/platform-operator/pkg/controller"

	"github.com/go-logr/logr"
)

const (
	platformFinalizer = "platform.zerosdp.alo.dev/finalizer"
)

// PlatformReconciler reconciles a Platform object
type PlatformReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=zerosdp.alo.dev,resources=platforms,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=zerosdp.alo.dev,resources=platforms/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=zerosdp.alo.dev,resources=platforms/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch

func (r *PlatformReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	platform := &zerosdpv1.Platform{}
	if err := r.Get(ctx, req.NamespacedName, platform); err != nil {
		if errors.IsNotFound(err) {
			log.Info("Platform not found")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get platform")
		return ctrl.Result{}, err
	}

	return r.reconcile(ctx, platform, log)
}

func (r PlatformReconciler) reconcile(ctx context.Context, platform *zerosdpv1.Platform, log logr.Logger) (ctrl.Result, error) {
	log.V(1).Info("Reconciling")
	meta.SetStatusCondition(&platform.Status.Conditions, metav1.Condition{
		Type:    string(zerosdpv1.Reconciling),
		Status:  metav1.ConditionTrue,
		Reason:  "Reconciling",
		Message: "Starting Reconciliation",
	})
	if err := r.Update(ctx, platform); err != nil {
		log.Error(err, "Failed to update platform")
		return ctrl.Result{}, err
	}

	// Initialize status
	if platform.Status.Conditions == nil || len(platform.Status.Conditions) == 0 {
		log.Info("Initializing status")
		meta.SetStatusCondition(&platform.Status.Conditions, metav1.Condition{
			Type:    string(zerosdpv1.Ready),
			Status:  metav1.ConditionUnknown,
			Reason:  "Reconciling",
			Message: "Starting Reconciliation",
		})
	}
	r.Status().Update(ctx, platform)

	// Check if status is unknown
	if meta.IsStatusConditionPresentAndEqual(platform.Status.Conditions, string(zerosdpv1.Ready), metav1.ConditionUnknown) {
		log.Info("Status is unknown")

		return _ctrl.ShortRequeue(), nil
	}

	// Add finalizer
	if !controllerutil.ContainsFinalizer(platform, platformFinalizer) {
		log.Info("Adding finalizer")
		if ok := controllerutil.AddFinalizer(platform, platformFinalizer); !ok {
			log.Info("Finalizer list not updated")
			return ctrl.Result{Requeue: true}, nil
		}

		if err := r.Update(ctx, platform); err != nil {
			log.Error(err, "Failed to update platform")
			return ctrl.Result{}, err
		}
	}

	// Uninstall
	if platform.GetDeletionTimestamp() != nil {
		if controllerutil.ContainsFinalizer(platform, platformFinalizer) {
			log.Info("Performing finalizer operations")

			meta.SetStatusCondition(&platform.Status.Conditions, metav1.Condition{
				Type:    string(zerosdpv1.Uninstalling),
				Status:  metav1.ConditionTrue,
				Reason:  "Finalizing",
				Message: "Starting Finalization",
			})

			if err := r.Status().Update(ctx, platform); err != nil {
				log.Error(err, "Failed to update status")
				return ctrl.Result{}, err
			}

			// Execute finalizer
			if err := r.finalize(ctx, platform, log); err != nil {
				log.Error(err, "Failed to finalize")
				return ctrl.Result{}, err
			}

			// Remove finalizer
			if ok := controllerutil.RemoveFinalizer(platform, platformFinalizer); !ok {
				log.Info("Finalizer list not updated")
				return ctrl.Result{Requeue: true}, nil
			}

			if err := r.Update(ctx, platform); err != nil {
				log.Error(err, "Failed to update platform")
				return ctrl.Result{}, err
			}

			return ctrl.Result{}, nil
		}
	}

	meta.SetStatusCondition(&platform.Status.Conditions, metav1.Condition{
		Type:    string(zerosdpv1.Reconciling),
		Status:  metav1.ConditionFalse,
		Reason:  "Reconciled",
		Message: "Finished Reconciliation",
	})
	if err := r.Update(ctx, platform); err != nil {
		log.Error(err, "Failed to update platform")
		return ctrl.Result{}, err
	}
	log.V(1).Info("Reconciled")

	return ctrl.Result{}, nil
}

func (r PlatformReconciler) finalize(ctx context.Context, platform *zerosdpv1.Platform, log logr.Logger) error {
	log.V(1).Info("Finalizing")
	defer log.V(1).Info("Finalized")

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PlatformReconciler) SetupWithManager(mgr ctrl.Manager, maxConcurrentReconciles int) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&zerosdpv1.Platform{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: maxConcurrentReconciles}).
		Complete(r)
}
