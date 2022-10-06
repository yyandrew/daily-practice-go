package user

import (
	"context"
	"dailypractice/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Id               primitive.ObjectID `bson:"_id" json:"id"`
	Email            string             `json:"email"`
	EcryptedPassword string             `bson:"cryptedPassword" json:"ecryptedPassword"`
}

type Userslice struct {
	Users []User `json:"users"`
}

func All() Userslice {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	utils.CheckError(err)

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	connection := client.Database("daily-practice").Collection("users")
	users := Userslice{}

	cur, err := connection.Find(ctx, bson.D{})
	utils.CheckError((err))
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var result User
		err := cur.Decode(&result)
		utils.CheckError(err)

		users.Users = append(users.Users, result)
	}

	return users
}
