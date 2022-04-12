package resource

import (
	"go-demo/pkg/calico"
	"net/http"

	"github.com/gin-gonic/gin"
	"k8s.io/klog"
)

func OpenNamespaceNetworkPolicy(ctx *gin.Context) {
	appG := Gin{C: ctx}
	name := ctx.Param("name")
	code, err := calico.OpenNamespaceNetworkPolicy(name)

	if err == nil {
		appG.Response(http.StatusOK, code, nil)
	} else {
		klog.Errorln(err)
		appG.Response(http.StatusBadRequest, code, nil)
	}
}
func CloseNamespaceNetworkPolicy(ctx *gin.Context) {
	appG := Gin{C: ctx}
	name := ctx.Param("name")
	code, err := calico.CloseNamespaceNetworkPolicy(name)

	if err == nil {
		appG.Response(http.StatusOK, code, nil)
	} else {
		klog.Errorln(err)
		appG.Response(http.StatusBadRequest, code, nil)
	}
}

func AddNetsBlackList(ctx *gin.Context) {
	appG := Gin{C: ctx}
	net := ctx.Param("net")
	code, err := calico.AddNetsBlackList(net)

	if err == nil {
		appG.Response(http.StatusOK, code, nil)
	} else {
		klog.Errorln(err)
		appG.Response(http.StatusBadRequest, code, nil)
	}
}

func DeleteNetsBlackList(ctx *gin.Context) {
	appG := Gin{C: ctx}
	net := ctx.Param("net")
	code, err := calico.DeleteNetsBlackList(net)

	if err == nil {
		appG.Response(http.StatusOK, code, nil)
	} else {
		klog.Errorln(err)
		appG.Response(http.StatusBadRequest, code, nil)
	}
}

func AddAllowAllNamespace(ctx *gin.Context) {
	appG := Gin{C: ctx}
	name := ctx.Param("name")
	code, err := calico.AddAllowAllNamespace(name)

	if err == nil {
		appG.Response(http.StatusOK, code, nil)
	} else {
		klog.Errorln(err)
		appG.Response(http.StatusBadRequest, code, nil)
	}
}

func DeleteAllowAllNamespace(ctx *gin.Context) {
	appG := Gin{C: ctx}
	name := ctx.Param("name")
	code, err := calico.DeleteAllowAllNamespace(name)

	if err == nil {
		appG.Response(http.StatusOK, code, nil)
	} else {
		klog.Errorln(err)
		appG.Response(http.StatusBadRequest, code, nil)
	}
}

func GetGlobalNetworkPolicy(ctx *gin.Context) {
	appG := Gin{C: ctx}
	code, data, err := calico.GetGlobalNetworkPolicy()

	if err == nil {
		appG.Response(http.StatusOK, code, data)
	} else {
		klog.Errorln(err)
		appG.Response(http.StatusBadRequest, code, data)
	}
}
