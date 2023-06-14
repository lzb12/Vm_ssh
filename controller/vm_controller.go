package controller

import (
	"VmSSH/servicer"
	"github.com/gin-gonic/gin"
)

func VMController(c *gin.Context)  {
	url := c.Query("url")
	username := c.Query("username")
	password := c.Query("password")

	//fmt.Println(url,username,password)

	vm := servicer.VmWare{}
	servicer.NewVmWare(url,username,password)
	vm.GetAllVmClient()
	//allhost,err  := vm.GetAllHost()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(allhost)

}