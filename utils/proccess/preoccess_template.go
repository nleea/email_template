package proccess

import (
	"context"
	"encoding/json"
	"fmt"
	CO "sequency/config"
	M "sequency/models"
	PR "sequency/utils/mq"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type WorkflowsMQ struct {
	MQ *PR.ConnectionMQ
	DB *mongo.Database
}

type SendEmailStrunc struct {
	Email_info interface{} `bson:"email_info"`
	Data       any         `json:"data"`
}

func (c *WorkflowsMQ) ProcessTemplate(params M.ProcessParams) {

	envs := CO.ConfigEnv()
	collection := c.DB.Collection(envs["ATLAS_DB_AGGREGATION"])
	var resulst M.Aggregation
	err := collection.FindOne(context.TODO(), bson.M{"aggregation_name": params.Process.Aggregation_template}).Decode(&resulst)

	if err != nil {
		fmt.Println("error", err)
		return
	}

	collectionToAggregate := c.DB.Collection(resulst.Collection)

	switch params.Process.Type {
	case "email":

		if params.Process.Send_automatically != nil || params.Exec != nil {

			cursor, err := collectionToAggregate.Aggregate(context.TODO(), resulst.Aggregation)

			if err != nil {
				fmt.Println("error", err)
				return
			}

			var dataAggregation []SendEmailStrunc
			cursor.All(context.TODO(), &dataAggregation)

			resulstt, errM := json.Marshal(&dataAggregation)

			if errM != nil {
				fmt.Println("error", err)
				return
			}

			c.MQ.SendMessage(resulstt)

		} else {
			g := make([]M.ActionsWorkflow, 1)
			g[0] = params.Process
			SaveActions(envs["WORKFLOW_STATUS"], g, params.WorkflowId, params.StatusId, c)
		}

	case "decision":

		for i := range params.Process.Branches {
			agg := MakeAggregation(params.Process.Branches[i].Conditions)

			y := bson.D{{Key: "$match", Value: bson.D{{Key: "$and", Value: agg}}}}

			cursor, err := collectionToAggregate.Aggregate(context.TODO(), mongo.Pipeline{y})

			if err != nil {
				fmt.Println("error", err)
				return
			}

			var t []bson.M
			cursor.All(context.TODO(), &t)

			if t != nil {
				SaveActions(envs["WORKFLOW_STATUS"], params.Process.Branches[i].Actions, params.WorkflowId, params.StatusId, c)
				break
			}
		}
	}
}

func MakeAggregation(conditions []M.Conditions) []primitive.D {
	agg := make([]bson.D, len(conditions))
	for c := range conditions {
		f := bson.D{
			{Key: conditions[c].Field,
				Value: bson.D{{
					Key:   "$" + conditions[c].Condition,
					Value: conditions[c].Value,
				}}},
		}

		agg[c] = f
	}

	return agg
}

func SaveActions(collection string, actions []M.ActionsWorkflow, workflowId string, statusId interface{}, c *WorkflowsMQ) {
	for x := range actions {

		if actions[x].Type == "email" {
			collectionToSaveStatus := c.DB.Collection(collection)
			filter := bson.D{{Key: "_id", Value: statusId}}
			update := bson.D{{Key: "$push", Value: bson.D{{Key: "actions", Value: actions[x]}}}}
			collectionToSaveStatus.UpdateOne(context.TODO(), filter, update)
			continue
		}
		paramsTo := M.ProcessParams{Process: actions[x], WorkflowId: workflowId, StatusId: statusId}
		c.ProcessTemplate(paramsTo)
	}
}
