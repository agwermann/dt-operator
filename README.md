# dt-operator

```
make manifests && make install && kubectl apply -f config/samples/
```

```
kubectl api-resources
```

```
kubectl get twincomponent
kubectl get twininstance
kubectl get messagebrokers
```

## Description
// TODO(user): An in-depth paragraph about your project and overview of use

## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Running on the cluster
1. Install Instances of Custom Resources:

```sh
kubectl apply -f config/samples/
```

2. Build and push your image to the location specified by `IMG`:
	
```sh
make docker-build docker-push IMG=<some-registry>/dt-operator:tag
```
	
3. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/dt-operator:tag
```

### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
UnDeploy the controller to the cluster:

```sh
make undeploy
```

## Contributing
// TODO(user): Add detailed information on how you would like others to contribute to this project

### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/) 
which provides a reconcile function responsible for synchronizing resources untile the desired state is reached on the cluster 

### Test It Out
1. Install the CRDs into the cluster:

```sh
make install
```

2. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

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

## Setup Project

Create Digital Twin Operator Project:

```bash
kubebuilder init --domain digitaltwin --repo github.com/agwermann/dt-operator
```

Create Digital Twin Definition API resources:

```bash
kubebuilder create api --group dtd --version v0 --kind TwinComponent
kubebuilder create api --group dtd --version v0 --kind TwinInstance
```

Create Digital Twin Core API Resources:

```bash
kubebuilder create api --group core --version v0 --kind MessagingGateway
kubebuilder create api --group core --version v0 --kind MessageBroker
kubebuilder create api --group core --version v0 --kind EventStore
```

## Pre-Requisites

1. Configure your Kubernetes cluster. You can run the platform in [Kind](https://kind.sigs.k8s.io/) in your local computer.

```sh
kind create cluster
```

2. Deploy ScillaDB for the Event Store.

```sh

```

3. Install Knative and Istio dependencies.

- Knative Serving:

```sh
kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.8.0/serving-crds.yaml
kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.8.0/serving-core.yaml
kubectl get pods --namespace knative-serving
```

- Knative Eventing:

```sh
kubectl apply -f https://github.com/knative/eventing/releases/download/knative-v1.8.0/eventing-crds.yaml
kubectl apply -f https://github.com/knative/eventing/releases/download/knative-v1.8.0/eventing-core.yaml
kubectl get pods --namespace knative-eventing
```

4. Install Camel-k

```sh
kubectl -n default create secret docker-registry external-registry-secret --docker-username <DOCKER_USERNAME> --docker-password <DOCKER_PASSWORD> -n dtserverless
kamel install --operator-image=docker.io/apache/camel-k:1.10.3 --olm=false -n dtserverless --global --registry docker.io --organization agwermann --registry-secret external-registry-secret --force
```

## Install Digital Twin Platform

Now that you have configured the above pre-requisites, you can install the platform.

```sh

```

## Message Format

The platform uses Cloud Event specification for describing event data within the platform boundaries. The physical twins messages are expected to follow Cloud Event specification as well. In case the below format is not sent, the platform will assume the content sent is the data field in JSON format, and it will populate the rest of the mandatory fields.

```json
{
    "id": "41c9afea-02a1-48a0-ad3c-cb9faae86551",
    "type" : "com.digitaltwin.telemetry.temperature",
    "subject": "SensorData",
    "source": "Sensor-01012023",
    "time" : "2018-04-05T17:31:00Z",
    // "component": "Sensor", TBD
    "datacontenttype": "application/json",
    "data": {
        "temperature": 20
    }
}
```

- Adicionar a nível de Component

- Factory1
    - Machine1
        - Sensor1
        - Sensor2
    - Machine2
        - Sensor3
        - Sensor4

- Factory2
    - Machine3
        - Sensor1
        - Sensor2
    - Machine4
        - Sensor3
        - Sensor4
