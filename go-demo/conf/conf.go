package conf

const LISTENPORT = "40080"

const STATIC_PATH = "/"

var (
	EtcdEndpoints      = "10.19.0.13:2379"
	EtcdKeyFile        = "./conf/ssl/etcd-key.pem"
	EtcdCertFile       = "./conf/ssl/etcd.pem"
	EtcdCACertFile     = "./conf/ssl/ca.pem"
	KubeconfigPath     = "./conf/kube/config"
	LeaseLockNamespace = "kube-system"
	LeaseLockName      = "go-demo"
)

const (
	LABEL_ALLOW_ALL   = "allow-all"
	LABEL_NAME        = "name"
	LABEL_OPEN_POLICY = "open-policy"
)

var LabelKeys = []string{
	LABEL_ALLOW_ALL,
	LABEL_NAME,
	LABEL_OPEN_POLICY,
}

var AllowAllNamespacesInit = []string{}

const AREA_GLOBAL_EGRESS_DENY = "area-global-egress-deny"

var GlobalNetworkSetPATH = map[string]string{
	AREA_GLOBAL_EGRESS_DENY: "./conf/yaml/area-global-egress-deny.yaml",
}

const (
	GLOBAL_DEFAULT        = "global-default"
	GLOBAL_ALLOW_ALL      = "global-allow-all"
	GLOBAL_INTERNAL_ALLOW = "global-internal-allow"
	GLOBAL_EGRESS_DENY    = "global-egress-deny"
)

var GlobalNetworkPolicyPath = map[string]string{
	GLOBAL_DEFAULT:        "./conf/yaml/global-default.yaml",
	GLOBAL_ALLOW_ALL:      "./conf/yaml/global-allow-all.yaml",
	GLOBAL_INTERNAL_ALLOW: "./conf/yaml/global-internal-allow.yaml",
	GLOBAL_EGRESS_DENY:    "./conf/yaml/global-egress-deny.yaml",
}
