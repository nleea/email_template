package router

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	CO "sequency/config"
	M "sequency/models"
	PR "sequency/utils/mq"
	UTP "sequency/utils/proccess"

	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type WorkflowsDB struct {
	DB *mongo.Database
	MQ *PR.ConnectionMQ
}

func (c *WorkflowsDB) GetWorkflows(ctx *gin.Context) {

	envs := CO.ConfigEnv()

	collection := c.DB.Collection(envs["emails"])

	cursor, err := collection.Find(context.TODO(), bson.D{})

	if err != nil {
		fmt.Println("There was error", err)
		ctx.JSON(400, gin.H{"message": err})
		return
	}

	var resulst []M.Workflows

	err2 := cursor.All(context.TODO(), &resulst)
	if err2 != nil {
		fmt.Println("There was a error", err2)
	}

	ctx.JSON(200, gin.H{"data": resulst})
}

func (c *WorkflowsDB) SaveWorkflows(ctx *gin.Context) {
	envs := CO.ConfigEnv()

	collection := c.DB.Collection(envs["ATLAS_DB_SEQUENCE"])

	newDocument := M.Workflows{ID: primitive.NewObjectID()}
	err := ctx.BindJSON(&newDocument)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err})
		return
	}

	_, erri := collection.InsertOne(context.TODO(), newDocument)

	if erri != nil {
		ctx.JSON(400, gin.H{"message": err})
		return
	}

	ctx.JSON(200, gin.H{"message": "Ok"})
}

func (c *WorkflowsDB) SaveAggregation(ctx *gin.Context) {
	envs := CO.ConfigEnv()

	collection := c.DB.Collection(envs["ATLAS_DB_AGGREGATION"])

	newDocument := M.Aggregation{}
	err := ctx.BindJSON(&newDocument)

	if err != nil {
		ctx.JSON(400, gin.H{"message": err})
		return
	}

	_, erri := collection.InsertOne(context.TODO(), newDocument)

	if erri != nil {
		ctx.JSON(400, gin.H{"message": err})
		return
	}

	ctx.JSON(200, gin.H{"message": "Ok"})
}

func (c *WorkflowsDB) UploadTemplate(ctx *gin.Context) {
	file, err := ctx.FormFile("file")

	if err != nil {
		ctx.JSON(400, gin.H{"error": err})
		return
	}

	path, errgwd := os.Getwd()

	if errgwd != nil {
		ctx.JSON(400, gin.H{"error": errgwd})
		return
	}

	filePath := filepath.Join(path, "templates", file.Filename)

	errupload := ctx.SaveUploadedFile(file, filePath)

	if errupload != nil {
		fmt.Println("There was error uploading the file", errupload)
		ctx.JSON(400, gin.H{"error": errupload.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Ok"})
}

func (c *WorkflowsDB) StartTemplate(ctx *gin.Context) {
	envs := CO.ConfigEnv()
	collection := c.DB.Collection(envs["WORKFLOW_STATUS"])

	workflowID := ctx.Param("workflow_id")

	sequence := c.DB.Collection(envs["ATLAS_DB_SEQUENCE"])

	var wokflows M.Workflows

	id, errWorkId := primitive.ObjectIDFromHex(workflowID)

	if errWorkId != nil {
		ctx.JSON(400, gin.H{"error": errWorkId.Error()})
		return
	}

	errFind := sequence.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&wokflows)

	if errFind != nil {
		ctx.JSON(400, gin.H{"error": errFind.Error()})
		return
	}

	resulst, err := collection.InsertOne(context.TODO(), M.WorkflowStatus{Workflow: workflowID, Actions: []M.ActionsWorkflow{},
		History: []M.WorkflowHistory{}, Next_action: "", Timestamp: ""})

	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	rt := UTP.WorkflowsMQ{MQ: c.MQ, DB: c.DB}

	for i := range wokflows.Actions {
		params := M.ProcessParams{WorkflowId: workflowID, Process: wokflows.Actions[i], StatusId: resulst.InsertedID}
		rt.ProcessTemplate(params)
	}

	ctx.JSON(200, gin.H{"message": "Ok"})
}

func (c *WorkflowsDB) ExecWorkflow(ctx *gin.Context) {
	envs := CO.ConfigEnv()
	collection := c.DB.Collection(envs["WORKFLOW_STATUS"])

	workflowId := ctx.Param("workflow_id")

	var data M.WorkflowStatus
	work, _ := primitive.ObjectIDFromHex(workflowId)
	fmt.Println(collection.Name())
	err := collection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: work}}).Decode(&data)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if data.Timestamp == "" && len(data.Actions) > 1 {

		// action := data.History[0]
		// nextAction := data.Actions[0]
		t := time.Now().UTC()
		// NextActionTime := t.Format("1h")
		fmt.Println(t)
		// _, err := collection.UpdateOne(context.TODO(), bson.M{"_id": work}, bson.D{{Key: "$set", Value: bson.A{bson.D{{Key: "next_action", Value: nextAction.ID}}, bson.D{{Key: "timestamp", Value: NextActionTime}}}}})

		// if err != nil {
		// 	ctx.JSON(400, gin.H{"error": err.Error()})
		// 	return
		// }
	}
	ctx.JSON(200, gin.H{"message": "Ok"})

}
