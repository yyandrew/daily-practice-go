package tip

import (
	"context"
	"dailypractice/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Tip struct {
	Id       primitive.ObjectID `bson:"_id" json:"id"`
	Context  string             `json:"context"`
	ImageURL string             `json:"imageURL"`
	Category string             `json:"category"`
}

type Tipslice struct {
	Tips []Tip `json:"tips"`
}

func All() Tipslice {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	utils.CheckError(err)

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	connection := client.Database("daily-practice").Collection("tips")

	var s Tipslice

	cur, err := connection.Find(ctx, bson.D{})
	utils.CheckError((err))

	for cur.Next(ctx) {
		var result = Tip{}
		err := cur.Decode(&result)
		utils.CheckError(err)

		s.Tips = append(s.Tips, result)
	}

	return s
}
