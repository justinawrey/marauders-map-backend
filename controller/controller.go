package controller

import (
	"gopkg.in/mgo.v2"
	"net/http"
	"encoding/json"
	"io/ioutil"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/cpen321/groupii-back/model"
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
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var location model.Location
	err = json.Unmarshal(body, &location)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = controller.Session.PutUserLoc(uuid, location)
	if err != nil {
		if err == mgo.ErrNotFound {
			http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}

func (controller *Controller) GetUserLoc(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	uuid := model.UUID(ps.ByName("uuid"))
	location := controller.Session.GetUserLoc(uuid)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(location)
}

func (controller *Controller) PutUser(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	body, _ := ioutil.ReadAll(req.Body)	
	var user model.User
	json.Unmarshal(body, &user)
	controller.Session.PutUser(user)
}

func (controller *Controller) GetUser(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	uuid := model.UUID(ps.ByName("uuid"))
	user := controller.Session.GetUser(uuid)
	w.Header().Set("Content-Type", "application/json")
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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(friends)
}

func (controller *Controller) GetAllUsers(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	users := controller.Session.GetAllUsers()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (controller *Controller) CleanUp() {
	controller.Session.CleanUp()
}
