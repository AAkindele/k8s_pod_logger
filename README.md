A simple Kubernetes controller that logs the lifecycle of pods. This is a learning exercise to get more experience with writing kubernetes controllers.
Kubernetes already has built-in events that will track a lot of things including the lifecycle of pods.

Setup instructions assume you already have access to a Kubernetes cluster. Minikube is an option if an existing cluster is not available. https://kubernetes.io/docs/tasks/tools/install-minikube/

Clone the repository

Install the controller in the cluster. The controller in this repo filter pods by default. The filter argument can be removed.
```yaml
args:
- -f=hello-world
```
```
kubectl apply -f deployment.yml
```

View the pod logs
```
kubectl get pods --namespace k8s-pod-logger
// there should be at least one pod running
kubectl logs [INSERT POD NAME HERE] --namespace k8s-pod-logger --tail 20 --follow
```

Create some other pods. There are examples in this repo.
```
kubectl apply -f examples/pods-with-label.yml
kubectl apply -f examples/pods-without-label.yml
```

The controller logs will look something like this
```
[SKIP] "k8s-pod-logger/k8s-pod-logger-78c848978f-54l4x". Pod does not have the required label value. Skipping.
[QUEUE] "k8s-pod-logger/example-with-label-7978688659-flj8j". Pod has the required label value. Adding to work queue.
[ADD] k8s-pod-logger/example-with-label-7978688659-flj8j
[QUEUE] "k8s-pod-logger/example-with-label-7978688659-flj8j". Pod has the required label value. Adding to work queue.
[UPDATE] k8s-pod-logger/example-with-label-7978688659-flj8j
[QUEUE] "k8s-pod-logger/example-with-label-7978688659-flj8j". Pod has the required label value. Adding to work queue.
[UPDATE] k8s-pod-logger/example-with-label-7978688659-flj8j
[QUEUE] "k8s-pod-logger/example-with-label-7978688659-flj8j". Pod has the required label value. Adding to work queue.
[UPDATE] k8s-pod-logger/example-with-label-7978688659-flj8j
[SKIP] "k8s-pod-logger/example-without-label-5bf4cc9bc-zqkmm". Pod does not have the required label value. Skipping.
[SKIP] "k8s-pod-logger/example-without-label-5bf4cc9bc-zqkmm". Pod does not have the required label value. Skipping.
[SKIP] "k8s-pod-logger/example-without-label-5bf4cc9bc-zqkmm". Pod does not have the required label value. Skipping.
[SKIP] "k8s-pod-logger/example-without-label-5bf4cc9bc-zqkmm". Pod does not have the required label value. Skipping.
```
