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
	Longitude float64
	Latitude  float64
}

type User struct {
	UUID     UUID
	Name     string
	Email    string
	PhotoURL string
	Friends  []UUID
	Location Location
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
