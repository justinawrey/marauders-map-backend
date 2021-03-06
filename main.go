package main

import (
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/cpen321/groupii-back/controller"
)

func main() {
	router := httprouter.New()
	controller := controller.NewController()
	defer controller.CleanUp()

	router.PUT("/user/:uuid", controller.PutUser)
	router.DELETE("/user/:uuid", controller.DeleteUser)
	router.GET("/user/:uuid", controller.GetUser)
	router.GET("/user", controller.GetAllUsers)

	router.PUT("/location/:uuid", controller.PutUserLoc)
	router.GET("/location/:uuid", controller.GetUserLoc)

	router.PUT("/friend/:uuid/:friendid", controller.PutFriend)
	router.DELETE("/friend/:uuid/:friendid", controller.DeleteFriend)
	router.GET("/friend/:uuid", controller.GetFriends)

	router.GET("/search", controller.SearchTextQuery)
	router.GET("/density", controller.GetDensityMetrics)
	router.GET("/heatmap.png", controller.GetHeatmapPNG)
	router.GET("/heatmap.kml", controller.GetHeatmapKML)

	// check if we are running through heroku or on localhost
	port := ":8088"
	if mongoPort := os.Getenv("PORT"); mongoPort != "" {
		port = ":" + mongoPort
	}

	log.Fatalln(http.ListenAndServe(port, router))
}
