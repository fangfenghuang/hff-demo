package k8sclient

import (
	"gopkg.in/ffmt.v1"
)

func Test() {
	//测试回滚
	//err := Rollout("statefulset", "hff", "hff-sts", 11)
	//err := Rollout("deployment", "hffns", "test-scheduler", 25)
	// err := Rollout("daemonset", "hffns", "nginx-ds", 1)
	// if err != nil {
	// 	klog.Errorln(err)
	// }
	//测试节点信息
	// node, err := K8sClientSet.CoreV1().Nodes().Get("tztest", v1.GetOptions{})
	// if err != nil {
	// 	klog.Errorln(err)
	// 	return
	// }
	//ffmt.Puts(node)
	// ffmt.Puts(node.Status.Capacity.Cpu())
	// ffmt.Puts(node.Status.Allocatable.Cpu())

	ffmt.Puts("---------------------------------")
}
