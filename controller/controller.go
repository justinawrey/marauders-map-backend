package controller

import (
	"net/http"
	"encoding/json"
	"io/ioutil"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/justinawrey/groupii-back/model"
)

type Controller struct {
	Session *model.MgoSession
}

func New() *Controller {
	return &Controller{
		Session: model.New(),
	}
}

func (controller *Controller) PutUserLoc(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	uuid := model.UUID(ps.ByName("uuid"))
	body, _ := ioutil.ReadAll(req.Body)
	var location model.Location
	json.Unmarshal(body, &location)
	controller.Session.PutUserLoc(uuid, location)
}

func (controller *Controller) GetUserLoc(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	uuid := model.UUID(ps.ByName("uuid"))
	location := controller.Session.GetUserLoc(uuid)
	json.NewEncoder(w).Encode(location)
}

func (controller *Controller) PutUser(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	body, _ := ioutil.ReadAll(req.Body)	
	var user model.User
	json.Unmarshal(body, &user)
	controller.Session.PutUser(user)
}

func (controller *Controller) GetUser(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	uuid := model.UUID(ps.ByName("uuid"))
	user := controller.Session.GetUser(uuid)
	json.NewEncoder(w).Encode(user)
}

func (controller *Controller) DeleteUser(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	uuid := model.UUID(ps.ByName("uuid"))
	controller.Session.DeleteUser(uuid)
}

func (controller *Controller) PutFriend(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	uuid := model.UUID(ps.ByName("uuid"))
	friendId := model.UUID(ps.ByName("friendid"))
	controller.Session.PutFriend(uuid, friendId)
}

func (controller *Controller) DeleteFriend(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	uuid := model.UUID(ps.ByName("uuid"))
	friendId := model.UUID(ps.ByName("friendid"))
	controller.Session.DeleteFriend(uuid, friendId)
}

func (controller *Controller) GetFriends(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	uuid := model.UUID(ps.ByName("uuid"))
	friends := controller.Session.GetFriends(uuid)
	json.NewEncoder(w).Encode(friends)
}

func (controller *Controller) CleanUp() {
	controller.Session.CleanUp()
}
