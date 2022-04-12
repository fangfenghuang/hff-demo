package resource

import (
	"go-demo/pkg/e"
	"net/http"

	"github.com/gin-gonic/gin"
	"k8s.io/klog"
)

type Parms struct {
	Weight string `json:"weight"`
}

func HelloWorld(ctx *gin.Context) {
	appG := Gin{C: ctx}
	id := ctx.Param("id")                          //取得URL中参数
	name := ctx.Query("name")                      //查询请求URL后面的参数
	sex := ctx.DefaultQuery("sex", "女")            //查询请求URL后面的参数，如果没有填写默认值
	age := ctx.PostForm("age")                     //从表单中查询参数(form-data)
	height := ctx.DefaultPostForm("height", "160") //从表单中查询参数，如果没有填写默认值
	var parm Parms
	err := ctx.BindJSON(&parm)
	if err != nil {
		klog.Errorln(err)
	}
	data := map[string]string{
		"id":     id,
		"name":   name,
		"sex":    sex,
		"age":    age,
		"height": height,
		"weight": parm.Weight,
	}

	appG.Response(http.StatusOK, e.SUCCESS, data)
}
