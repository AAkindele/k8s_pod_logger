package main

import (
	"flag"
	"fmt"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type LogEvent struct {
	ResourceName string
	Action       string
}

func main() {
	// parse command line flags
	podFilter := flag.String("f", "", "limit logging to pods that match this filter")
	flag.Parse()

	options := controllerOptions{
		PodFilter: *podFilter,
	}
	fmt.Println("podFilter", *podFilter)

	// create a k8s client
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// create informer that will be used to watch for pods changes
	factory := informers.NewSharedInformerFactory(clientset, 0)
	podsInformer := factory.Core().V1().Pods()
	informer := podsInformer.Informer()

	// the pods we care about will be pushed here
	podQueue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	// on changes to pods, perform actions
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			updateHandler(obj, "ADD", cache.MetaNamespaceKeyFunc, podQueue, options)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			updateHandler(newObj, "UPDATE", cache.MetaNamespaceKeyFunc, podQueue, options)
		},
		DeleteFunc: func(obj interface{}) {
			updateHandler(obj, "DELETE", cache.DeletionHandlingMetaNamespaceKeyFunc, podQueue, options)
		},
	})

	// start the controller
	stopCh := make(chan struct{})
	defer close(stopCh)
	controller := NewController(podQueue, informer, options)
	controller.Run(stopCh)
}

// process the change event for pod. queue the name if it has required label.
type namespaceFuncType = func(interface{}) (string, error)

func updateHandler(obj interface{}, action string, namespaceFunc namespaceFuncType, queue workqueue.RateLimitingInterface, options controllerOptions) {
	namespacedName, _ := namespaceFunc(obj)
	pod := obj.(*coreV1.Pod)
	labelValue := pod.Labels["k8s-pod-logger"]

	if (options.PodFilter == "") || (options.PodFilter == labelValue) {
		fmt.Printf("[QUEUE] %q. Pod has the required label value. Adding to work queue.\n", namespacedName)
		queue.Add(LogEvent{
			ResourceName: namespacedName,
			Action:       action,
		})
	} else {
		fmt.Printf("[SKIP] %q. Pod does not have the required label value. Skipping.\n", namespacedName)
	}
}
