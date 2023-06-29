package deployment

import (
	"ContainerCloudFirst/lib"
	"context"
	"fmt"
	v1 "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetImages(dep v1.Deployment) string {
	return GetImagesByPod(dep.Spec.Template.Spec.Containers)
}
func GetImagesByPod(container []core.Container) string {
	images := container[0].Image
	if imgLen := len(container); imgLen > 1 {
		images = fmt.Sprintf("+其他%d个镜像", imgLen-1)
	}
	return images

}
func GetLabels(m map[string]string) string {
	labels := ""
	// aa=xxx,xxx=xx
	for k, v := range m {
		if labels != "" {
			labels += ","
		}
		labels += fmt.Sprintf("%s=%s", k, v)
	}
	return labels
}
func GetReplicas(item v1.Deployment) [3]int32 {
	return [3]int32{item.Status.Replicas, item.Status.AvailableReplicas, item.Status.UnavailableReplicas}
}

func GetRsLableByDeploymentListWatch(dep *v1.Deployment, rslist []*v1.ReplicaSet) ([]map[string]string, error) {
	ret := make([]map[string]string, 0)
	for _, item := range rslist {
		if IsRsFromDep(dep, *item) {
			s, err := metav1.LabelSelectorAsMap(item.Spec.Selector)
			if err != nil {
				return nil, err
			}
			ret = append(ret, s)
		}
	}
	return ret, nil
}

// 获取这个deployment 的rs标签 普通调用API
func GetRsLableByDeployment(dep *v1.Deployment) string {
	selector, _ := metav1.LabelSelectorAsSelector(dep.Spec.Selector)
	listOpt := metav1.ListOptions{
		LabelSelector: selector.String(),
	}
	rs, _ := lib.K8sClient.AppsV1().ReplicaSets(dep.Namespace).List(context.Background(), listOpt)
	for _, item := range rs.Items {
		if IsCurrentRsByDep(dep, item) {
			asSelector, err := metav1.LabelSelectorAsSelector(item.Spec.Selector)
			if err != nil {
				return ""
			}
			return asSelector.String()
		}
	}
	return ""

}

// 判断 rs 是否属于 某个 dep
func IsRsFromDep(dep *v1.Deployment, set v1.ReplicaSet) bool {
	for _, ref := range set.OwnerReferences {
		if ref.Kind == "Deployment" && ref.Name == dep.Name {
			return true
		}
	}
	return false
}

// IsCurrentRsByDep 判断rs是不是当前最新版的deployment 的rs
func IsCurrentRsByDep(dep *v1.Deployment, set v1.ReplicaSet) bool {
	if set.ObjectMeta.Annotations["deployment.kubernetes.io/revision"] != dep.ObjectMeta.Annotations["deployment.kubernetes.io/revision"] {
		return false
	}
	return IsRsFromDep(dep, set)

}
