package calico

import (
	"errors"
	"fmt"
	"go-demo/conf"
	"go-demo/pkg/e"
	"go-demo/pkg/k8sclient"

	v3 "github.com/projectcalico/libcalico-go/lib/apis/v3"
	"github.com/projectcalico/libcalico-go/lib/options"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

//清空网络策略
func ClearAllNetworkPolicy() error {
	nps, err := ClientV3.NetworkPolicies().List(Ctx, options.ListOptions{})
	if err != nil {
		klog.Errorln(err)
	}
	//ffmt.Puts(nps)
	for _, np := range nps.Items {
		klog.Infof("del np [%s/%s]", np.Name, np.Namespace)
		_, err := ClientV3.NetworkPolicies().Delete(Ctx, np.Namespace, np.Name, options.DeleteOptions{})
		if err != nil {
			klog.Errorln(err)
			return err
		}

	}
	gnps, err := ClientV3.GlobalNetworkPolicies().List(Ctx, options.ListOptions{})
	if err != nil {
		klog.Errorln(err)
	}
	//ffmt.Puts(gnps)
	for _, gnp := range gnps.Items {
		klog.Infof("del gnp [%s]", gnp.Name)
		_, err := ClientV3.GlobalNetworkPolicies().Delete(Ctx, gnp.Name, options.DeleteOptions{})
		if err != nil {
			klog.Errorln(err)
			return err
		}
	}

	gnss, err := ClientV3.GlobalNetworkSets().List(Ctx, options.ListOptions{})
	if err != nil {
		klog.Errorln(err)
	}
	//ffmt.Puts(gnps)
	for _, gns := range gnss.Items {
		klog.Infof("del gns [%s]", gns.Name)
		_, err := ClientV3.GlobalNetworkSets().Delete(Ctx, gns.Name, options.DeleteOptions{})
		if err != nil {
			klog.Errorln(err)
			return err
		}
	}

	return nil
}

//一键清空所有网络策略及节点标签
func CleanUp() error {
	//删除所有网络策略
	klog.Infoln("CleanUp...")
	err := ClearAllNetworkPolicy()
	if err != nil {
		return err
	}
	allNs, err := k8sclient.K8sClientSet.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	//清理节点标签
	for _, n := range allNs.Items {
		for _, k := range conf.LabelKeys {
			if _, ok := n.Labels[k]; ok {
				klog.Infof("rm ns %s label :%s", n.Name, k)
				if err := k8sclient.DeleteNamespaceLabels(n.Name, k); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// TODO: 优化校验
func CheckNetworkPolicyExist() (bool, error) {
	ns, err := k8sclient.K8sClientSet.CoreV1().Namespaces().Get("kube-system", metav1.GetOptions{})
	if err != nil {
		return true, err
	}
	if value := ns.Labels[conf.LABEL_ALLOW_ALL]; value == "true" {
		klog.Infoln("[jump]: calico networkpolicy init already")
		return true, nil
	}
	return false, nil
}

// 初始化网络策略
func CalicoNetworkPolicyInit() error {
	/**
	1.判断是否已经初始化
	2.清空网络策略
	3.给节点打上标签
	4.初始化网络策略
	**/
	klog.Infoln("CalicoNetworkPolicyInit...")
	// 1.校验是否已初始化
	if npExist, err := CheckNetworkPolicyExist(); npExist {
		return err
	}

	// 2.cleanup
	if err := CleanUp(); err != nil {
		return err
	}

	// 3.给节点打标签
	allNs, err := k8sclient.K8sClientSet.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, n := range allNs.Items {
		if err := k8sclient.UpdateNamespaceLabels(n.Name, conf.LABEL_NAME, n.Name); err != nil {
			return err
		}
	}
	for _, n := range conf.AllowAllNamespacesInit {
		if err := k8sclient.UpdateNamespaceLabels(n, conf.LABEL_ALLOW_ALL, "true"); err != nil {
			return err
		}
	}

	//4.初始化网络策略
	if err := ParseAndCreateGlobalNetworkSet(conf.GlobalNetworkSetPATH[conf.AREA_GLOBAL_EGRESS_DENY]); err != nil {
		return err
	}
	for _, fp := range conf.GlobalNetworkPolicyPath {
		if err := ParseAndCreateGlobalNetworkPolicy(fp); err != nil {
			return err
		}
	}

	return nil
}

func InitGlobalNetworkPolicy() {
	if err := CalicoSetup(); err != nil {
		klog.Fatal(err)
	}
	if err := CalicoNetworkPolicyInit(); err != nil {
		if err := CleanUp(); err != nil {
			klog.Fatal(err)
		}
		klog.Fatal(err)
	}
}

func OpenNamespaceNetworkPolicy(name string) (int, error) {
	klog.Infof("open ns [%s] networkpolicy", name)
	//1.添加global-internal-allow，如果已存在则跳过
	gnp, err := ClientV3.GlobalNetworkPolicies().Get(Ctx, conf.GLOBAL_INTERNAL_ALLOW, options.GetOptions{})
	if err != nil {
		return e.ERROR_INTERNAL_SERVER_ERROR, err
	}
	flag := false
	for i := 0; i < len(gnp.Spec.Ingress); i++ {
		if gnp.Spec.Ingress[i].Destination.NamespaceSelector == fmt.Sprintf("name == '%s'", name) {
			flag = true
		}
	}
	if !flag {
		gnp.Spec.Ingress = append(gnp.Spec.Ingress, v3.Rule{
			Action: v3.Allow,
			Destination: v3.EntityRule{
				NamespaceSelector: fmt.Sprintf("name == '%s'", name),
			},
			Source: v3.EntityRule{
				NamespaceSelector: fmt.Sprintf("name == '%s'", name),
			},
		})
		_, err = ClientV3.GlobalNetworkPolicies().Update(Ctx, gnp, options.SetOptions{})
		if err != nil {
			return e.ERROR_INTERNAL_SERVER_ERROR, err
		}
	}

	//2.添加ns标签，name、open-policy，如果已存在则跳过
	if err := k8sclient.UpdateNamespaceLabels(name, conf.LABEL_NAME, name); err != nil {
		return e.ERROR_INTERNAL_SERVER_ERROR, err
	} //TODO:此步骤可能可以省去，在创建ns时添加？
	if err := k8sclient.UpdateNamespaceLabels(name, conf.LABEL_ALLOW_ALL, "true"); err != nil {
		return e.ERROR_INTERNAL_SERVER_ERROR, err
	}

	return e.SUCCESS, nil
}

func CloseNamespaceNetworkPolicy(name string) (int, error) {
	klog.Infof("close ns [%s] networkpolicy", name)
	//1.global-internal-allow删除一项
	gnp, err := ClientV3.GlobalNetworkPolicies().Get(Ctx, conf.GLOBAL_INTERNAL_ALLOW, options.GetOptions{})
	if err != nil {
		return e.ERROR_INTERNAL_SERVER_ERROR, err
	}

	for i := 0; i < len(gnp.Spec.Ingress); i++ {
		if gnp.Spec.Ingress[i].Destination.NamespaceSelector == fmt.Sprintf("name == '%s'", name) {
			gnp.Spec.Ingress = append(gnp.Spec.Ingress[:i], gnp.Spec.Ingress[i+1:]...)
			i--
		}
	}

	_, err = ClientV3.GlobalNetworkPolicies().Update(Ctx, gnp, options.SetOptions{})
	if err != nil {
		return e.ERROR_INTERNAL_SERVER_ERROR, err
	}

	//2.删除ns标签，open-policy
	ns, err := k8sclient.K8sClientSet.CoreV1().Namespaces().Get(name, metav1.GetOptions{})
	if err != nil {
		return e.ERROR_INTERNAL_SERVER_ERROR, err
	}

	if value := ns.Labels[conf.LABEL_ALLOW_ALL]; value == "true" {
		if err := k8sclient.DeleteNamespaceLabels(name, conf.LABEL_ALLOW_ALL); err != nil {
			return e.ERROR_INTERNAL_SERVER_ERROR, err
		}
	}

	return e.SUCCESS, nil
}

func AddNetsBlackList(net string) (int, error) {
	klog.Infof("add nets blacklist: [%s]", net)
	/**1.  TODO: check net
		ip/cidr格式是否合法
		是否禁止输入的网段？0.0.0.0/0？
	**/

	//2.添加GlobalNetworkSet，如果已存在则跳过
	gns, err := ClientV3.GlobalNetworkSets().Get(Ctx, conf.AREA_GLOBAL_EGRESS_DENY, options.GetOptions{})
	if err != nil {
		return e.ERROR_INTERNAL_SERVER_ERROR, err
	}
	for _, v := range gns.Spec.Nets {
		if v == net {
			klog.Infof("net [%s] exist in %s", net, conf.AREA_GLOBAL_EGRESS_DENY)
			return e.SUCCESS, nil
		}
	}
	gns.Spec.Nets = append(gns.Spec.Nets, net)
	_, err = ClientV3.GlobalNetworkSets().Update(Ctx, gns, options.SetOptions{})
	if err != nil {
		return e.ERROR_INTERNAL_SERVER_ERROR, err
	}

	return e.SUCCESS, nil
}

func DeleteNetsBlackList(net string) (int, error) {
	klog.Infof("del nets blacklist: [%s]", net)
	/**1.  TODO: check net
		ip/cidr格式是否合法
		是否禁止输入的网段？0.0.0.0/0？
	**/

	//2. 删除GlobalNetworkSet相关项
	gns, err := ClientV3.GlobalNetworkSets().Get(Ctx, conf.AREA_GLOBAL_EGRESS_DENY, options.GetOptions{})
	if err != nil {
		return e.ERROR_INTERNAL_SERVER_ERROR, err
	}
	flag := false
	for i := 0; i < len(gns.Spec.Nets); i++ {
		if gns.Spec.Nets[i] == net {
			gns.Spec.Nets = append(gns.Spec.Nets[:i], gns.Spec.Nets[i+1:]...)
			i--
			flag = true
		}
	}
	if !flag {
		klog.Infof("[jump]: net [%s] not exist in %s", net, conf.AREA_GLOBAL_EGRESS_DENY)
		return e.SUCCESS, nil
	}
	_, err = ClientV3.GlobalNetworkSets().Update(Ctx, gns, options.SetOptions{})
	if err != nil {
		return e.ERROR_INTERNAL_SERVER_ERROR, err
	}

	return e.SUCCESS, nil
}

func AddAllowAllNamespace(name string) (int, error) {
	//检查是否未开启网络策略：是否有open-policy label
	ns, err := k8sclient.K8sClientSet.CoreV1().Namespaces().Get(name, metav1.GetOptions{})
	if err != nil {
		return e.ERROR_INTERNAL_SERVER_ERROR, err
	}
	if _, ok := ns.Labels[conf.LABEL_OPEN_POLICY]; ok {
		return e.INVALID_PARAMS, errors.New("allow all forbbiden to open-policy namespace")
	}

	if err := k8sclient.UpdateNamespaceLabels(name, conf.LABEL_NAME, name); err != nil {
		return e.ERROR_INTERNAL_SERVER_ERROR, err
	} //TODO:此步骤可能可以省去，在创建ns时添加？

	if err := k8sclient.UpdateNamespaceLabels(name, conf.LABEL_ALLOW_ALL, "true"); err != nil {
		return e.ERROR_INTERNAL_SERVER_ERROR, err
	}
	return e.SUCCESS, nil
}

func DeleteAllowAllNamespace(name string) (int, error) {
	// TODO: 检查容器云内部ns不允许上传
	for _, v := range conf.AllowAllNamespacesInit {
		if v == name {
			return e.INVALID_PARAMS, fmt.Errorf("ns %s is forbbidon for delete from %s", name, conf.LABEL_ALLOW_ALL)
		}
	}

	ns, err := k8sclient.K8sClientSet.CoreV1().Namespaces().Get(name, metav1.GetOptions{})
	if err != nil {
		return e.ERROR_INTERNAL_SERVER_ERROR, err
	}
	if value := ns.Labels[conf.LABEL_ALLOW_ALL]; value == "true" {
		if err := k8sclient.DeleteNamespaceLabels(name, conf.LABEL_ALLOW_ALL); err != nil {
			return e.ERROR_INTERNAL_SERVER_ERROR, err
		}
	} else {
		klog.Infof("[jump]: label %s not exist in ns %s", conf.LABEL_ALLOW_ALL, name)
		return e.SUCCESS, nil
	}
	return e.SUCCESS, nil
}

func GetGlobalNetworkPolicy() (int, *v3.GlobalNetworkPolicyList, error) {
	gnpList, err := ClientV3.GlobalNetworkPolicies().List(Ctx, options.ListOptions{})
	if err != nil {
		klog.Errorln(err)
		return e.ERROR_INTERNAL_SERVER_ERROR, gnpList, err
	}
	return e.SUCCESS, gnpList, nil
}
