package mongo

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	_database   = "test_db"
	_uri        = "mongodb://127.0.0.1:27017/admin"
	_collection = "test"
)

type doc struct {
	Name  string `bson:"name"`
	Value string `bson:"value"`
}

func TestDatabase_UpsertCollectionIndexMany(t *testing.T) {

	db, err := NewDatabase(context.Background(), _uri, _database)
	if err != nil {
		t.Fatal(err)
	}

	indexs := []Index{
		{
			Collection: _collection,
			Keys:       bson.D{{Key: "name", Value: 1}},
		},
		{
			Collection: _collection,
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

	cursor, err := db.Collection(_collection).Indexes().List(nil)
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
