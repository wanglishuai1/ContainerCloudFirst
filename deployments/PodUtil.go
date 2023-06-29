package deployment

import (
	"ContainerCloudFirst/lib"
	"context"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// 判断POD是否就绪
func GetPodIsReady(pod v1.Pod) bool {
	for _, condition := range pod.Status.Conditions {
		if condition.Type == "ContainersReady" && condition.Status != "True" {
			return false
		}
	}
	for _, rg := range pod.Spec.ReadinessGates {
		for _, condition := range pod.Status.Conditions {
			if condition.Type == rg.ConditionType && condition.Status != "True" {
				return false
			}
		}
	}
	return true
}
func DeletePod(ns string, podName string) error {
	return lib.K8sClient.CoreV1().Pods(ns).
		Delete(context.Background(), podName, metav1.DeleteOptions{})
}
