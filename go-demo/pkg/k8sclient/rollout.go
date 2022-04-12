package k8sclient

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/klog"
	"k8s.io/kubectl/pkg/polymorphichelpers"
)

const (
	DaemonSet   = "daemonset"
	StatefulSet = "statefulset"
	Deployment  = "deployment"
)

var rollbackKind = map[string]schema.GroupKind{
	DaemonSet:   {Group: "apps", Kind: "DaemonSet"},
	StatefulSet: {Group: "apps", Kind: "StatefulSet"},
	Deployment:  {Group: "apps", Kind: "Deployment"},
}

func Rollout(resource, namespace, name string, toRevision int64) error {

	kind, ok := rollbackKind[resource]
	if !ok {
		return fmt.Errorf("resource type error: %v", resource)
	}

	rollbacker, err := polymorphichelpers.RollbackerFor(kind, K8sClientSet)
	if err != nil {
		return fmt.Errorf("error getting Rollbacker for a %v: %v", kind.String(), err)
	}
	var obj runtime.Object
	switch resource {
	case DaemonSet:
		obj, err = K8sClientSet.AppsV1().DaemonSets(namespace).Get(name, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("failed to retrieve DaemonSet %s: %s", name, err.Error())
		}
	case StatefulSet:
		obj, err = K8sClientSet.AppsV1().StatefulSets(namespace).Get(name, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("failed to retrieve Statefulset %s: %s", name, err.Error())
		}
	case Deployment:
		obj, err = K8sClientSet.AppsV1().Deployments(namespace).Get(name, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("failed to retrieve Deployment %s: %s", name, err.Error())
		}
	}
	_, err = rollbacker.Rollback(obj, nil, toRevision, false)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	klog.Infof("%v %v/%v rollout to revision %d done", resource, namespace, name, toRevision)
	return nil
}
