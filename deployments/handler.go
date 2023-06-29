package deployment

import (
	"ContainerCloudFirst/core"
	"ContainerCloudFirst/lib"
	"context"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ReqHandlers(r *gin.Engine) {
	r.POST("/update/deployment/scale", incrReplicas)
	r.POST("/core/deployments", ListAllDeployments)
	r.POST("/core/pods", ListPodsByDeployment)
	r.DELETE("/core/pods", DeletePOD)
	r.GET("/core/pods_json", GetPodJson)
}

// 删除POD
func DeletePOD(c *gin.Context) {
	ns := c.DefaultQuery("namespace", "default")
	podName := c.DefaultQuery("pod", "")
	if podName == "" || ns == "" {
		panic("error ns or pod")
	}
	lib.CheckError(DeletePod(ns, podName))
	c.JSON(200, gin.H{"message": "Ok"})

}
func GetPodJson(c *gin.Context) {
	ns := c.DefaultQuery("namespace", "default")
	podName := c.DefaultQuery("pod", "")
	if podName == "" || ns == "" {
		panic("error")
	}
	if pod := core.PodMap.Get(ns, podName); pod == nil {
		panic("no such pod" + podName)
	} else {
		c.JSON(200, pod)
	}

}
func ListAllDeployments(c *gin.Context) {
	ns := c.DefaultQuery("namespace", "default")
	c.JSON(200, gin.H{"message": "Ok", "result": ListAll(ns)})
}

// 根据dep获取pod列表
func ListPodsByDeployment(c *gin.Context) {
	ns := c.DefaultQuery("namespace", "default")
	depname := c.DefaultQuery("deployment", "default")
	dep, err := core.DepMap.GetDeployment(ns, depname)
	lib.CheckError(err)
	rsList, err := core.RsMap.ListByNameSpace(ns) //根据namespace 获取到所有的rs
	lib.CheckError(err)
	labels, err := GetRsLableByDeploymentListWatch(dep, rsList)
	lib.CheckError(err)
	c.JSON(200, gin.H{
		"message": "Ok", "result": ListPodsByLabel(ns, labels),
	})
}
func incrReplicas(c *gin.Context) {
	req := struct {
		Namespace  string `json:"namespace" binding:"required,min=1"`
		Deployment string `json:"deployment" binding:"required,min=1"`
		Dec        bool   `json:"dec"`
	}{}
	lib.CheckError(c.ShouldBindJSON(&req))
	ctx := context.Background()
	scale, err := lib.K8sClient.AppsV1().Deployments(req.Namespace).GetScale(ctx, req.Deployment, v1.GetOptions{})
	lib.CheckError(err)
	if req.Dec {
		scale.Spec.Replicas--
	} else {
		scale.Spec.Replicas++
	}
	_, err = lib.K8sClient.AppsV1().Deployments(req.Namespace).UpdateScale(ctx, req.Deployment, scale, v1.UpdateOptions{})
	lib.CheckError(err)
	lib.Success("ok", c)
}
