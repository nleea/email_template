package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	CO "sequency/config"
)

func ConnectDB() *mongo.Database {

	envs := CO.ConfigEnv()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(envs["ATLAS_URI"]).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)
	fmt.Println(envs["ATLAS_URI"])
	if err != nil {
		panic(err)
	}
	// defer func() {
	// 	if err = client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()
	var result bson.M
	if err := client.Database(envs["ATLAS_DB"]).RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	return client.Database(envs["ATLAS_DB"])
}
