package main

import (
	"VmSSH/config"
	"VmSSH/router"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	router.Router.InitApiRouter(r)
	yaml := config.Yaml{}
	yaml.LoadToml()
	addressBind := fmt.Sprintf("%s:%d", config.Conf.Server.Host ,config.Conf.Server.Port)
	r.Run(addressBind)

}
