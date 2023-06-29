package deployment

import (
	"ContainerCloudFirst/lib"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"log"
)

func Create() {
	ctx := context.Background()

	listopt := metav1.ListOptions{}
	depList, err := lib.K8sClient.AppsV1().Deployments("default").List(ctx, listopt)

	for _, item := range depList.Items { //遍历所有deployment
		fmt.Println(item.Name)
	}
	ngxDep := &v1.Deployment{} //我们要创建的deployment
	b, _ := ioutil.ReadFile("yamls/nginx.yaml")
	ngxJson, _ := yaml.ToJSON(b)
	lib.CheckError(json.Unmarshal(ngxJson, ngxDep))

	createopt := metav1.CreateOptions{}

	_, err = lib.K8sClient.AppsV1().Deployments("myweb").
		Create(ctx, ngxDep, createopt)

	if err != nil {
		log.Fatal(err)
	}
}
