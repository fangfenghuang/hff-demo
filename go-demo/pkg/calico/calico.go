package calico

import (
	"context"
	"encoding/json"
	"go-demo/conf"
	"io/ioutil"

	"github.com/projectcalico/libcalico-go/lib/apiconfig"
	v3 "github.com/projectcalico/libcalico-go/lib/apis/v3"
	"github.com/projectcalico/libcalico-go/lib/clientv3"
	"github.com/projectcalico/libcalico-go/lib/options"
	ffmt "gopkg.in/ffmt.v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	yaml2 "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/klog"
)

var ClientV3 clientv3.Interface
var Ctx context.Context

func CalicoSetup() error {
	var err error
	ClientV3, err = clientv3.New(apiconfig.CalicoAPIConfig{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{},
		Spec: apiconfig.CalicoAPIConfigSpec{
			DatastoreType: apiconfig.EtcdV3,
			EtcdConfig: apiconfig.EtcdConfig{
				EtcdEndpoints:  conf.EtcdEndpoints,
				EtcdKeyFile:    conf.EtcdKeyFile,
				EtcdCertFile:   conf.EtcdCertFile,
				EtcdCACertFile: conf.EtcdCACertFile,
			}},
	})

	if err != nil {
		return err
	}
	Ctx = context.Background()
	klog.Infoln("setup calico client succeed")
	return nil
}

func GetNetworkPolicie(name, namespace string) *v3.NetworkPolicy {
	var np *v3.NetworkPolicy
	np, err := ClientV3.NetworkPolicies().Get(Ctx, name, namespace, options.GetOptions{})
	if err != nil {
		klog.Errorln(err)
		return nil
	}
	ffmt.Puts(np)
	return np
}

func ParseAndCreateGlobalNetworkPolicy(fp string) error {
	klog.Infof("parse gnp yaml:%s ", fp)
	data, err := ioutil.ReadFile(fp)
	if err != nil {
		return err
	}
	if data, err = yaml2.ToJSON(data); err != nil {
		return err
	}
	var gnp v3.GlobalNetworkPolicy
	if err := json.Unmarshal(data, &gnp); err != nil {
		return err
	}
	//ffmt.Puts(gnp)
	_, err = ClientV3.GlobalNetworkPolicies().Create(Ctx, &gnp, options.SetOptions{})
	if err != nil {
		return err
	}
	return nil
}

func ParseAndCreateGlobalNetworkSet(file string) error {
	klog.Infof("parse gns yaml:%s ", file)
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	if data, err = yaml2.ToJSON(data); err != nil {
		return err
	}
	gns := v3.GlobalNetworkSet{}
	if err := json.Unmarshal(data, &gns); err != nil {
		return err
	}
	//ffmt.Puts(gnp)
	_, err = ClientV3.GlobalNetworkSets().Create(Ctx, &gns, options.SetOptions{})
	if err != nil {
		return err
	}
	return nil
}
