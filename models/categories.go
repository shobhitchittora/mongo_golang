package main

import "gopkg.in/mgo.v2/bson"

type categories struct {
	ID    bson.ObjectId `json:"id" bson:"_id"`
	Title string        `json:"title" bson:"title"`
	Links []string      `json:"links" bson:"links"`
}
