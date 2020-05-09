package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)
type Container struct{
	Name string
	Image string
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

func GetDeployment(hub *Hub, w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Inside GetDeployment func\n")
	fmt.Printf("dep %s \n", mux.Vars(r)["dep"])
	kubeclient := GetKubeClient()
	newDepl, err := kubeclient.AppsV1().Deployments(mux.Vars(r)["ns"]).Get(mux.Vars(r)["dep"],metav1.GetOptions{})
	if err != nil{
		panic(err)
	}
	hub.broadcast <- getDeployEvent(newDepl)
	//_, _ = w.Write(getDeployEvent(newDepl))
	return
}

func GetDeployments(hub *Hub, w http.ResponseWriter, r *http.Request) {
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
		fmt.Println(dep.Name)
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
			/*deploy := DeploymentEvent{
				Name: newDepl.Name,
				Namespace: newDepl.Namespace,
				ReplicaCurrent: fmt.Sprintf("%d", newDepl.Status.AvailableReplicas),
				ReplicaDesired: fmt.Sprintf("%d", newDepl.Status.Replicas),
				Containers: []Container{},
			}
			for _ ,c := range newDepl.Spec.Template.Spec.Containers{
				deploy.Containers = append(deploy.Containers, Container{c.Name,c.Image})
			}

			b, err := json.Marshal(deploy)
			if err != nil {
				fmt.Println(err)
				return
			}*/

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
		deploy.Containers = append(deploy.Containers, Container{c.Name,c.Image})
	}

	b, err := json.Marshal(deploy)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return b
}