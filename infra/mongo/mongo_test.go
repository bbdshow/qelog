package mongo

import (
	"context"
	"fmt"
	"strings"
	"testing"

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

func TestDatabase_UpsertCollectionIndexMany(t *testing.T) {

	db, err := NewDatabase(context.Background(), uri, database)
	if err != nil {
		t.Fatal(err)
	}

	indexs := []Index{
		{
			Collection: collection,
			Keys:       bson.D{{Key: "name", Value: 1}},
		},
		{
			Collection: collection,
			Name:       "name_value",
			Keys:       bson.D{{Key: "name", Value: 1}, {Key: "value", Value: -1}},
			Unique:     true,
			Background: true,
		}}

	err = db.UpsertCollectionIndexMany(indexs)
	if err != nil {
		t.Fatal(err)
		return
	}

	cursor, err := db.Collection(collection).Indexes().List(nil)
	if err != nil {
		t.Fatal(err)
		return
	}
	var value interface{}
next:
	if cursor.Next(nil) {
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
