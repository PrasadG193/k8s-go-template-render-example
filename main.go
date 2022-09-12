package main

import (
	"bytes"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/printers"
)

// resolveGoTemplate resolves go template value from the k8s resource object
func resolveGoTemplate(obj runtime.Object, goTemplateStr string) (string, error) {
	var buff bytes.Buffer
	jp, err := printers.NewGoTemplatePrinter([]byte(goTemplateStr))
	if err != nil {
		return "", nil
	}
	err = jp.PrintObj(obj, &buff)
	return buff.String(), err
}

func printGoTemplateValues(obj runtime.Object, goTemplateStr string) {
	value, err := resolveGoTemplate(obj, goTemplateStr)
	if err != nil {
		panic(err)
	}
	fmt.Println(goTemplateStr, "--->", value)
}

func main() {
	// Tests
	printGoTemplateValues(getDeploy(), "{{ .spec }}")
	printGoTemplateValues(getDeploy(), "{{ (index .spec.template.spec.containers 0).image }}")
	printGoTemplateValues(getDeploy(), "{{ .spec.replicas }}")
	printGoTemplateValues(getDeploy(), "{{ $available := false }}{{ range $condition := $.status.conditions }}{{ if and (eq .type \"Available\") (eq .status \"True\")  }}{{ $available = true }}{{ end }}{{ end }}{{ $available }}")
}

func getDeploy() *appsv1.Deployment {
	replicas := int32(2)
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "demo",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "demo",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						corev1.Container{
							Name:  "web",
							Image: "nginx:1.12",
							Ports: []corev1.ContainerPort{
								corev1.ContainerPort{
									Name:          "http",
									HostPort:      0,
									ContainerPort: 80,
									Protocol:      corev1.Protocol("TCP"),
								},
							},
							Resources:       corev1.ResourceRequirements{},
							ImagePullPolicy: corev1.PullPolicy("IfNotPresent"),
						},
					},
				},
			},
			Strategy:        appsv1.DeploymentStrategy{},
			MinReadySeconds: 0,
		},
		Status: appsv1.DeploymentStatus{
			AvailableReplicas: 1,
			Conditions: []appsv1.DeploymentCondition{
				{
					LastTransitionTime: metav1.Now(),
					LastUpdateTime:     metav1.Now(),
					Message:            "Deployment has minimum availability.",
					Reason:             "MinimumReplicasAvailable",
					Status:             "True",
					Type:               "Available",
				},
				{
					LastTransitionTime: metav1.Now(),
					LastUpdateTime:     metav1.Now(),
					Message:            "ReplicaSet test-deployment-xxxxx has successfully progressed.",
					Reason:             "NewReplicaSetAvailable",
					Status:             "True",
					Type:               "Progressing",
				},
			},
		},
	}
}
