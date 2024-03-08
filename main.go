package main

import (
	"github.com/gin-gonic/gin"
	DB "sequency/db"
	R "sequency/routes"
	NSQ "sequency/utils/nsq"
)

func main() {

	r := gin.Default()

	DB.ConnectDB()

	R.Routes(r, "v1")

	go NSQ.ProcessOrderNSQ()

	r.Run()
}
