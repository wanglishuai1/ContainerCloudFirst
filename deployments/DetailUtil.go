package deployment

import (
	"ContainerCloudFirst/lib"
	"context"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetPodByDep(ns string, dep *v1.Deployment) []*Pod {
	ctx := context.Background()
	listopt := metav1.ListOptions{
		LabelSelector: GetRsLableByDeployment(dep),
	}
	list, err := lib.K8sClient.CoreV1().Pods(ns).List(ctx, listopt)
	lib.CheckError(err)
	pods := make([]*Pod, len(list.Items))
	for i, item := range list.Items {
		pods[i] = &Pod{
			Name:       item.Name,
			CreateTime: item.CreationTimestamp.Format("2006-01-02 15:04:05"),
			Node:       item.Spec.NodeName,
			Images:     GetImagesByPod(item.Spec.Containers),
		}
	}
	return pods
}
func GetDeployment(ns string, name string) *Deployment {
	ctx := context.Background()

	getopt := metav1.GetOptions{}
	dep, err := lib.K8sClient.AppsV1().Deployments(ns).Get(ctx, name, getopt)
	lib.CheckError(err)
	return &Deployment{
		Name:       dep.Name,
		Namespace:  dep.Namespace,
		Images:     GetImages(*dep),
		CreateTime: dep.CreationTimestamp.Format("2006-01-02 15:04:05"),
		Pods:       GetPodByDep(ns, dep),
		Replicas:   GetReplicas(*dep),
	}
}
