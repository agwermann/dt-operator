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

package core

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	corev0 "github.com/agwermann/dt-operator/apis/core/v0"
	v0 "github.com/agwermann/dt-operator/apis/core/v0"
	"github.com/agwermann/dt-operator/controllers/core/broker"
)

// MessageBrokerReconciler reconciles a MessageBroker object
type MessageBrokerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=core.digitaltwin,resources=messagebrokers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core.digitaltwin,resources=messagebrokers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core.digitaltwin,resources=messagebrokers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MessageBroker object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *MessageBrokerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO: Creation and deletion work fine.
	// TODO: Need to understand if there is another way to do, if not we must create a MQTTMessageBroker resource instead of keep it generic
	logger := log.FromContext(ctx).WithValues("MessageBroker", req.NamespacedName)

	messageBroker := &corev0.MessageBroker{}
	err := r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, messageBroker)

	broker := broker.NewMqttMessageBroker(logger)

	// Delete scenario
	if err != nil {
		if errors.IsNotFound(err) {
			err = r.Delete(context.TODO(), broker.GetBrokerDeployment(req.Name, req.Namespace), &client.DeleteOptions{})

			if err != nil {
				logger.Error(err, fmt.Sprintf("Error while deleting broker %s deployment", req.Name))
			}

			err = r.Delete(context.TODO(), broker.GetBrokerConfigMap(req.Name, req.Namespace), &client.DeleteOptions{})

			if err != nil {
				logger.Error(err, fmt.Sprintf("Error while deleting broker %s config map", req.Name))
			}

			err = r.Delete(context.TODO(), broker.GetBrokerService(req.Name, req.Namespace), &client.DeleteOptions{})

			if err != nil {
				logger.Error(err, fmt.Sprintf("Error while deleting broker %s service", req.Name))
			}
		}
		return ctrl.Result{}, nil
	}

	// Create / Update scenarios
	logger.Info(fmt.Sprintf("Creating Message Broker of type %s", messageBroker.Spec.Type))

	// Check if there is already a mqtt broker (it can only have one)
	err = r.Create(context.TODO(), broker.GetBrokerDeployment(req.Name, req.Namespace), &client.CreateOptions{})
	err = r.Create(context.TODO(), broker.GetBrokerConfigMap(req.Name, req.Namespace), &client.CreateOptions{})
	err = r.Create(context.TODO(), broker.GetBrokerService(req.Name, req.Namespace), &client.CreateOptions{})

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MessageBrokerReconciler) ResolveBrokerType(ctx context.Context, req ctrl.Request, brokerType v0.MessageBrokerType) (broker.MessageBroker, error) {
	var messageBroker broker.MessageBroker

	logger := log.FromContext(ctx).WithValues("MessageBroker", req.NamespacedName)

	switch brokerType {
	case v0.MESSAGE_BROKER_MQTT:
		messageBroker = broker.NewMqttMessageBroker(logger)
	case v0.MESSAGE_BROKER_AMQP:
		errorMessage := fmt.Sprintf("Broker Type %s still not supported", v0.MESSAGE_BROKER_AMQP)
		return messageBroker, errors.NewBadRequest(errorMessage)
	case v0.MESSAGE_BROKER_KAFKA:
		errorMessage := fmt.Sprintf("Broker Type %s still not supported", v0.MESSAGE_BROKER_KAFKA)
		return messageBroker, errors.NewBadRequest(errorMessage)
	default:
		errorMessage := fmt.Sprintf("Message broker type %s is invalid: it must be one of the following: %s, %s, %s", brokerType, v0.MESSAGE_BROKER_AMQP, v0.MESSAGE_BROKER_KAFKA, v0.MESSAGE_BROKER_MQTT)
		return messageBroker, errors.NewBadRequest(errorMessage)
	}

	return messageBroker, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MessageBrokerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev0.MessageBroker{}).
		Complete(r)
}
