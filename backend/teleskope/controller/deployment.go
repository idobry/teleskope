package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

type Container struct{
	Name string
	Image string
	Envs []string
}

type deployments struct {
	ID []string
}

type DeploymentEvent struct{
	Name string
	Namespace string
	Containers []Container
	ReplicaCurrent string
	ReplicaDesired string
}

func GetDeployment(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Inside GetDeployment func\n")
	vars := mux.Vars(r)
	ns := vars["ns"]
	dep := vars["dep"]
	fmt.Printf("dep %s in $s \n", dep, ns)
	kubeclient := GetKubeClient()
	newDepl, err := kubeclient.AppsV1().Deployments(mux.Vars(r)["ns"]).Get(mux.Vars(r)["dep"],metav1.GetOptions{})
	if err != nil{
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(getDeployEvent(newDepl))
	if err != nil{
		panic(err)
	}
	return
}

func GetDeployments(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Inside GetDeployments func\n")
	fmt.Printf("ns %s \n", mux.Vars(r)["ns"])
	kubeclient := GetKubeClient()
	deploymetns, err := kubeclient.AppsV1().Deployments(mux.Vars(r)["ns"]).List(metav1.ListOptions{})
	if err != nil{
		panic(err)
	}
	depList := deployments{
		[]string{},
	}
	for _, dep := range deploymetns.Items{
		depList.ID = append(depList.ID, dep.Name)
	}
	b, err := json.Marshal(depList)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, _ = w.Write(b)
	return
}

func StreamDeployments(h *Hub) {
	fmt.Printf("Inside StreamDeployments func\n")
	kubeclient := GetKubeClient()

	factory := informers.NewSharedInformerFactory(kubeclient, 0)
	deploymentInformer := factory.Apps().V1().Deployments().Informer()
	stopper := make(chan struct{})
	defer close(stopper)

	deploymentInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {},
		UpdateFunc: func(old, new interface{}) {
			newDepl := new.(*appsv1.Deployment)
			fmt.Printf("deploymet %s changed\n", newDepl.Name)

			h.broadcast <- getDeployEvent(newDepl)

			return
		},
		DeleteFunc: func(obj interface{}) {},
	})

	deploymentInformer.Run(stopper)
}

func getDeployEvent(newDepl *appsv1.Deployment) []byte{
	deploy := DeploymentEvent{
		Name: newDepl.Name,
		Namespace: newDepl.Namespace,
		ReplicaCurrent: fmt.Sprintf("%d", newDepl.Status.AvailableReplicas),
		ReplicaDesired: fmt.Sprintf("%d", newDepl.Status.Replicas),
		Containers: []Container{},
	}
	for _ ,c := range newDepl.Spec.Template.Spec.Containers{
		var con = Container{}
		con.Name = c.Name
		con.Image = c.Image
		con.Envs = []string{}
		for _, env := range c.Env{
			con.Envs = append(con.Envs, env.Name+":"+env.Value)
		}
		deploy.Containers = append(deploy.Containers, con)
	}

	b, err := json.Marshal(deploy)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return b
}
