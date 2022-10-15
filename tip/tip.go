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
	UserId   primitive.ObjectID `bson:"user_id" json:"user_id"`
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

func Create(content string, category string, imageUrl string, user_id string) (interface{}, error) {
	newTip := Tip{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := getCollection(ctx)
	res, err := collection.InsertOne(ctx, bson.D{{"context", content}, {"category", category}, {"imageUrl", imageUrl}, {"user_id", user_id}})

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Printf("res %v", res.InsertedID)

	collection.FindOne(ctx, bson.D{{"_id", res.InsertedID}}).Decode(&newTip)
	fmt.Printf("new tip: %v", newTip)

	return newTip, nil
}

func FilterByUserId(user_id string, category string, content string) (Tipslice, error) {
	s := Tipslice{Tips: make([]Tip, 0)}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := getCollection(ctx)
	filter := bson.D{{"user_id", user_id}}
	if category != "" {
		filter = append(filter, bson.E{"category", category})
	}
	if content != "" {
		filter = append(filter, bson.E{"context", primitive.Regex{Pattern: content, Options: ""}})
	}
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return s, err
	}
	for cur.Next(ctx) {
		var result = Tip{}
		err := cur.Decode(&result)
		if err != nil {
			fmt.Println(err)
			continue
		}
		s.Tips = append([]Tip{result}, s.Tips...)
	}
	return s, nil
}
