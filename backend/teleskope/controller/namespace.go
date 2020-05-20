package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

type NamespaceList struct{
	ID []string
}

func GetNamespaces(hub *Hub, w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Inside GetNamespaces func\n")
	kubeclient := GetKubeClient()

	namespaces, err := kubeclient.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil{
		panic(err)
	}
	nsList := NamespaceList{
		[]string{},
	}
	for _, ns := range namespaces.Items{
		nsList.ID = append(nsList.ID, ns.Name)
	}
	b, err := json.Marshal(nsList)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, _ = w.Write(b)
	return
}


func WatchNamespaces(h *Hub) {
	fmt.Printf("Inside StreamDeployments func\n")
	kubeclient := GetKubeClient()

	factory := informers.NewSharedInformerFactory(kubeclient, 0)
	namespaceInformer := factory.Core().V1().Namespaces().Informer()
	stopper := make(chan struct{})
	defer close(stopper)

	namespaceInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {},
		UpdateFunc: func(old, new interface{}) {
			h.broadcast <- []byte("")
			return
		},
		DeleteFunc: func(obj interface{}) {},
	})

	namespaceInformer.Run(stopper)
}
