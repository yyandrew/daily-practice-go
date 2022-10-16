package tip

import (
	"context"
	"dailypractice/utils"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Tip struct {
	Id        primitive.ObjectID `bson:"_id" json:"_id"`
	Context   string             `json:"context"`
	ImageURL  string             `json:"imageUrl"`
	Category  string             `json:"category"`
	UserId    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Deletable bool               `json:"deletable"`
}

type Tipslice struct {
	Tips []Tip `json:"tips"`
}

var IMG_PATH string

func init() {
	err := godotenv.Load()
	utils.CheckError(err)
	IMG_PATH = os.Getenv("IMG_PATH")
}

func getCollection(ctx context.Context) *mongo.Collection {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	utils.CheckError(err)

	collection := client.Database("daily-practice").Collection("tips")
	return collection
}

func All(category string, content string, user_id string) Tipslice {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := getCollection(ctx)

	s := Tipslice{Tips: make([]Tip, 0)}

	filter := bson.D{{"category", category}}
	if content != "" {
		filter = append(filter, bson.E{"context", content})
	}
	cur, err := collection.Find(ctx, filter)
	utils.CheckError((err))

	for cur.Next(ctx) {
		var result = Tip{}
		err := cur.Decode(&result)
		utils.CheckError(err)
		if user_id != "" && result.UserId.Hex() == user_id {
			result.Deletable = true
		} else {
			result.Deletable = false
		}

		s.Tips = append([]Tip{result}, s.Tips...)
	}

	return s
}

func Delete(id string) (Tip, bool) {
	ok := true
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	deletedTip := Tip{}
	collection := getCollection(ctx)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return deletedTip, false
	}

	collection.FindOneAndDelete(ctx, bson.D{{"_id", objectId}}).Decode(&deletedTip)
	fmt.Printf("image!!! %s\n", IMG_PATH+deletedTip.ImageURL)
	if deletedTip.ImageURL != "" {
		err = os.Remove(IMG_PATH + deletedTip.ImageURL)
		if err != nil {
			fmt.Errorf(err.Error())
			return deletedTip, false
		}
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

func FindById(id string) (Tip, bool) {
	ok := true
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tip := Tip{}

	objectId, err := primitive.ObjectIDFromHex(id)
	utils.CheckError(err)
	filter := bson.D{{"_id", objectId}}
	collection := getCollection(ctx)
	collection.FindOne(ctx, filter).Decode(&tip)

	return tip, ok
}
