package main

import (
	"github.com/gin-gonic/gin"
	DB "sequency/db"
	R "sequency/routes"
	RB "sequency/utils/mq"
)

func main() {

	r := gin.Default()

	DB.ConnectDB()
	connection := RB.MQ()

	a := RB.ConnectionMQ{MQ: connection}

	R.Routes(r, &a, "v1")

	go a.PollMq()

	defer connection.Close()

	r.Run()
}
