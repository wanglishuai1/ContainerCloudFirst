package deployment

import (
	"ContainerCloudFirst/core"
	"ContainerCloudFirst/lib"
)

// 显示 所有
func ListAll(namespace string) (ret []*Deployment) {
	//ctx := context.Background()
	//listopt := metav1.ListOptions{}
	//depList, err := lib.K8sClient.AppsV1().Deployments(namespace).List(ctx, listopt)
	depList, err := core.DepMap.ListByNS(namespace)
	lib.CheckError(err)
	for _, item := range depList { //遍历所有deployment

		ret = append(ret, &Deployment{Name: item.Name,
			Images:   GetImages(*item),
			Replicas: GetReplicas(*item),
		})
	}
	return
}
func ListPodsByLabel(ns string, labels []map[string]string) (ret []*Pod) {
	list, err := core.PodMap.ListByLabels(ns, labels)
	lib.CheckError(err)
	for _, pod := range list {
		ret = append(ret, &Pod{
			Name:       pod.Name,
			Images:     GetImagesByPod(pod.Spec.Containers),
			Node:       pod.Spec.NodeName,
			CreateTime: pod.CreationTimestamp.Format("2006-01-02 15:04:05"),
			Phase:      string(pod.Status.Phase),
			Message:    core.EventMap.GetMessage(pod.Namespace, "Pod", pod.Name),
			IsReady:    GetPodIsReady(*pod),
			IP: []string{
				pod.Status.PodIP,
				pod.Status.HostIP,
			},
		})

	}
	return
}
