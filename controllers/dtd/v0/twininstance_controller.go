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
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	kserving "knative.dev/serving/pkg/apis/serving/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	dtdv0 "github.com/agwermann/dt-operator/apis/dtd/v0"
)

// TwinInstanceReconciler reconciles a TwinInstance object
type TwinInstanceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=dtd.digitaltwin,resources=twininstances,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=dtd.digitaltwin,resources=twininstances/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=dtd.digitaltwin,resources=twininstances/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the TwinInstance object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *TwinInstanceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx).WithValues("TwinInstance", req.NamespacedName)

	logger.Info("New TwinInstance identified")

	twinInstance := &dtdv0.TwinInstance{}
	err := r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, twinInstance)

	if err != nil {
		logger.Error(err, "Not possible to get Twin Instance"+twinInstance.Name)
	}

	kService := &kserving.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "serving.knative.dev/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      twinInstance.Name,
			Namespace: twinInstance.Namespace,
		},
		Spec: kserving.ServiceSpec{
			ConfigurationSpec: kserving.ConfigurationSpec{
				Template: kserving.RevisionTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Name: twinInstance.Name + "-01",
					},
					Spec: kserving.RevisionSpec{
						PodSpec: twinInstance.Spec.Template.Spec,
					},
				},
			},
		},
	}

	err = r.Create(ctx, kService, &client.CreateOptions{})

	if err != nil {
		logger.Error(err, fmt.Sprintf("Error while reconcile"))
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TwinInstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dtdv0.TwinInstance{}).
		Complete(r)
}
