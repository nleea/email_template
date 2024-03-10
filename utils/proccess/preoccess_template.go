package proccess

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	CO "sequency/config"
	DB "sequency/db"
	M "sequency/models"
	PR "sequency/utils/mq"
)

type WorkflowsMQ struct {
	MQ *PR.ConnectionMQ
}

func (c *WorkflowsMQ) ProcessTemplate(process M.ActionsWorkflow, workflowId string, statusId interface{}) {
	envs := CO.ConfigEnv()
	collection := DB.CLIENT_DB.Collection(envs["ATLAS_DB_AGGREGATION"])

	switch process.Type {
	case "email":

		if process.Send_automatically != nil {
			var resulst M.Aggregation
			err := collection.FindOne(context.TODO(), bson.M{"aggregation_name": process.Aggregation_template}).Decode(&resulst)

			if err != nil {
				fmt.Println("error", err)
				return
			}

			collectionToAggregate := DB.CLIENT_DB.Collection(resulst.Collection)
			cursor, err := collectionToAggregate.Aggregate(context.TODO(), resulst.Aggregation)

			if err != nil {
				fmt.Println("error", err)
				return
			}

			var dataAggregation interface{}
			cursor.All(context.TODO(), &dataAggregation)

			resulstt, errM := json.Marshal(&dataAggregation)

			if errM != nil {
				fmt.Println("error", err)
				return
			}

			c.MQ.SendMessage(resulstt)

		} else {
			collectionToSaveStatusE := DB.CLIENT_DB.Collection(envs["WORKFLOW_STATUS"])
			filter := bson.D{{Key: "_id", Value: statusId}}
			update := bson.D{{Key: "$addToSet", Value: bson.D{{Key: "actions", Value: process}}}}
			collectionToSaveStatusE.UpdateOne(context.TODO(), filter, update)
		}

	case "decision":

		for i := range process.Branches {
			agg := make([]bson.D, len(process.Branches[i].Conditions))
			for c := range process.Branches[i].Conditions {
				f := bson.D{
					{Key: process.Branches[i].Conditions[c].Field,
						Value: bson.D{{
							Key:   "$" + process.Branches[i].Conditions[c].Condition,
							Value: process.Branches[i].Conditions[c].Value,
						}}},
				}

				agg[c] = f
			}

			y := bson.D{{Key: "$match", Value: bson.D{{Key: "$and", Value: agg}}}}

			collectionTo := DB.CLIENT_DB.Collection(envs["ATLAS_DB_AGGREGATION"])
			var g M.Aggregation

			collectionTo.FindOne(context.TODO(), bson.M{"aggregation_name": process.Aggregation_template}).Decode(&g)
			collectionToAggregate := DB.CLIENT_DB.Collection(g.Collection)

			cursor, err := collectionToAggregate.Aggregate(context.TODO(), mongo.Pipeline{y})

			if err != nil {
				fmt.Println("error", err)
				return
			}

			var t []bson.M
			cursor.All(context.TODO(), &t)

			if t != nil {

				for x := range process.Branches[i].Actions {
					if process.Branches[i].Actions[x].Type == "email" {
						collectionToSaveStatus := DB.CLIENT_DB.Collection(envs["WORKFLOW_STATUS"])
						filter := bson.D{{Key: "_id", Value: statusId}}
						update := bson.D{{Key: "$addToSet", Value: bson.D{{Key: "actions", Value: process.Branches[i].Actions[i]}}}}
						collectionToSaveStatus.UpdateOne(context.TODO(), filter, update)
						continue
					}
					c.ProcessTemplate(process.Branches[i].Actions[x], workflowId, statusId)
				}
				break
			}

		}
	}
}
