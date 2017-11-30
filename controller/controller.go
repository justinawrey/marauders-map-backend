package controller

import (
	"encoding/json"
	"image"
	"image/png"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"time"

	mgo "gopkg.in/mgo.v2"

	"github.com/dustin/go-heatmap"
	"github.com/dustin/go-heatmap/schemes"
	"github.com/julienschmidt/httprouter"
	"gitlab.com/cpen321/groupii-back/model"
	"go.skia.org/infra/perf/go/kmeans"
)

type Controller struct {
	Session *model.MgoSession
}

func NewController() *Controller {
	return &Controller{
		Session: model.NewSession(),
	}
}

func randomInRange(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min) + min
}

type observation struct {
	longitude float64
	latitude  float64
}

func (obs observation) Distance(c kmeans.Clusterable) float64 {
	other := c.(observation)
	return math.Sqrt((obs.longitude-other.longitude)*(obs.longitude-other.longitude) +
		(obs.latitude-other.latitude)*(obs.latitude-other.latitude))
}

func (obs observation) AsClusterable() kmeans.Clusterable {
	return obs
}

// calculateCentroid implements CalculateCentroid.
func calculateCentroid(members []kmeans.Clusterable) kmeans.Centroid {
	var sumX = 0.0
	var sumY = 0.0
	length := float64(len(members))

	for _, m := range members {
		sumX += m.(observation).longitude
		sumY += m.(observation).latitude
	}
	return observation{longitude: sumX / length, latitude: sumY / length}
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
		checkForResourceNotFound(w, err)
		return
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

func (controller *Controller) GetDensityMetrics(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	users, err := controller.Session.GetAllUsers()
	if err != nil {
		checkForResourceNotFound(w, err)
		return
	}

	// set up observations for kmeans
	var observations []kmeans.Clusterable
	for _, user := range users {
		observations = append(observations, observation{
			longitude: user.Location.Longitude,
			latitude:  user.Location.Latitude,
		})
	}

	// set up centroids for kmeans
	var centroids []kmeans.Centroid
	for i := 0; i < 100000; i++ {
		long := float64(randomInRange(-180, 180))
		lat := float64(randomInRange(-90, 90))
		centroids = append(centroids, observation{
			longitude: long,
			latitude:  lat,
		})
	}

	// get updated clusters
	for i := 0; i < 100000; i++ {
		centroids = kmeans.Do(observations, centroids, calculateCentroid)
	}
	clusters, _ := kmeans.GetClusters(observations, centroids)

	var retData []struct {
		Longitude float64
		Latitude  float64
		Radius    float64
	}
	for _, cluster := range clusters {
		retData = append(retData, struct {
			Longitude float64
			Latitude  float64
			Radius    float64
		}{
			Longitude: cluster[0].(observation).longitude,
			Latitude:  cluster[0].(observation).latitude,
			Radius:    float64(len(cluster) - 1), // first element of cluster is centroid, this is not actual data
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(retData)
}

func (controller *Controller) GetHeatmapKML(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	users, err := controller.Session.GetAllUsers()
	if err != nil {
		checkForResourceNotFound(w, err)
		return
	}
	points := []heatmap.DataPoint{}
	for _, user := range users {
		points = append(points, heatmap.P(user.Location.Longitude,
			user.Location.Latitude))
	}
	w.Header().Set("Content-Type", "application/vnd.google-earth.kml+xml")
	heatmap.KML(image.Rect(0, 0, 1024, 1024), points, 200, 128, schemes.AlphaFire, "https://maraudersss.herokuapp.com/heatmap.png", w)
}

func (controller *Controller) GetHeatmapPNG(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	users, err := controller.Session.GetAllUsers()
	if err != nil {
		checkForResourceNotFound(w, err)
		return
	}
	points := []heatmap.DataPoint{}
	for _, user := range users {
		points = append(points, heatmap.P(user.Location.Longitude,
			user.Location.Latitude))
	}
	mapimg := heatmap.Heatmap(image.Rect(0, 0, 1024, 1024), points, 200, 128, schemes.AlphaFire)
	w.Header().Set("Content-Type", "image/png")
	png.Encode(w, mapimg)
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
