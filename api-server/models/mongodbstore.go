package models

import (
	"fmt"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoStore struct {
	session *mgo.Session
	dbname  string
	colname string
}

func NewMongoStore(sess *mgo.Session, dbName string, collectionName string) *MongoStore {
	if sess == nil {
		panic("nil pointer passed for session")
	}
	return &MongoStore{
		session: sess,
		dbname:  dbName,
		colname: collectionName,
	}
}

// GetTextPostByID returns a TextPost struct for the post with the
// provided ID.
func (ms *MongoStore) GetTextPostByID(id bson.ObjectId) (*TextPost, error) {
	result := &TextPost{}
	col := ms.session.DB(ms.dbname).C(ms.colname)
	if err := col.Find(bson.M{"_id": id}).One(&result); err != nil {
		return nil, fmt.Errorf("error finding user: %v", err)
	}
	return result, nil
}

// Insert writes given post to database.
func (ms *MongoStore) InsertTextPost(newPost *TextPost) error {
	col := ms.session.DB(ms.dbname).C(ms.colname)
	if err := col.Insert(newPost); err != nil {
		return fmt.Errorf("error inserting new post to mongodb: %v", err)
	}
	return nil
}

// FetchAllShort returns slice of all posts excluding their
// body field.
func (ms *MongoStore) FetchAllShort() ([]*PostShort, error) {
	longPost := &TextPost{}
	shortSlice := make([]*PostShort, 0)
	col := ms.session.DB(ms.dbname).C(ms.colname)
	iterVal := col.Find(bson.M{}).Iter()
	for iterVal.Next(longPost) {
		postShort := &PostShort{
			ID:        longPost.ID,
			Author:    longPost.Author,
			Title:     longPost.Title,
			Created:   longPost.Created,
			Edited:    longPost.Edited,
			Publish:   longPost.Publish,
			DraftMode: longPost.DraftMode,
			Tags:      longPost.Tags,
			Views:     longPost.Views,
		}
		shortSlice = append(shortSlice, postShort)
	}
	if err := iterVal.Err(); err != nil {
		return nil, err
	}
	return shortSlice, nil
}

// DeleteTextPost will delete post with given ID
func (ms *MongoStore) DeletePost(postID bson.ObjectId) error {
	col := ms.session.DB(ms.dbname).C(ms.colname)
	if err := col.RemoveId(postID); err != nil {
		return fmt.Errorf("error deleting post: %v", err)
	}
	return nil
}

func (ms *MongoStore) UpdateTextPost(postID bson.ObjectId, updates *TextPostUpdates) (*TextPost, error) {
	postToUpdate, err := ms.GetTextPostByID(postID)
	if err != nil {
		return nil, err
	}
	if err := postToUpdate.ApplyUpdates(updates); err != nil {
		return nil, fmt.Errorf("error applying updates to post: %v", err)
	}

	change := mgo.Change{
		Update:    bson.M{"$set": postToUpdate},
		ReturnNew: true,
	}
	// TODO: NEED TO TEST IF THIS ERASES POST BODY
	result := &TextPost{}
	col := ms.session.DB(ms.dbname).C(ms.colname)
	if _, err := col.FindId(postID).Apply(change, result); err != nil {
		return nil, fmt.Errorf("error updating record: %v", err)
	}
	return result, nil
}

// TODO: Backup database often and export "off-site"
