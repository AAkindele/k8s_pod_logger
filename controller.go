package main

import (
	"fmt"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type Controller struct {
	Queue    workqueue.RateLimitingInterface
	Informer cache.SharedInformer
	Options  controllerOptions
}

type controllerOptions struct {
	PodFilter string
}

func NewController(
	queue workqueue.RateLimitingInterface,
	informer cache.SharedInformer,
	options controllerOptions,
) *Controller {
	return &Controller{
		Queue:    queue,
		Informer: informer,
		Options:  options,
	}
}

// entry point for the controller
func (controller *Controller) Run(stopCh chan struct{}) {
	defer runtime.HandleCrash()
	defer controller.Queue.ShutDown()

	fmt.Println("starting controller")
	fmt.Printf("filtering for pods with the label 'k8s-pod-logger=%s'\n", controller.Options.PodFilter)

	// start listening for changes on the provided listener
	go controller.Informer.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh, controller.Informer.HasSynced) {
		fmt.Println("Timeout waiting for cache to sync.")
		return
	}

	go controller.runWorker()

	<-stopCh
	fmt.Println("stopping controller")
}

func (controller *Controller) runWorker() {
	for controller.processNextWorkItem() {
	}
}

func (controller *Controller) processNextWorkItem() bool {
	// block until an item is returned, or queue tells us to shutdown
	obj, shutdown := controller.Queue.Get()
	if shutdown {
		return false
	}

	err := func(obj interface{}) error {
		defer controller.Queue.Done(obj)
		logEvent, ok := obj.(LogEvent)

		if !ok {
			// requeue failed items with provided backoff strategy
			controller.Queue.AddRateLimited(obj)
			return fmt.Errorf("error processing %q: requeuing", logEvent.ResourceName)
		}

		fmt.Printf("[%s] %s\n", logEvent.Action, logEvent.ResourceName)

		controller.Queue.Forget(obj)
		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
	}

	return true
}
