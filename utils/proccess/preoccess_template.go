package proccess

import (
	"context"
	"fmt"
	CO "sequency/config"
	DB "sequency/db"
	M "sequency/models"
	// NSQ "sequency/utils/nsq"
	// "time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func ProcessTemplate(process M.ActionsWorkflow, workflowId string, statusId interface{}) {
	envs := CO.ConfigEnv()
	collection := DB.CLIENT_DB.Collection(envs["ATLAS_DB_AGGREGATION"])

	switch process.Type {
	case "email":

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

		// NSQ.SendMessageToNSQ("emails", M.MessageNSQ{Name: "Test", Content: "TEst", Timestamp: time.Now().String()})

	case "decision":

		for i := range process.Branches {
			agg := []bson.A{}
			for c := range process.Branches[i].Conditions {
				f := bson.A{
					bson.D{
						{Key: process.Branches[i].Conditions[c].Field,
							Value: bson.D{{
								Key:   "$" + process.Branches[i].Conditions[c].Condition,
								Value: process.Branches[i].Conditions[c].Value,
							}}},
					},
				}

				agg = append(agg, f)
			}

			y := bson.D{{Key: "$match", Value: bson.D{{Key: "$and", Value: agg[0]}}}}


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
			fmt.Println("desicion", t)

		}

		// _, err := collection.Aggregate(context.TODO(), bson.M{})

		// if err != nil {
		// 	fmt.Println("error", err)
		// 	return
		// }
	}
}
