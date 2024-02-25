package handler

import (
	"context"
	"fmt"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"strings"
	"time"
)

type Job struct {
	JobName           string            `form:"jobName" json:"jobName" binding:"required"`
	ContainerImage    string            `form:"containerImage" json:"containerImage" binding:"required"`
	Args              string            `form:"args" json:"args" binding:"required"`
	NodeAffinityValue string            `form:"nodeAffinityValue" json:"nodeAffinityValue" binding:"required"`
	CPUtLimit         string            `form:"cpuLimit" json:"cpuLimit"`
	CPURequest        string            `form:"cpuRequest" json:"cpuRequest"`
	MemLimit          string            `form:"memLimit" json:"memLimit"`
	MemRequest        string            `form:"memRequest" json:"memRequest"`
	Label             map[string]string `form:"label" json:"label"`
}

func CreateJob(pvcName, nodeName, dataPath, s3Path, containerImage, args string) {
	t := time.Now().UTC()
	jobName := pvcName + t.Format("2024-02-25-21")
	json := Job{
		JobName:           jobName,
		ContainerImage:    containerImage,
		Args:              args,
		NodeAffinityValue: nodeName,
		CPUtLimit:         "1000m",
		MemLimit:          "800Mi",
		CPURequest:        "100m",
		MemRequest:        "400Mi",
	}

	fmt.Printf("Args: %s %s %s\n", json.JobName, json.ContainerImage, json.Args)

	// create job
	jobs := KubernetesClientset.AppsV1().Deployments("quant-job")

	var replicas int32 = 1

	json.Label["job-name"] = json.JobName

	jobSpec := &appv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      json.JobName,
			Namespace: "default",
			Labels:    json.Label,
		},
		Spec: appv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: json.Label,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   json.JobName,
					Labels: json.Label,
				},
				Spec: v1.PodSpec{
					Affinity: &v1.Affinity{
						NodeAffinity: &v1.NodeAffinity{
							RequiredDuringSchedulingIgnoredDuringExecution: &v1.NodeSelector{
								NodeSelectorTerms: []v1.NodeSelectorTerm{
									{
										MatchExpressions: []v1.NodeSelectorRequirement{
											{
												Key:      "type",
												Operator: v1.NodeSelectorOpIn,
												Values:   []string{json.NodeAffinityValue},
											},
										},
									},
								},
							},
						},
					},
					Containers: []v1.Container{
						{
							Name:            json.JobName,
							Image:           json.ContainerImage,
							Args:            strings.Split(json.Args, " "),
							ImagePullPolicy: "Always",
							Resources: v1.ResourceRequirements{
								Limits: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse(json.CPUtLimit),
									v1.ResourceMemory: resource.MustParse(json.MemLimit),
								},
								Requests: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse(json.CPURequest),
									v1.ResourceMemory: resource.MustParse(json.MemRequest),
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "pvcData",
									MountPath: "/data",
								},
								{
									Name:      "s3Data",
									MountPath: "/s3data",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "pvcData",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: pvcName, // PVC name
								},
							},
						},
						{
							Name: "s3Data",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: "s3data", // PVC name
								},
							},
						},
					},
					RestartPolicy:    v1.RestartPolicyAlways,
					ImagePullSecrets: []v1.LocalObjectReference{{Name: "gcp-gitlab-registry"}},
				},
			},
		},
	}

	_, err := jobs.Create(context.TODO(), jobSpec, metav1.CreateOptions{})
	if err != nil {
		log.Println("Failed to create K8s job." + err.Error())
		return
	}

	//print job details
	log.Println("Created K8s job successfully")

	return
}
