package calico

import (
	"go-demo/pkg/e"

	v3 "github.com/projectcalico/libcalico-go/lib/apis/v3"
	"github.com/projectcalico/libcalico-go/lib/options"
	ffmt "gopkg.in/ffmt.v1"
	"k8s.io/klog"
)

// 删除np
func DeleteNetworkPolicy(namespace, name string) (int, error) {
	klog.Infof("del np [%s/%s]...", namespace, name)
	_, err := ClientV3.NetworkPolicies().Delete(Ctx, namespace, name, options.DeleteOptions{})
	if err != nil {
		klog.Errorln(err)
		return e.ERROR_INTERNAL_SERVER_ERROR, err
	}
	return e.SUCCESS, nil
}

// 添加np
func CreateNetworkPolicy(policy *v3.NetworkPolicy) (int, error) {
	klog.Infof("Add np [%s/%s]...", policy.Name, policy.Namespace)
	ffmt.Puts(policy)

	_, err := ClientV3.NetworkPolicies().Create(Ctx, policy, options.SetOptions{})
	if err != nil {
		klog.Errorln(err)
		return e.ERROR_INTERNAL_SERVER_ERROR, err
	}
	return e.SUCCESS, nil
}

func GetNetworkPolicy() (int, *v3.NetworkPolicyList, error) {
	npList, err := ClientV3.NetworkPolicies().List(Ctx, options.ListOptions{})
	if err != nil {
		klog.Errorln(err)
		return e.ERROR_INTERNAL_SERVER_ERROR, npList, err
	}
	return e.SUCCESS, npList, nil
}
