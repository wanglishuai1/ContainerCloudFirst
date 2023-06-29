package core

import (
	"ContainerCloudFirst/lib"
	"fmt"
	"k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"sync"
)

type DeploymentMap struct {
	data sync.Map // [key string] []*v1.Deployment    key=>namespace
}

func (this *DeploymentMap) Add(dep *v1.Deployment) {
	if list, ok := this.data.Load(dep.Namespace); ok {
		list = append(list.([]*v1.Deployment), dep)
		this.data.Store(dep.Namespace, list)
	} else {
		this.data.Store(dep.Namespace, []*v1.Deployment{dep})
	}
}

// 更新
func (this *DeploymentMap) Update(dep *v1.Deployment) error {
	if list, ok := this.data.Load(dep.Namespace); ok {
		for i, range_dep := range list.([]*v1.Deployment) {
			if range_dep.Name == dep.Name {
				list.([]*v1.Deployment)[i] = dep
			}
		}
		return nil
	}
	return fmt.Errorf("deployment-%s not found", dep.Name)
}

// 删除
func (this *DeploymentMap) Delete(dep *v1.Deployment) {
	if list, ok := this.data.Load(dep.Namespace); ok {
		for i, range_dep := range list.([]*v1.Deployment) {
			if range_dep.Name == dep.Name {
				newList := append(list.([]*v1.Deployment)[:i], list.([]*v1.Deployment)[i+1:]...)
				this.data.Store(dep.Namespace, newList)
				break
			}
		}
	}
}

func (this *DeploymentMap) ListByNS(ns string) ([]*v1.Deployment, error) {
	if list, ok := this.data.Load(ns); ok {
		return list.([]*v1.Deployment), nil
	}
	return nil, fmt.Errorf("record not found")
}
func (this *DeploymentMap) GetDeployment(ns string, depname string) (*v1.Deployment, error) {
	if list, ok := this.data.Load(ns); ok {
		for _, item := range list.([]*v1.Deployment) {
			if item.Name == depname {
				return item, nil
			}
		}
	}
	return nil, fmt.Errorf("record not found")
}

var DepMap *DeploymentMap //作为全局对象
func init() {
	DepMap = &DeploymentMap{}
}

type DepHandler struct{}

func (this *DepHandler) OnAdd(obj interface{}) {
	DepMap.Add(obj.(*v1.Deployment))
}
func (this *DepHandler) OnUpdate(oldObj, newObj interface{}) {
	err := DepMap.Update(newObj.(*v1.Deployment))
	if err != nil {
		fmt.Println(err)
	}
}
func (this *DepHandler) OnDelete(obj interface{}) {
	if d, ok := obj.(*v1.Deployment); ok {
		DepMap.Delete(d)
	}
}

func InitDeployment() {
	fact := informers.NewSharedInformerFactory(lib.K8sClient, 0)

	depInformer := fact.Apps().V1().Deployments()
	depInformer.Informer().AddEventHandler(&DepHandler{})

	podInformer := fact.Core().V1().Pods()
	podInformer.Informer().AddEventHandler(&PodHandler{})

	rsInformer := fact.Apps().V1().ReplicaSets()
	rsInformer.Informer().AddEventHandler(&RsHandler{})

	eventInformer := fact.Core().V1().Events()
	eventInformer.Informer().AddEventHandler(&EventHandler{})

	fact.Start(wait.NeverStop)

}
