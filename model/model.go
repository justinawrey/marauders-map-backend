package model

import (
	"os"

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
	// check if we are running through heroku or localhost
	var mgoSession *MgoSession
	if mongoDbURI := os.Getenv("MONGODB_URI"); mongoDbURI != "" {
		mgoSession = dialHeroku(mongoDbURI)
	} else {
		mgoSession = dialLocalDB()
	}
	return mgoSession
}

func dialHeroku(mongoDbURI string) *MgoSession {
	dialInfo, _ := mgo.ParseURL(mongoDbURI)
	session, _ := mgo.DialWithInfo(dialInfo)
	db := session.DB(dialInfo.Database)
	collection := db.C("users")
	return &MgoSession{
		CurrSession:    session,
		CurrDB:         db,
		CurrCollection: collection,
	}
}

func dialLocalDB() *MgoSession {
	session, _ := mgo.Dial("localhost")
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

func (mongoSession *MgoSession) PutUser(user User) error {
	_, err := mongoSession.CurrCollection.Upsert(bson.M{"uuid": user.UUID}, user)
	return err
}

func (mongoSession *MgoSession) GetUser(id UUID) (User, error) {
	var user User
	err := mongoSession.CurrCollection.Find(bson.M{"uuid": id}).One(&user)
	return user, err
}

func (mongoSession *MgoSession) DeleteUser(id UUID) error {
	return mongoSession.CurrCollection.Remove(bson.M{"uuid": id})
}

func (mongoSession *MgoSession) PutUserLoc(id UUID, location Location) error {
	toUpdate := bson.M{"uuid": id}
	update := bson.M{"$set": bson.M{"location": location}}
	return mongoSession.CurrCollection.Update(toUpdate, update)
}

func (mongoSession *MgoSession) GetUserLoc(id UUID) (Location, error) {
	user, err := mongoSession.GetUser(id)
	return user.Location, err
}

func (mongoSession *MgoSession) PutFriend(id UUID, friendId UUID) error {
	toUpdate := bson.M{"uuid": id}
	update := bson.M{"$addToSet": bson.M{"friends": friendId}}
	return mongoSession.CurrCollection.Update(toUpdate, update)
}

func (mongoSession *MgoSession) DeleteFriend(id UUID, friendId UUID) error {
	toUpdate := bson.M{"uuid": id}
	update := bson.M{"$pull": bson.M{"friends": friendId}}
	return mongoSession.CurrCollection.Update(toUpdate, update)
}

func (mongoSession *MgoSession) GetFriends(id UUID) ([]User, error) {
	user, err := mongoSession.GetUser(id)
	if err != nil {
		return nil, err
	}
	friendIds := user.Friends
	var friends []User
	err = mongoSession.CurrCollection.Find(bson.M{"uuid": bson.M{"$in": friendIds}}).All(&friends)
	return friends, err
}

func (mongoSession *MgoSession) GetAllUsers() ([]User, error) {
	var users []User
	err := mongoSession.CurrCollection.Find(nil).All(&users)
	return users, err
}

func (mongoSession *MgoSession) CleanUp() {
	mongoSession.CurrSession.Close()
}
