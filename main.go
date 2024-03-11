package main

import (
	"github.com/gin-gonic/gin"
	DB "sequency/db"
	R "sequency/routes"
	RB "sequency/utils/mq"
)

func main() {

	r := gin.Default()

	DBCONNECT := DB.ConnectDB()
	connection := RB.MQ()

	MQCONNECT := RB.ConnectionMQ{MQ: connection}

	routes := R.RoutesDe{DB: DBCONNECT, MQ: MQCONNECT}

	routes.Routes(r, "v1")

	go MQCONNECT.PollMq()

	defer connection.Close()

	r.Run()
}
