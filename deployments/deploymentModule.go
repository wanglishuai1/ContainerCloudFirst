package deployment

type Pod struct {
	Namespace  string
	Name       string
	Images     string
	Node       string
	IsReady    bool
	CreateTime string
	Phase      string
	Message    string
	IP         []string
}
type Deployment struct {
	Namespace  string
	Name       string
	Replicas   [3]int32
	Images     string
	CreateTime string
	Pods       []*Pod
}
