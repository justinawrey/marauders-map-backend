// PRECONDITIONS:
// db on which tests are performed does not contain user with id: thisuserdoesnotexist
// db on which tests are performed contains user with id: thisuserdoesexist

package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/cpen321/groupii-back/model"
)

var testController *Controller

func TestMain(m *testing.M) {
	testController = NewController() // setup before running tests
	exitVal := m.Run()
	testController.CleanUp() // teardown after running tests
	os.Exit(exitVal)
}

func TestPutUserLoc(t *testing.T) {
	putUserLocHandle := httprouter.Handle(testController.PutUserLoc)
	getUserLocHandle := httprouter.Handle(testController.GetUserLoc)
	badLocationQuery := `{"longitude": "edf", "latitude": 1234}`
	correctLocationQuery := `{"longitude": 111, "latitude": 222}`

	// put location to user that does not exist
	resp := makeRequest("PUT",
		"/location/:uuid",
		"/location/thisuserdoesnotexist",
		bytes.NewBuffer([]byte(correctLocationQuery)),
		putUserLocHandle)
	if resp.Code != http.StatusNotFound {
		t.Error("expected 404 response code, got " +
			strconv.Itoa(resp.Code))
	}

	// put badly formatted location
	resp = makeRequest("PUT",
		"/location/:uuid",
		"/location/thisuserdoesexist",
		bytes.NewBuffer([]byte(badLocationQuery)),
		putUserLocHandle)
	if resp.Code != http.StatusBadRequest {
		t.Error("expected 400 response code, got " +
			strconv.Itoa(resp.Code))
	}

	// successful put
	resp = makeRequest("PUT",
		"/location/:uuid",
		"/location/thisuserdoesexist",
		bytes.NewBuffer([]byte(correctLocationQuery)),
		putUserLocHandle)
	if resp.Code != http.StatusNoContent {
		t.Error("expected 204 response code, got " +
			strconv.Itoa(resp.Code))
	}

	// verify that long and lat for user thisuserdoesexist actually changed
	resp = makeRequest("GET",
		"/location/:uuid",
		"/location/thisuserdoesexist",
		bytes.NewBuffer([]byte("")),
		getUserLocHandle)

	var location model.Location
	json.Unmarshal(resp.Body.Bytes(), &location)
	correctLocationResponse := model.Location{Longitude: 111, Latitude: 222}

	if resp.Code != http.StatusOK {
		t.Error("expected 200 response code, got " +
			strconv.Itoa(resp.Code))
	}
	if location != correctLocationResponse {
		t.Error("expected & retrieved locations do not match")
	}
}

func TestGetUserLoc(t *testing.T) {
	getUserLocHandle := httprouter.Handle(testController.GetUserLoc)

	// get location from user that does not exist
	resp := makeRequest("GET",
		"/location/:uuid",
		"/location/thisuserdoesnotexist",
		bytes.NewBuffer([]byte("")),
		getUserLocHandle)

	if resp.Code != http.StatusNotFound {
		t.Error("expected 404 response code, got " +
			strconv.Itoa(resp.Code))
	}

	// successful get location
	resp = makeRequest("GET",
		"/location/:uuid",
		"/location/thisuserdoesexist",
		bytes.NewBuffer([]byte("")),
		getUserLocHandle)

	var location model.Location
	json.Unmarshal(resp.Body.Bytes(), &location)
	correctLocationResponse := model.Location{Longitude: 111, Latitude: 222}

	if resp.Code != http.StatusOK {
		t.Error("expected 200 response code, got " +
			strconv.Itoa(resp.Code))
	}

	if location != correctLocationResponse {
		t.Error("expected and retrieved locations do not match")
	}
}

func TestPutUser(t *testing.T) {
	putUserHandle := httprouter.Handle(testController.PutUser)

 // put badly formatted json

 // successful put user
}

//func TestGetUser(t *testing.T) {
//	getUserHandle := httprouter.Handle(testController.GetUser)

// do tests..
//}

//func TestDeleteUser(t *testing.T) {
//	deleteUserHandle := httprouter.Handle(testController.DeleteUser)

// do tests..
//}

//func TestPutFriend(t *testing.T) {
//	putFriendHandle := httprouter.Handle(testController.PutFriend)

// do tests..
//}

//func TestDeleteFriend(t *testing.T) {
//	deleteFriendHandle := httprouter.Handle(testController.DeleteFriend)

// do tests..
//}

//func TestGetFriends(t *testing.T) {
//	getFriendsHandle := httprouter.Handle(testController.GetFriends)

// do tests..
//}

//func TestGetAllUsers(t *testing.T) {
//	getAllUsersHandle := httprouter.Handle(testController.GetAllUsers)

// do tests..
//}

func makeRequest(method string, uri string, req_uri string, body *bytes.Buffer, handle httprouter.Handle) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(method, req_uri, body)
	testRouter := httprouter.New()
	testRouter.Handle(method, uri, handle)
	testRouter.ServeHTTP(recorder, req)
	return recorder
}
