package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	zerosdpv1alpha1 "github.com/audacioustux/zerosdp/api/v1alpha1"
	"github.com/audacioustux/zerosdp/pkg/logging"
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
	log := r.Log.WithValues("request", req)

	log.Debug("Reconciling")
	defer log.Debug("Reconciled")

	// get the platform object
	platform := &zerosdpv1alpha1.Platform{}
	if err := r.Get(ctx, req.NamespacedName, platform); err != nil {
		log.Info("unable to fetch Platform", "error", err)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// TODO: Initiate the FSM

	// initialize the platform
	if result, err := r.initPlatform(platform); err != nil || ShouldRequeue(result) {
		log.Info("unable to initialize Platform", "error", err)
		return result, err
	}

	return ctrl.Result{}, nil
}

func ShouldRequeue(r ctrl.Result) bool {
	return r.Requeue || r.RequeueAfter > 0
}

func (r PlatformReconciler) initPlatform(platform *zerosdpv1alpha1.Platform) (ctrl.Result, error) {
	log := r.Log.WithValues("platform", platform.Name)

	log.Debug("Initializing")
	defer log.Debug("Initialized")

	return ctrl.Result{}, nil
}

func (r *PlatformReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&zerosdpv1alpha1.Platform{}).
		Complete(r)
}
