package tip

import (
	"context"
	"dailypractice/utils"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Tip struct {
	Id       primitive.ObjectID `bson:"_id" json:"_id"`
	Context  string             `json:"context"`
	ImageURL string             `json:"imageUrl"`
	Category string             `json:"category"`
}

type Tipslice struct {
	Tips []Tip `json:"tips"`
}

func getCollection(ctx context.Context) *mongo.Collection {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	utils.CheckError(err)

	collection := client.Database("daily-practice").Collection("tips")
	return collection
}

func All(category string) Tipslice {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := getCollection(ctx)

	s := Tipslice{Tips: make([]Tip, 0)}

	cur, err := collection.Find(ctx, bson.D{{"category", category}})
	utils.CheckError((err))

	for cur.Next(ctx) {
		var result = Tip{}
		err := cur.Decode(&result)
		utils.CheckError(err)

		s.Tips = append([]Tip{result}, s.Tips...)
	}

	return s
}

func Delete(id string) (Tip, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	deletedTip := Tip{}
  fmt.Printf("id: %s\n", id)

	objectId, err := primitive.ObjectIDFromHex(id)
	utils.CheckError(err)
	filter := bson.D{{"_id", objectId}}
	collection := getCollection(ctx)
	collection.FindOne(ctx, filter).Decode(&deletedTip)
	utils.CheckError(err)

	res, err := collection.DeleteOne(ctx, filter)
	utils.CheckError(err)
	fmt.Printf("res: %+v", res)
	ok := true
	if err != nil {
		ok = false
	}

	return deletedTip, ok
}
