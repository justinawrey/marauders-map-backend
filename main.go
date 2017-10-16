package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/cpen321/groupii-back/controller"
)

func main() {
	router := httprouter.New()
	controller := controller.New()
	defer controller.CleanUp()

	router.PUT("/location/:uuid", controller.PutUserLoc)
	router.GET("/location/:uuid", controller.GetUserLoc)

	router.PUT("/user/:uuid", controller.PutUser)
	router.GET("/user/:uuid", controller.GetUser)
	router.DELETE("/user/:uuid", controller.DeleteUser)

	router.PUT("/friend/:uuid/:friendid", controller.PutFriend)
	router.DELETE("/friend/:uuid/:friendid", controller.DeleteFriend)
	router.GET("/friend/:uuid", controller.GetFriends)
	
	http.ListenAndServe(":9898", router)
}
