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
	location, err := controller.Session.GetUserLoc(uuid)
	if err != nil {
		if err == mgo.ErrNotFound {
			http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(location)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (controller *Controller) PutUser(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	body, err := ioutil.ReadAll(req.Body)	
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var user model.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = controller.Session.PutUser(user)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (controller *Controller) GetUser(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	uuid := model.UUID(ps.ByName("uuid"))
	user, err := controller.Session.GetUser(uuid)
	if err != nil {
		if err == mgo.ErrNotFound {
			http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
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
