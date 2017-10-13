package controller

import (
	"net/http"

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
}

func (controller *Controller) GetUserLoc(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
}

func (controller *Controller) PutUser(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
}

func (controller *Controller) GetUser(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
}

func (controller *Controller) PutFriend(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
}

func (controller *Controller) DeleteFriend(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
}

func (controller *Controller) GetFriends(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
}

func (controller *Controller) CleanUp() {
}
