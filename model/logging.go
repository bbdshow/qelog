package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Logging struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	UniqueKey     string             `bson:"uk"`
	Module        string             `bson:"m"`
	IP            string             `bson:"ip"`
	Level         int                `bson:"l"`
	Short         string             `bson:"s"`
	Full          string             `json:"f"`
	Condition1    string             `bson:"c1"`
	Condition2    string             `bson:"c2"`
	Condition3    string             `bson:"c3"`
	MillTimeStamp int64              `bson:"t"` // 毫秒
}
