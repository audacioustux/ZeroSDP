package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	zerosdpv1alpha1 "github.com/audacioustux/zerosdp/api/v1alpha1"
	"github.com/audacioustux/zerosdp/pkg/logging"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PlatformReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logging.Logger
}

// +kubebuilder:rbac:groups=zerosdp.alo.dev,resources=platforms,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=zerosdp.alo.dev,resources=platforms/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=zerosdp.alo.dev,resources=platforms/finalizers,verbs=update

func (r *PlatformReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("platform", req.NamespacedName)

	// Fetch the Platform instance
	platform := &zerosdpv1alpha1.Platform{}
	if err := r.Get(ctx, req.NamespacedName, platform); err != nil {
		log.Info("unable to fetch Platform object", "error", err)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Update status
	if err := r.updateStatus(ctx, platform, log); err != nil {
		log.Info("unable to update Platform status", "error", err)
		return ctrl.Result{}, err
	}

	// Reconcile actual state with desired state
	return r.reconcile(ctx, platform, log)
}

func (r PlatformReconciler) reconcile(ctx context.Context, platform *zerosdpv1alpha1.Platform, log logging.Logger) (ctrl.Result, error) {
	log = log.WithValues("context", "reconcile")

	log.Debug("Reconciling")
	defer log.Debug("Reconciled")

	// Log current status
	log.Debug("Current status", "status", platform.Status)

	// Check if status is unknown
	if meta.IsStatusConditionPresentAndEqual(platform.Status.Conditions, string(zerosdpv1alpha1.Ready), metav1.ConditionUnknown) {

		return ctrl.Result{Requeue: true}, nil
	}

	return ctrl.Result{}, nil
}

func ShouldRequeue(r ctrl.Result) bool {
	return r.Requeue || r.RequeueAfter > 0
}

func (r PlatformReconciler) updateStatus(ctx context.Context, platform *zerosdpv1alpha1.Platform, log logging.Logger) error {
	log = log.WithValues("context", "updateStatus")

	log.Debug("Updating status")
	defer log.Debug("Updated status")

	// Initialize status
	if platform.Status.Conditions == nil || len(platform.Status.Conditions) == 0 {
		log.Info("Initializing status")
		meta.SetStatusCondition(&platform.Status.Conditions, metav1.Condition{
			Type:    string(zerosdpv1alpha1.Ready),
			Status:  metav1.ConditionUnknown,
			Reason:  "Reconciling",
			Message: "Starting Reconciliation",
		})
	}

	return r.Status().Update(ctx, platform)
}

func (r *PlatformReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&zerosdpv1alpha1.Platform{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: 3}).
		Complete(r)
}
