package main

import (
	"github.com/gin-gonic/gin"
	DB "sequency/db"
	R "sequency/routes"
)

func main() {

	r := gin.Default()

	DB.ConnectDB()

	R.Routes(r, "v1")

	r.Run()
}
