package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type TextPost struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	Author    string        `json:"author"`
	Created   time.Time     `json:"created"`
	Edited    time.Time     `json:"edited"`
	Publish   time.Time     `json:"publish"` // Can set to publish in future
	DraftMode bool          `json:"draftmode"`
	Body      string        `json:"body"`
	Tags      []string      `json:"tags"`
	Views     int           `json:"views"`
}

func NewTextPost(author string, body string) *TextPost {
	return &TextPost{
		ID:        bson.NewObjectId(),
		Author:    author,
		Created:   time.Now(),
		DraftMode: true,
		Body:      body,
		Tags:      make([]string, 0),
		Views:     0,
	}
}
