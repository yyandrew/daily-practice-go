package user

import (
	"context"
	"dailypractice/utils"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id               primitive.ObjectID `bson:"_id" json:"id"`
	Email            string             `json:"email"`
	EcryptedPassword string             `bson:"cryptedPassword"`
}

type Userslice struct {
	Users []User `json:"users"`
}

var (
	ctx        context.Context
	cancel     context.CancelFunc
	collection *mongo.Collection
)

func init() {
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	utils.CheckError(err)

	collection = client.Database("daily-practice").Collection("users")
}

func All() Userslice {
	users := Userslice{}

	cur, err := collection.Find(ctx, bson.D{})
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

func FindByEmail(email string) (User, error) {
	filter := bson.D{{"email", email}}
	user := User{}
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		fmt.Println(err)
	}
	return user, err
}

func (u *User) AuthByPassword(plainPW string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.EcryptedPassword), []byte(plainPW))
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
