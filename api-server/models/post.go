package models

import (
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

type UserTextPost struct {
	Author    string    `json:"author"`
	Title     string    `json:"title"`
	Publish   time.Time `json:"publish"` // Can set to publish in future
	DraftMode bool      `json:"draftmode"`
	Body      string    `json:"body"`
	Tags      []string  `json:"tags"`
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

func (utp *UserTextPost) GenPostMetaData() *TextPost {
	newPost := NewTextPost(utp.Author, utp.Title, utp.Body)
	newPost.Publish = utp.Publish
	newPost.DraftMode = utp.DraftMode
	if utp.Tags != nil {
		newPost.Tags = utp.Tags
	}
	return newPost
}

func (tp *TextPost) ApplyUpdates(updates *TextPostUpdates) error {
	if len(updates.Body) > 0 {
		tp.Body = updates.Body
	}
	tp.DraftMode = updates.DraftMode
	if !updates.Publish.IsZero() {
		tp.Publish = updates.Publish
	}
	if len(updates.Tags) > 0 {
		tp.Tags = updates.Tags
	}
	if len(updates.Title) > 0 {
		tp.Title = updates.Title
	}
	return nil
}
