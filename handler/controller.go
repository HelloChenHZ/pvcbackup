package handler

import (
	"fmt"
	"log"

	v1 "k8s.io/api/core/v1"
	v1informer "k8s.io/client-go/informers/core/v1"
	v1lister "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
)

const backupImageName = "backupimage:v1.0"
const recoveryImageName = "recoveryimage:v1.0"

type getPodController struct {
	PoLister v1lister.PodLister
}

var GetPodController = &getPodController{}

func (c *getPodController) InitGetPodController(podInformer v1informer.PodInformer) {
	c.PoLister = podInformer.Lister()
	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.onAdd,
		UpdateFunc: c.onUpdate,
		DeleteFunc: c.onDelete,
	})
	return
}

func (*getPodController) Run(stopCh chan struct{}) {
	<-stopCh
}

func (*getPodController) onAdd(obj interface{}) {
	pod := obj.(*v1.Pod)
	log.Println(pod.Name + " is add")
}

func (*getPodController) onUpdate(oldObj interface{}, newObj interface{}) {
	pod := newObj.(*v1.Pod)
	log.Println(pod.Name + " is update,the status is " + string(pod.Status.Phase))
}

func (*getPodController) onDelete(obj interface{}) {
	pod := obj.(*v1.Pod)
	log.Println(pod.Name + " is delete")
	for _, volume := range pod.Spec.Volumes {
		if volume.PersistentVolumeClaim != nil {
			fmt.Println("Pod Name: " + pod.GetName())
			fmt.Println("PVC Name: " + volume.PersistentVolumeClaim.ClaimName)
			fmt.Println("Node Name:", pod.Spec.NodeName)
			pvcName := volume.PersistentVolumeClaim.ClaimName
			nodeName := pod.Spec.NodeName
			dataPath := "/data"
			s3Path := pvcName + pod.Name
			// TODO: read value from backup resource
			createJob(pvcName, nodeName, dataPath, s3Path, backupImageName, "-d "+dataPath+"-s "+s3Path)
		}
	}
}
