package k8sclient

import (
	"encoding/json"
	"go-demo/conf"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
)

type RemoveStringValue struct {
	Op   string `json:"op"`
	Path string `json:"path"`
}

func InitClientSet() {
	// init clientset for k8s operation
	// if kubeconf path is empty, Using the inClusterConfig.
	config, err := clientcmd.BuildConfigFromFlags("", conf.KubeconfigPath)
	if err != nil {
		klog.Fatal(err)
	}

	K8sClientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatal(err)
	}
}

//添加或更新标签
func UpdateNamespaceLabels(name string, labelKey string, labelNewValue string) error {
	log.Printf("> label ns [%s]: %s=%s \n", name, labelKey, labelNewValue)
	ns, err := K8sClientSet.CoreV1().Namespaces().Get(name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	labels := ns.Labels
	if len(ns.Labels) == 0 {
		labels = make(map[string]string)
	}
	labels[labelKey] = labelNewValue

	patchData := map[string]interface{}{"metadata": map[string]map[string]string{"labels": labels}}
	playLoadBytes, _ := json.Marshal(patchData)

	_, err = K8sClientSet.CoreV1().Namespaces().Patch(name, types.StrategicMergePatchType, playLoadBytes)
	if err != nil {
		return err
	}
	return nil
}

func DeleteNamespaceLabels(name string, labelKey string) error {
	log.Printf("> del ns label [%s]: %s- \n", name, labelKey)
	payloads := []RemoveStringValue{
		{
			Op:   "remove",
			Path: "/metadata/labels/" + labelKey,
		},
	}

	playLoadBytes, _ := json.Marshal(payloads)

	_, err := K8sClientSet.CoreV1().Namespaces().Patch(name, types.JSONPatchType, playLoadBytes)
	if err != nil {
		return err
	}

	return nil
}
