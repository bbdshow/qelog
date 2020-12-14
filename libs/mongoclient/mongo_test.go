package mongoclient

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	uri        = "mongodb://localhost:27017"
	database   = "mongo_test"
	collection = "test"
)

type doc struct {
	Name  string `bson:"name"`
	Value string `bson:"value"`
}

func TestNewMongoClientByURI(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := NewMongoClientByURI(ctx, uri)
	if err != nil {
		t.Fatal(err)
	}

	i := rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10000000)
	coll := client.Database(database).Collection(collection)
	doc := doc{
		Name:  "TestNewMongoClientByURI",
		Value: strconv.Itoa(int(i)),
	}

	result, err := coll.InsertOne(ctx, doc)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(result.InsertedID)
}

func TestCreateCollectionIndexMany(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := NewMongoClientByURI(ctx, uri)
	if err != nil {
		t.Fatal(err)
	}
	database := client.Database(database)

	many := []Index{
		{
			Collection: collection,
			Keys: bson.M{
				"name": 1,
			},
		},
		{
			Collection: collection,
			Name:       "name_value",
			Keys: bson.M{
				"value": -1,
				"name":  1,
			},
			Unique:     true,
			Background: true,
		}}

	err = UpsertCollectionIndexMany(database, many)
	if err != nil {
		t.Fatal(err)
		return
	}

	cursor, err := database.Collection(collection).Indexes().List(ctx)
	if err != nil {
		t.Fatal(err)
		return
	}
	var value interface{}
next:
	if cursor.Next(ctx) {
		err = cursor.Decode(&value)
		if err != nil {
			t.Log(err) //return
		}

		if strings.Index(fmt.Sprintf("%#v", value), "name_value") != -1 {
			return
		}
		goto next
	}

	t.Fail()
}
