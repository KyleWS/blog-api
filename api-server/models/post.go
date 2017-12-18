package models

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type TextPost struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	Author    string        `json:"author"`
	Title     string        `json:"title"`
	Created   time.Time     `json:"created"`
	Edited    time.Time     `json:"edited"`
	Publish   time.Time     `json:"publish"` // Can set to publish in future
	DraftMode bool          `json:"draftmode"`
	Body      string        `json:"body"`
	Tags      []string      `json:"tags"`
	Views     int           `json:"views"`
}

// TextPostUpdates reflects the certain fields of a text post that are
// mutable.
type TextPostUpdates struct {
	Title     string    `json:"title"`
	Publish   time.Time `json:"publish"` // Can set to publish in future
	DraftMode bool      `json:"draftmode"`
	Body      string    `json:"body"`
	Tags      []string  `json:"tags"`
}

// PostShort is used to display all the posts without having to
// load the whole post body (which could be quite log)
type PostShort struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	Author    string        `json:"author"`
	Title     string        `json:"title"`
	Created   time.Time     `json:"created"`
	Edited    time.Time     `json:"edited"`
	Publish   time.Time     `json:"publish"` // Can set to publish in future
	DraftMode bool          `json:"draftmode"`
	Tags      []string      `json:"tags"`
	Views     int           `json:"views"`
}

func NewTextPost(author string, title string, body string) *TextPost {
	return &TextPost{
		ID:        bson.NewObjectId(),
		Author:    author,
		Title:     title,
		Created:   time.Now(),
		DraftMode: true,
		Body:      body,
		Tags:      make([]string, 0),
		Views:     0,
	}
}

func (tp *TextPost) ApplyUpdates(updates *TextPostUpdates) error {
	if len(updates.Body) > 0 {
		tp.Body = updates.Body
	} else {
		return fmt.Errorf("error cannot update to empty body")
	}
	tp.DraftMode = updates.DraftMode
	if !updates.Publish.IsZero() {
		tp.Publish = updates.Publish
	}
	if len(updates.Tags) > 0 {
		tp.Tags = updates.Tags
	} else {
		return fmt.Errorf("error cannot set tags to empty list")
	}
	if len(updates.Title) > 0 {
		tp.Title = updates.Title
	} else {
		return fmt.Errorf("error cannot set title to empty")
	}
	return nil
}
