
# Create cluster
kind create cluster --config scripts/kind-config.yaml

# Load docker images
kind load docker-image agwermann/manufactoring-service:0.1

# Install Knative and Istio dependencies
kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.8.0/serving-crds.yaml
kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.8.0/serving-core.yaml
kubectl get pods --namespace knative-serving

# Install Istio
kubectl apply -l knative.dev/crd-install=true -f https://github.com/knative/net-istio/releases/download/knative-v1.8.0/istio.yaml
kubectl apply -f https://github.com/knative/net-istio/releases/download/knative-v1.8.0/istio.yaml
kubectl apply -f https://github.com/knative/net-istio/releases/download/knative-v1.8.0/net-istio.yaml
kubectl --namespace istio-system get service istio-ingressgateway
kubectl get pods --namespace knative-serving
kubectl get pods --namespace istio-system

# Install Eventing dependencies
kubectl apply -f https://github.com/knative/eventing/releases/download/knative-v1.8.0/eventing-crds.yaml
kubectl apply -f https://github.com/knative/eventing/releases/download/knative-v1.8.0/eventing-core.yaml
kubectl get pods --namespace knative-eventing

# Create namespace
kubectl apply -f config/samples/0-dt-namespaces.yaml

# Install Camel-K
kubectl create secret docker-registry external-registry-secret --docker-username agwermann --docker-password <DOCKER_PASSWORD> -n default
kamel install --operator-image=docker.io/apache/camel-k:1.10.3 --olm=false -n default --global --registry docker.io --organization agwermann --registry-secret external-registry-secret --force

# HTTP Bin
kubectl apply -f https://github.com/istio/istio/raw/master/samples/httpbin/httpbin.yaml
