package model

import (
   "gopkg.in/mgo.v2/bson"
   mgo "gopkg.in/mgo.v2"
   "fmt"
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


func (ms *MongoStore) GetByID(id bson.ObjectId) (*TextPost, error) {
	result := &TextPost{}
	col := ms.session.DB(ms.dbname).C(ms.colname)
	if err := col.Find(bson.M{"_id": id}).One(&result); err != nil {
		return nil, fmt.Errof("error finding user: %v", err)
	}
	return result, nil
}
