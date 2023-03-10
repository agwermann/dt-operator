
# Create cluster
kind create cluster

# Load docker images
kind load docker-image agwermann/manufactoring-service:0.1

# Install Knative and Istio dependencies
kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.8.0/serving-crds.yaml
kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.8.0/serving-core.yaml
kubectl get pods --namespace knative-serving

# Install Eventing dependencies
kubectl apply -f https://github.com/knative/eventing/releases/download/knative-v1.8.0/eventing-crds.yaml
kubectl apply -f https://github.com/knative/eventing/releases/download/knative-v1.8.0/eventing-core.yaml
kubectl get pods --namespace knative-eventing

# Create namespace
kubectl apply -f config/samples/0-dt-core.yaml

# Install Camel-K
kubectl -n default create secret docker-registry external-registry-secret --docker-username agwermann --docker-password <PASSWORD> -n dt-core
kamel install --operator-image=docker.io/apache/camel-k:1.10.3 --olm=false -n dt-core --global --registry docker.io --organization agwermann --registry-secret external-registry-secret --force
