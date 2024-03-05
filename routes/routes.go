package routes

import (
	"github.com/gin-gonic/gin"
	ROUTES "sequency/routes/router"
)

func Routes(R *gin.Engine, path string) {
	router := R.Group(path)

	router.GET("/wokflows", ROUTES.GetWorkflows)
}
