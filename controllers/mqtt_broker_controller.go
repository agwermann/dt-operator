package controllers

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const BROKER_CONFIG_MAP_NAME = "mqtt-broker-config"
const BROKER_DEPLOYMENT_NAME = "mqtt-broker-deployment"
const BROKER_SERVICE_NAME = "mqtt-broker-service"
const BROKER_NAMESPACE = "mqtt"

var BROKER_CONFIG_MAP_KEY = types.NamespacedName{
	Name:      BROKER_CONFIG_MAP_NAME,
	Namespace: BROKER_NAMESPACE,
}

var BROKER_DEPLOYMENT_KEY = types.NamespacedName{
	Name:      BROKER_DEPLOYMENT_NAME,
	Namespace: BROKER_NAMESPACE,
}

var BROKER_SERVICE_KEY = types.NamespacedName{
	Name:      BROKER_SERVICE_NAME,
	Namespace: BROKER_NAMESPACE,
}

func buildLabels(appLabel string) map[string]string {
	return map[string]string{
		"app": appLabel,
	}
}

func (r *TwinServiceReconciler) deleteBrokerDeployment(ctx context.Context, request reconcile.Request) (*reconcile.Result, error) {
	logger := log.FromContext(ctx).WithValues("TwinService", request.NamespacedName)

	err := r.Delete(context.TODO(), r.buildBrokerConfigMapDefinition())

	if err != nil {
		logger.Error(err, "Error while deleting broker config map")
	}

	err = r.Delete(context.TODO(), r.buildBrokerDeploymentDefinition())

	if err != nil {
		logger.Error(err, "Error while deleting broker deployment")
	}

	err = r.Delete(context.TODO(), r.buildBrokerServiceDefinition())

	if err != nil {
		logger.Error(err, "Error while deleting broker service")
	}

	return nil, nil
}

func (r *TwinServiceReconciler) applyBrokerDeployment(ctx context.Context, request reconcile.Request) (*reconcile.Result, error) {
	logger := log.FromContext(ctx).WithValues("TwinService", request.NamespacedName)

	logger.Info("Creating MQTT Broker")

	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: BROKER_NAMESPACE,
		},
	}

	err := r.Create(ctx, namespace)

	// TODO FIX
	if err != nil {
		logger.Info(`Namespace already exists: ` + BROKER_NAMESPACE)
	}

	// Create ConfigMap, if not exists
	configMap := &v1.ConfigMap{}

	err = r.Get(context.TODO(), BROKER_CONFIG_MAP_KEY, configMap)

	if err != nil && errors.IsNotFound(err) {
		err := r.Create(context.TODO(), r.buildBrokerConfigMapDefinition())

		if err != nil {
			logger.Error(err, `Error while creating broker config map: `+BROKER_CONFIG_MAP_NAME)
			return &reconcile.Result{}, err
		}
	}

	// Create Deployment, if not exists
	deployment := &appsv1.Deployment{}

	err = r.Get(context.TODO(), BROKER_DEPLOYMENT_KEY, deployment)

	if err != nil && errors.IsNotFound(err) {
		err := r.Create(context.TODO(), r.buildBrokerDeploymentDefinition())

		if err != nil {
			logger.Error(err, `Error while creating broker deployment: `+BROKER_CONFIG_MAP_NAME)
			return &reconcile.Result{}, err
		}
	}

	// Create Service, if not exists
	service := &v1.Service{}

	err = r.Get(context.TODO(), BROKER_SERVICE_KEY, service)

	if err != nil && errors.IsNotFound(err) {
		err := r.Create(context.TODO(), r.buildBrokerServiceDefinition())

		if err != nil {
			logger.Error(err, `Error while creating broker service: `+BROKER_SERVICE_NAME)
			return &reconcile.Result{}, err
		}
	}

	return nil, nil
}

func (r *TwinServiceReconciler) buildBrokerDeploymentDefinition() *appsv1.Deployment {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      BROKER_DEPLOYMENT_NAME,
			Namespace: BROKER_NAMESPACE,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: buildLabels(BROKER_DEPLOYMENT_NAME),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: buildLabels(BROKER_DEPLOYMENT_NAME),
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "mosquitto",
						Image: "eclipse-mosquitto:2.0",
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{}, // TODO
							Limits:   corev1.ResourceList{}, // TODO
						},
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 1883,
							},
						},
						// VolumeMounts: []corev1.VolumeMount{
						// 	{
						// 		Name:      "mosquitto-config",
						// 		MountPath: "/mosquitto/config/mosquitto.conf",
						// 		SubPath:   "mosquitto.conf",
						// 	},
						// },
					}},
					// Volumes: []corev1.Volume{
					// 	{
					// 		Name: "mosquitto-config",
					// 	},
					// },
				},
			},
		},
	}
	return deployment
}

func (r *TwinServiceReconciler) buildBrokerConfigMapDefinition() *v1.ConfigMap {
	configmap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      BROKER_CONFIG_MAP_NAME,
			Namespace: BROKER_NAMESPACE,
		},
		Data: map[string]string{
			"mosquitto.conf": `|-
			# Ip/hostname to listen to.
			# If not given, will listen on all interfaces
			#bind_address
		
			# Port to use for the default listener.
			port 1883
		
			# Allow anonymous users to connect?
			# If not, the password file should be created
			allow_anonymous true
		
			# The password file.
			# Use the "mosquitto_passwd" utility.
			# If TLS is not compiled, plaintext "username:password" lines bay be used
			# password_file /mosquitto/config/passwd
			`,
		},
	}
	return configmap
}

func (r *TwinServiceReconciler) buildBrokerServiceDefinition() *v1.Service {
	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      BROKER_SERVICE_NAME,
			Namespace: BROKER_NAMESPACE,
		},
		Spec: v1.ServiceSpec{
			Selector: buildLabels(BROKER_SERVICE_NAME),
			Ports: []v1.ServicePort{
				{
					Port: 1883,
					TargetPort: intstr.IntOrString{
						IntVal: 1883,
					},
				},
			},
		},
	}
	return service
}
