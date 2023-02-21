/*
Copyright 2022.

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

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	dtdlv0 "github.com/agwermann/dt-operator/apis/v0"
	v0 "github.com/agwermann/dt-operator/apis/v0"
)

// TwinServiceReconciler reconciles a TwinService object
type TwinServiceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=dtdl.digitaltwin,resources=twinservices,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=dtdl.digitaltwin,resources=twinservices/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=dtdl.digitaltwin,resources=twinservices/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the TwinService object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *TwinServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("Reconciling twin service")

	twinService := &v0.TwinService{}
	r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, twinService)

	if twinService.Spec.DataSource == "mqtt" {
		_, err := r.applyBrokerDeployment(ctx, req)

		if err != nil {
			logger.Error(err, "Error while creating broker")
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TwinServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dtdlv0.TwinService{}).
		Complete(r)
}
