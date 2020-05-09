package k8s

import (
	"flag"
	"os"
	"path/filepath"

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
		kubeconfig := flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		cfg, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		internal.ExitOnErr(err)
	} else {
		kubeconfig := flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		cfg, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		internal.ExitOnErr(err)
	}
	flag.Parse()


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
