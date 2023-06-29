package main

import (
	"ContainerCloudFirst/core"
	deployment "ContainerCloudFirst/deployments"
	"ContainerCloudFirst/lib"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	r := gin.New()
	r.Static("/static", "./static")
	r.LoadHTMLGlob("html/**/*")
	deployment.ReqHandlers(r)
	r.GET("/deployments", func(c *gin.Context) {
		c.HTML(http.StatusOK, "deployment_list.html",
			lib.DataBuilder().
				SetTitle("deployment列表").
				SetData("DepList", deployment.ListAll("default")))
	})
	r.GET("/deployments/:name", func(c *gin.Context) {
		c.HTML(http.StatusOK, "deployment_detail.html",
			lib.DataBuilder().
				SetTitle("deployment详情-"+c.Param("name")).
				SetData("DepDetail", deployment.GetDeployment("default", c.Param("name"))))
	})
	core.InitDeployment()
	r.Run(":8080")

}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

//func main() {
//
//	list, err := lib.K8sClient.CoreV1().Pods("default").List(context.Background(), metav1.ListOptions{})
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, pod := range list.Items {
//		fmt.Println(pod.Name)
//	}
//	log.Println("-----Service-------")
//
//	serviceList, err := lib.K8sClient.CoreV1().Services("default").List(context.Background(), metav1.ListOptions{})
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, service := range serviceList.Items {
//		fmt.Println(service.Name)
//	}
//	log.Println("---Node---------")
//	nodeList, err := lib.K8sClient.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, node := range nodeList.Items {
//		fmt.Println(node.Name)
//	}
//	log.Println("-----Deployment-------")
//	deployMentList, err := lib.K8sClient.AppsV1().Deployments("kube-system").List(context.Background(), metav1.ListOptions{})
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, deployment := range deployMentList.Items {
//		fmt.Println(deployment.Name, deployment.Namespace)
//	}
//
//	ngxDep := &v1.Deployment{}
//	b, _ := ioutil.ReadFile("yamls/nginx.yaml")
//	err = yaml.Unmarshal(b, ngxDep)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	createOpt := metav1.CreateOptions{}
//	_, err = lib.K8sClient.AppsV1().Deployments("default").Create(context.Background(), ngxDep, createOpt)
//	CheckError(err)
//
//}
