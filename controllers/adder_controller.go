/*
Copyright 2021 Sunil Thaha.

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

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	opsv1beta1 "github.com/sthaha/opdemo/api/v1beta1"
)

// AdderReconciler reconciles a Adder object
type AdderReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=ops.cee.org,resources=adders,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ops.cee.org,resources=adders/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=ops.cee.org,resources=adders/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Adder object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *AdderReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("adder", req.NamespacedName)

	// your logic here
	adder := &opsv1beta1.Adder{}
	err := r.Get(ctx, req.NamespacedName, adder)
	if err != nil {
		log.Error(err, "failed to get adder")
		return ctrl.Result{}, err
	}

	inputs := adder.Spec.Inputs
	result := adder.Status.Result
	log.Info("Adder", "name", adder.Name, "inputs", inputs, "result", result)

	var sum int32
	for _, v := range inputs {
		sum += v
	}

	// everytthing is as expected
	if sum == result {
		log.Info("Adder", "name", adder.Name, "sum", sum, "result", result)
		return ctrl.Result{}, err
	}

	adder.Status.Result = sum
	log.Info("Adder Updated", "name", adder.Name, "sum", sum, "result", adder.Status.Result)

	err = r.Status().Update(ctx, adder)
	if err != nil {
		log.Error(err, "Failed to update Memcached status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AdderReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&opsv1beta1.Adder{}).
		Complete(r)
}
