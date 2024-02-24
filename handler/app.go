package handler

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
)

var (
	KubernetesClientset *kubernetes.Clientset
)

func initK8sClient() {
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
	KubernetesClientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}

	fmt.Println("Init kubernetes client successfully!")
}

func GetNodeName(pvcName string) string {
	// 获取 PVC 对应的 Pod
	podList, err := KubernetesClientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("Error retrieving pod list: ", err.Error())
		os.Exit(1)
	}

	for _, pod := range podList.Items {
		for _, volume := range pod.Spec.Volumes {
			if volume.PersistentVolumeClaim != nil && volume.PersistentVolumeClaim.ClaimName == pvcName {
				fmt.Println("Pod Name:", pod.Name)
				fmt.Println("Node Name:", pod.Spec.NodeName)
				return pod.Spec.NodeName
			}
		}
	}

	return ""
}

func Init() {
	initK8sClient()
}
