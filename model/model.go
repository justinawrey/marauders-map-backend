package model

import (
	"gopkg.in/mgo.v2"
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
	UUID     UUID `json:"uuid" bson:"uuid"`
	Name     string `json:"name" bson:"name"`
	Email    string `json:"email" bson:"email"`
	PhotoURL string `json:"photoURL" bson:"photoURL"`
	Friends  []UUID `json:"friends" bson:"friends"`
	Location Location `json:"location" bson:"location"`
}

type UUID string

func New() *MgoSession {
	return &MgoSession{
	// initialize a handler to default session
	// initialize a handler to default DB
	// initialize a handler to default collection
	}
}

func (mongoSession *MgoSession) SwitchDB() {
}

func (mongoSession *MgoSession) SwitchCollection() {
}

func (mongoSession *MgoSession) PutUser(user User) {
}

func (mongoSession *MgoSession) GetUser(id UUID) User {
	return User{}
}

func (mongoSession *MgoSession) DeleteUser(id UUID) {
}

func (mongoSession *MgoSession) PutUserLoc(id UUID, location Location) {
}

func (mongoSession *MgoSession) GetUserLoc(id UUID) Location {
	return Location{}
}

func (mongoSession *MgoSession) PutFriend(id UUID, friendId UUID) {
}

func (mongoSession *MgoSession) DeleteFriend(id UUID, friendId UUID) {
}

func (mongoSession *MgoSession) GetFriends(id UUID) []User {
	return []User{}
}

func (mongoSession *MgoSession) CleanUp() {
}
