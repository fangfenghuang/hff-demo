package resource

import (
	"go-demo/pkg/calico"
	"net/http"

	"github.com/gin-gonic/gin"
	v3 "github.com/projectcalico/libcalico-go/lib/apis/v3"
	"k8s.io/klog"
)

func CreateNetworkPolicy(ctx *gin.Context) {
	appG := Gin{C: ctx}
	var parm *v3.NetworkPolicy
	err := ctx.BindJSON(&parm)
	code := 0
	if err == nil {
		code, err = calico.CreateNetworkPolicy(parm)
	}

	if err == nil {
		appG.Response(http.StatusOK, code, nil)
	} else {
		klog.Errorln(err)
		appG.Response(http.StatusBadRequest, code, nil)
	}
}

func DeleteNetworkPolicy(ctx *gin.Context) {
	appG := Gin{C: ctx}
	namespace := ctx.Param("namespace")
	name := ctx.Param("name")
	code, err := calico.DeleteNetworkPolicy(namespace, name)

	if err == nil {
		appG.Response(http.StatusOK, code, nil)
	} else {
		klog.Errorln(err)
		appG.Response(http.StatusBadRequest, code, nil)
	}
}

func GetNetworkPolicy(ctx *gin.Context) {
	appG := Gin{C: ctx}
	code, data, err := calico.GetNetworkPolicy()

	if err == nil {
		appG.Response(http.StatusOK, code, data)
	} else {
		klog.Errorln(err)
		appG.Response(http.StatusBadRequest, code, data)
	}
}
