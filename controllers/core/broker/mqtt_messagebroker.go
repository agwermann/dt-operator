package broker

import (
	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const MQTT_BROKER_CONFIG_MAP_SUFFIX = "-config"
const MQTT_BROKER_DEPLOYMENT_SUFFIX = "-deployment"
const MQTT_BROKER_SERVICE_SUFFIX = "-service"
const MQTT_BROKER_NAMESPACE = "mqtt"

func buildLabels(appLabel string) map[string]string {
	return map[string]string{
		"app": appLabel,
	}
}

type MessageBroker interface {
	GetBrokerDeployment(name string) *appsv1.Deployment
	GetBrokerConfigMap(name string) *v1.ConfigMap
	GetBrokerService(name string) *v1.Service
}

func NewMqttMessageBroker(logger logr.Logger) MessageBroker {
	return &mqttMessageBroker{
		logger: logger,
	}
}

type mqttMessageBroker struct {
	logger logr.Logger
}

func (b *mqttMessageBroker) CreateBroker() error {
	b.logger.Info("Creating MQTT Broker")

	// namespace := &corev1.Namespace{
	// 	ObjectMeta: metav1.ObjectMeta{
	// 		Name: BROKER_NAMESPACE,
	// 	},
	// }

	// err := r.Create(ctx, namespace)

	// if err != nil {
	// 	logger.Info(`Namespace already exists: ` + BROKER_NAMESPACE)
	// }

	// // Create ConfigMap, if not exists
	// configMap := &v1.ConfigMap{}

	// err = r.Get(context.TODO(), BROKER_CONFIG_MAP_KEY, configMap)

	// if err != nil && errors.IsNotFound(err) {
	// 	err := r.Create(context.TODO(), r.buildBrokerConfigMapDefinition())

	// 	if err != nil {
	// 		logger.Error(err, `Error while creating broker config map: `+BROKER_CONFIG_MAP_NAME)
	// 		return &reconcile.Result{}, err
	// 	}
	// }

	// // Create Deployment, if not exists
	// deployment := &appsv1.Deployment{}

	// err = r.Get(context.TODO(), BROKER_DEPLOYMENT_KEY, deployment)

	// if err != nil && errors.IsNotFound(err) {
	// 	err := r.Create(context.TODO(), r.buildBrokerDeploymentDefinition())

	// 	if err != nil {
	// 		logger.Error(err, `Error while creating broker deployment: `+BROKER_CONFIG_MAP_NAME)
	// 		return &reconcile.Result{}, err
	// 	}
	// }

	// // Create Service, if not exists
	// service := &v1.Service{}

	// err = r.Get(context.TODO(), BROKER_SERVICE_KEY, service)

	// if err != nil && errors.IsNotFound(err) {
	// 	err := r.Create(context.TODO(), r.buildBrokerServiceDefinition())

	// 	if err != nil {
	// 		logger.Error(err, `Error while creating broker service: `+BROKER_SERVICE_NAME)
	// 		return &reconcile.Result{}, err
	// 	}
	// }

	// return nil, nil
	return nil
}

func (b *mqttMessageBroker) DeleteBroker() error {
	// logger := log.FromContext(ctx).WithValues("TwinService", request.NamespacedName)

	// err := r.Delete(context.TODO(), r.buildBrokerConfigMapDefinition())

	// if err != nil {
	// 	logger.Error(err, "Error while deleting broker config map")
	// }

	// err = r.Delete(context.TODO(), r.buildBrokerDeploymentDefinition())

	// if err != nil {
	// 	logger.Error(err, "Error while deleting broker deployment")
	// }

	// err = r.Delete(context.TODO(), r.buildBrokerServiceDefinition())

	// if err != nil {
	// 	logger.Error(err, "Error while deleting broker service")
	// }

	// return nil, nil
	return nil
}

func (b *mqttMessageBroker) DoesBrokerExists() error {
	return nil
}

func (b *mqttMessageBroker) GetBrokerDeployment(name string) *appsv1.Deployment {

	deploymentName := name + MQTT_BROKER_DEPLOYMENT_SUFFIX

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Namespace: MQTT_BROKER_NAMESPACE,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: buildLabels(deploymentName),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: buildLabels(deploymentName),
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

func (b *mqttMessageBroker) GetBrokerConfigMap(name string) *v1.ConfigMap {

	configName := name + MQTT_BROKER_CONFIG_MAP_SUFFIX

	configmap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configName,
			Namespace: MQTT_BROKER_NAMESPACE,
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

func (b *mqttMessageBroker) GetBrokerService(name string) *v1.Service {
	serviceName := name + MQTT_BROKER_SERVICE_SUFFIX
	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: MQTT_BROKER_NAMESPACE,
		},
		Spec: v1.ServiceSpec{
			Selector: buildLabels(serviceName),
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
