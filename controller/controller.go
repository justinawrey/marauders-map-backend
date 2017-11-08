package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"gopkg.in/mgo.v2"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/cpen321/groupii-back/model"
)

type Controller struct {
	Session *model.MgoSession
}

func NewController() *Controller {
	return &Controller{
		Session: model.NewSession(),
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
		checkForResourceNotFound(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (controller *Controller) GetUserLoc(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	uuid := model.UUID(ps.ByName("uuid"))
	location, err := controller.Session.GetUserLoc(uuid)
	if err != nil {
		checkForResourceNotFound(w, err)
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
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (controller *Controller) GetUser(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	uuid := model.UUID(ps.ByName("uuid"))
	user, err := controller.Session.GetUser(uuid)
	if err != nil {
		checkForResourceNotFound(w, err)
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
	err := controller.Session.DeleteUser(uuid)
	if err != nil {
		checkForResourceNotFound(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (controller *Controller) PutFriend(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	uuid := model.UUID(ps.ByName("uuid"))
	friendId := model.UUID(ps.ByName("friendid"))
	err := controller.Session.PutFriend(uuid, friendId)
	if err != nil {
		checkForResourceNotFound(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (controller *Controller) DeleteFriend(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	uuid := model.UUID(ps.ByName("uuid"))
	friendId := model.UUID(ps.ByName("friendid"))
	err := controller.Session.DeleteFriend(uuid, friendId)
	if err != nil {
		checkForResourceNotFound(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (controller *Controller) GetFriends(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	uuid := model.UUID(ps.ByName("uuid"))
	friends, err := controller.Session.GetFriends(uuid)
	if err != nil {
		checkForResourceNotFound(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(friends)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (controller *Controller) GetAllUsers(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	users, err := controller.Session.GetAllUsers()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (controller *Controller) SearchTextQuery(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	urlQueryMap := req.URL.Query()
	_, ok := urlQueryMap["query"]
	if len(urlQueryMap) != 1 || !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return		
	}
	queryString := urlQueryMap["query"][0]
	users, err := controller.Session.SearchTextQuery(queryString)
	if err != nil {
		checkForResourceNotFound(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if len(users) > 0 {
		err = json.NewEncoder(w).Encode(users)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}		
	} else {
		w.Write([]byte("[]"))
	}
}

func (controller *Controller) CleanUp() {
	controller.Session.CleanUp()
}

func checkForResourceNotFound(w http.ResponseWriter, err error) {
	if err == mgo.ErrNotFound {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	} else {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
