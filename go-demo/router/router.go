package router

import (
	"net/http"

	"go-demo/conf"
	"go-demo/router/resource"

	"github.com/gin-gonic/gin"
	"k8s.io/klog"
)

func InitRouter() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(Cors())

	router.GET("/helloworld/:id", resource.HelloWorld)

	router.GET("/api/v1/globalnetworkpolicys", resource.GetGlobalNetworkPolicy)
	router.POST("/api/v1/globalnetworkpolicy/openpolicy/:namespace", resource.OpenNamespaceNetworkPolicy)
	router.POST("/api/v1/globalnetworkpolicy/closepolicy/:namespace", resource.CloseNamespaceNetworkPolicy)
	router.POST("/api/v1/globalnetworkpolicy/nets-blacklist/:namespace", resource.AddNetsBlackList)
	router.DELETE("/api/v1/globalnetworkpolicy/nets-blacklist/:namespace", resource.DeleteNetsBlackList)
	router.POST("/api/v1/globalnetworkpolicy/allow-all/:namespace", resource.AddAllowAllNamespace)
	router.DELETE("/api/v1/globalnetworkpolicy/allow-all/:namespace", resource.DeleteAllowAllNamespace)

	router.GET("/api/v1/networkpolicys", resource.GetNetworkPolicy)
	router.POST("/api/v1/networkpolicys", resource.CreateNetworkPolicy)
	router.DELETE("/api/v1/networkpolicys/:namespace/:name", resource.DeleteNetworkPolicy)

	router.Static("/api/v1/files", conf.STATIC_PATH)

	if err := router.Run("0.0.0.0:" + conf.LISTENPORT); err != nil {
		klog.Errorln("service start fail: ", err)
	}

}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, developerId")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
