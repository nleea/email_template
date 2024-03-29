package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	ROUTES "sequency/routes/router"
	EMAIL "sequency/utils/emails"
	RB "sequency/utils/mq"
)

type RoutesDe struct {
	DB *mongo.Database
	MQ RB.ConnectionMQ
}

func (c RoutesDe) Routes(R *gin.Engine, path string) {

	// Controller
	workflowRoutes := &ROUTES.WorkflowsDB{DB: c.DB, MQ: &c.MQ}

	// Routes
	router := R.Group(path)
	router.GET("/wokflows", workflowRoutes.GetWorkflows)
	router.POST("/create/workflow", workflowRoutes.SaveWorkflows)
	router.POST("/create/aggregation", workflowRoutes.SaveAggregation)
	router.POST("/upload/template", workflowRoutes.UploadTemplate)
	router.GET("/start/template/:workflow_id", workflowRoutes.StartTemplate)
	router.GET("/exec/template/:workflow_id", workflowRoutes.ExecWorkflow)

	router.POST("test/email", func(ctx *gin.Context) {
		EMAIL.SendEmail("egresados398@gmail.com", "neldecas12@gmail.com", "Test", "Test")
		ctx.JSON(200, gin.H{"message": "Ok"})
	})
}
