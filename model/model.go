package model

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MgoSession struct {
	CurrSession    *mgo.Session
	CurrDB         *mgo.Database
	CurrCollection *mgo.Collection
}

type Location struct {
	Longitude float64 `json:"longitude" bson:"longitude"`
	Latitude  float64 `json:"latitude" bson:"latitude"`
}

type User struct {
	UUID     UUID     `json:"uuid" bson:"uuid"`
	Name     string   `json:"name" bson:"name"`
	Email    string   `json:"email" bson:"email"`
	PhotoURL string   `json:"photoURL" bson:"photoURL"`
	Friends  []UUID   `json:"friends" bson:"friends"`
	Location Location `json:"location" bson:"location"`
}

type UUID string

func New() *MgoSession {
	session, _ := mgo.Dial("mongodb://root:root@localhost:27017")
	db := session.DB("marauders")
	collection := db.C("users")
	return &MgoSession{
		CurrSession:    session,
		CurrDB:         db,
		CurrCollection: collection,
	}
}

func (mongoSession *MgoSession) SwitchDB(db string) {
	mongoSession.CurrDB = mongoSession.CurrSession.DB(db)
}

func (mongoSession *MgoSession) SwitchCollection(collection string) {
	mongoSession.CurrCollection = mongoSession.CurrDB.C(collection)
}

func (mongoSession *MgoSession) PutUser(user User) {
	mongoSession.CurrCollection.Insert(user)
}

func (mongoSession *MgoSession) GetUser(id UUID) User {
	var user User
	mongoSession.CurrCollection.Find(bson.M{"uuid": id}).One(&user)
	return user
}

func (mongoSession *MgoSession) DeleteUser(id UUID) {
	mongoSession.CurrCollection.Remove(bson.M{"uuid": id})
}

func (mongoSession *MgoSession) PutUserLoc(id UUID, location Location) {
	toUpdate := bson.M{"uuid": id}
	update := bson.M{"$set": bson.M{"longitude": location.Longitude, "latitude": location.Latitude}}
	mongoSession.CurrCollection.Update(toUpdate, update)
}

func (mongoSession *MgoSession) GetUserLoc(id UUID) Location {
	var location Location
	mongoSession.CurrCollection.Find(bson.M{"uuid": id}).Select(bson.M{"location": 1}).One(&location)
	return location
}

func (mongoSession *MgoSession) PutFriend(id UUID, friendId UUID) {
	toUpdate := bson.M{"uuid": id}
	update := bson.M{"$addToSet": bson.M{"friends": friendId}}
	mongoSession.CurrCollection.Update(toUpdate, update)
}

func (mongoSession *MgoSession) DeleteFriend(id UUID, friendId UUID) {
	toUpdate := bson.M{"uuid": id}
	update := bson.M{"$pull": bson.M{"friends": friendId}}
	mongoSession.CurrCollection.Update(toUpdate, update)
}

func (mongoSession *MgoSession) GetFriends(id UUID) []User {
	var friendIds []UUID
	mongoSession.CurrCollection.Find(bson.M{"uuid": id}).Select(bson.M{"friends": 1}).One(&friendIds)
	var friends []User
	mongoSession.CurrCollection.Find(bson.M{"uuid": bson.M{"$in": friendIds}}).All(&friends)
	return friends
}

func (mongoSession *MgoSession) CleanUp() {
	mongoSession.CurrSession.Close()
}
