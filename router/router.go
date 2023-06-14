package router

import (
	"VmSSH/controller"
	"github.com/gin-gonic/gin"
)

var Router router

type router struct {

}

func (r *router) InitApiRouter(router *gin.Engine)  {
	//vm
	vm := router.Group("/api/v1/vm")
	vm.POST("/login",controller.VMController)
}