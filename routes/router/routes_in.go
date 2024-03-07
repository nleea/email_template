package router

import (
	"context"
	"fmt"
	CO "sequency/config"
	DB_CONNECT "sequency/db"
	M "sequency/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetWorkflows(ctx *gin.Context) {

	envs := CO.ConfigEnv()

	collection := DB_CONNECT.CLIENT_DB.Collection(envs["ATLAS_DB_SEQUENCE"])

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

func SaveWorkflows(ctx *gin.Context) {
	envs := CO.ConfigEnv()

	collection := DB_CONNECT.CLIENT_DB.Collection(envs["ATLAS_DB_SEQUENCE"])

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

func SaveAggregation(ctx *gin.Context) {
	envs := CO.ConfigEnv()

	collection := DB_CONNECT.CLIENT_DB.Collection(envs["ATLAS_DB_AGGREGATION"])

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
