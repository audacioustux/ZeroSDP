package utils

import (
	"time"

	"k8s.io/apimachinery/pkg/util/rand"
	ctrl "sigs.k8s.io/controller-runtime"
)

func ShouldRequeue(r ctrl.Result) bool {
	return r.Requeue || r.RequeueAfter > 0
}

func NewRequeueWithDelay(min int, max int, units time.Duration) ctrl.Result {
	seconds := rand.IntnRange(min, max)
	return ctrl.Result{RequeueAfter: time.Duration(seconds) * units}
}

func ShortRequeue() ctrl.Result {
	return NewRequeueWithDelay(2, 3, time.Second)
}

func LongRequeue() ctrl.Result {
	return NewRequeueWithDelay(2, 3, time.Minute)
}
