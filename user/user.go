package user

import (
	"context"
	// "dailypractice/tip"
	"dailypractice/tip"
	"dailypractice/utils"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id              MyObjectID `bson:"_id" json:"id"`
	Email           string     `json:"email"`
	CryptedPassword string     `bson:"cryptedPassword" json:"-"`
}
type MyObjectID string

func (id MyObjectID) MarshalBSONValue() (bsontype.Type, []byte, error) {
	p, err := primitive.ObjectIDFromHex(string(id))
	if err != nil {
		return bsontype.Null, nil, err
	}

	return bson.MarshalValue(p)
}

type Userslice struct {
	Users []User `json:"users"`
}

func getCollection(ctx context.Context) *mongo.Collection {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	utils.CheckError(err)

	collection := client.Database("daily-practice").Collection("users")
	return collection
}

func All() Userslice {
	users := Userslice{}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := getCollection(ctx)
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

func (u *User) Tips(category string, content string) (interface{}, error) {
	fmt.Println(u.Id)
	tips, err := tip.FilterByUserId(string(u.Id), category, content)
	if err != nil {
		fmt.Println(err)
		return tips, err
	}

	return tips, nil
}

func FindByEmail(email string) (User, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	filter := bson.D{{"email", email}}
	user := User{}
	collection := getCollection(ctx)
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		fmt.Println(err)
	}
	return user, err
}

func FindById(id string) (interface{}, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	objectId, err := primitive.ObjectIDFromHex(id)
	user := User{}
	if err != nil {
		return nil, err
	}

	filter := bson.D{{"_id", objectId}}
	collection := getCollection(ctx)
	err1 := collection.FindOne(ctx, filter).Decode(&user)
	if err1 != nil {
		fmt.Println(err)
		return nil, err
	}
	return user, err
}

func (u *User) AuthByPassword(plainPW string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.CryptedPassword), []byte(plainPW))
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func Save(email string, password string) (interface{}, bool) {
	user := User{}

	cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return nil, false
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := getCollection(ctx)
	res, err := collection.InsertOne(ctx, bson.D{{"email", email}, {"cryptedPassword", string(cryptedPassword)}})
	if err != nil {
		fmt.Println(err)
		return nil, false
	}

	collection.FindOne(ctx, bson.D{{"_id", res.InsertedID}}).Decode(&user)

	return user, true
}
