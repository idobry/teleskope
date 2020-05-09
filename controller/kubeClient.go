package controller

import (
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/idob/deploymentStatus/internal"
)

func GetKubeClient() *kubernetes.Clientset{
	var cfg *rest.Config
	var err error
	cfg, err = rest.InClusterConfig()
	home := homeDir()
	if err != nil && home != "" {
		cfg, err = clientcmd.BuildConfigFromFlags("", home+"/.kube/config")
		internal.ExitOnErr(err)
	} else {
		cfg, err = clientcmd.BuildConfigFromFlags("", "")
		internal.ExitOnErr(err)
	}

	kubeclient, err := kubernetes.NewForConfig(cfg)
	internal.ExitOnErr(err)

	return kubeclient
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
