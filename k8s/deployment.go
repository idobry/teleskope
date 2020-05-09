package k8s

import (
	"fmt"
	"net/http"

	"github.com/idob/deploymentStatus/internal"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

const manifestsDir = "manifests"

// Use an empty string to run on all namespaces
const namespace = ""
const newLabelKey = "new-label-to-add"
const tempLabelKey = "temporary"
const tempSuffix = "-temp"
const componentLabelKey = "component"

func GetDeploy(w http.ResponseWriter, r *http.Request) {

	kubeclient := GetKubeClient()

	deployments, err := kubeclient.AppsV1().Deployments(namespace).List(metav1.ListOptions{
		LabelSelector: componentLabelKey,
	})
	internal.ExitOnErr(err)
	_, err = w.Write([]byte(`{"hello" : "world1"}`))
	fmt.Printf("Got %d deployments.\n", len(deployments.Items))

	factory := informers.NewSharedInformerFactory(kubeclient, 0)
	deploymentInformer := factory.Apps().V1().Deployments().Informer()
	stopper := make(chan struct{})
	defer close(stopper)

	deploymentInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {},
		UpdateFunc: func(old, new interface{}) {
			newDepl := new.(*appsv1.Deployment)
			fmt.Printf("%s deployment updated.\n", newDepl.Name)
			//_, err = w.Write([]byte(`{"hello" : "`+ newDepl.Name +`"}`))

			return
		},
		DeleteFunc: func(obj interface{}) {},
	})

	deploymentInformer.Run(stopper)
}
