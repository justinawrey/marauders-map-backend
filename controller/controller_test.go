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
	badUserQuery := `{
		"uuid": "thisuserdoesexist",
		"name": 1234,
		"email": 23,
		"photoURL": "sdf/sdf/sdf/",
		"friends": ["12", "13", "14"],
		"location": {
		  "longitude": 129,
		  "latitude": 456
		}
	  }`

	goodUserQuery := `{
		"uuid": "thisuserdoesexist",
		"name": "jon doe",
		"email": "jon@jon.com",
		"photoURL": "sdf/sdf/sdf/",
		"friends": ["12", "13", "14"],
		"location": {
		  "longitude": 129,
		  "latitude": 456
		}
	  }`

	// put badly formatted user
	resp := makeRequest("PUT",
		"/user/:uuid",
		"/user/thisuserdoesexist",
		bytes.NewBuffer([]byte(badUserQuery)),
		putUserHandle)

	if resp.Code != http.StatusBadRequest {
		t.Error("expected 400 response code, got " +
			strconv.Itoa(resp.Code))
	}

	// successful put user
	resp = makeRequest("PUT",
		"/user/:uuid",
		"/user/thisuserdoesexist",
		bytes.NewBuffer([]byte(goodUserQuery)),
		putUserHandle)

	if resp.Code != http.StatusNoContent {
		t.Error("expected 204 response code, got " +
			strconv.Itoa(resp.Code))
	}
}

func TestGetUser(t *testing.T) {
	getUserHandle := httprouter.Handle(testController.GetUser)

	// get user that does not exist
	resp := makeRequest("GET",
		"/user/:uuid",
		"/user/thisuserdoesnotexist",
		bytes.NewBuffer([]byte("")),
		getUserHandle)

	if resp.Code != http.StatusNotFound {
		t.Error("expected 404 response code, got " +
			strconv.Itoa(resp.Code))
	}

	// successful get user
	resp = makeRequest("GET",
		"/user/:uuid",
		"/user/thisuserdoesexist",
		bytes.NewBuffer([]byte("")),
		getUserHandle)

	if resp.Code != http.StatusOK {
		t.Error("expected 200 response code, got " +
			strconv.Itoa(resp.Code))
	}
}

func TestDeleteUser(t *testing.T) {
	deleteUserHandle := httprouter.Handle(testController.DeleteUser)
	putUserHandle := httprouter.Handle(testController.PutUser)
	goodUserQuery := `{
		"uuid": "thisuserdoesexist",
		"name": "jon doe",
		"email": "jon@jon.com",
		"photoURL": "sdf/sdf/sdf/",
		"friends": ["12", "13", "14"],
		"location": {
		  "longitude": 129,
		  "latitude": 456
		}
	  }`

	// delete user that does not exist
	resp := makeRequest("DELETE",
		"/user/:uuid",
		"/user/thisuserdoesnotexist",
		bytes.NewBuffer([]byte("")),
		deleteUserHandle)

	if resp.Code != http.StatusNotFound {
		t.Error("expected 404 response code, got " +
			strconv.Itoa(resp.Code))
	}

	// successful delete user
	resp = makeRequest("DELETE",
		"/user/:uuid",
		"/user/thisuserdoesexist",
		bytes.NewBuffer([]byte("")),
		deleteUserHandle)

	if resp.Code != http.StatusNoContent {
		t.Error("expected 204 response code, got " +
			strconv.Itoa(resp.Code))
	}

	// put the test user back
	resp = makeRequest("PUT",
		"/user/:uuid",
		"/user/thisuserdoesexist",
		bytes.NewBuffer([]byte(goodUserQuery)),
		putUserHandle)

	if resp.Code != http.StatusNoContent {
		t.Error("expected 204 response code, got " +
			strconv.Itoa(resp.Code))
	}
}

func TestPutFriend(t *testing.T) {
	putFriendHandle := httprouter.Handle(testController.PutFriend)

	// add friend to a user that does not exist
	resp := makeRequest("PUT",
		"/friend/:uuid/:friendid",
		"/friend/thisuserdoesnotexist/thisuserdoesexist",
		bytes.NewBuffer([]byte("")),
		putFriendHandle)

	if resp.Code != http.StatusNotFound {
		t.Error("Expected 404 response code, got " +
			strconv.Itoa(resp.Code))
	}

	// add friend to a user that does exist
	resp = makeRequest("PUT",
		"/friend/:uuid/:friendid",
		"/friend/thisuserdoesexist/12345",
		bytes.NewBuffer([]byte("")),
		putFriendHandle)

	if resp.Code != http.StatusNoContent {
		t.Error("Expected 204 response code, got " +
			strconv.Itoa(resp.Code))
	}
}

func TestDeleteFriend(t *testing.T) {
	deleteFriendHandle := httprouter.Handle(testController.DeleteFriend)

	// delete friend from a user that does not exist
	resp := makeRequest("DELETE",
		"/friend/:uuid/:friendid",
		"/friend/thisuserdoesnotexist/thisuserdoesexist",
		bytes.NewBuffer([]byte("")),
		deleteFriendHandle)

	if resp.Code != http.StatusNotFound {
		t.Error("Expected 404 response code, got " +
			strconv.Itoa(resp.Code))
	}

	// delete friend from a user that does exist
	resp = makeRequest("DELETE",
		"/friend/:uuid/:friendid",
		"/friend/thisuserdoesexist/12",
		bytes.NewBuffer([]byte("")),
		deleteFriendHandle)

	if resp.Code != http.StatusNoContent {
		t.Error("Expected 204 response code, got " +
			strconv.Itoa(resp.Code))
	}
}

func TestGetFriends(t *testing.T) {
	getFriendsHandle := httprouter.Handle(testController.GetFriends)

	// get friends from a user that does not exist
	resp := makeRequest("GET",
		"/friend/:uuid",
		"/friend/thisuserdoesnotexist",
		bytes.NewBuffer([]byte("")),
		getFriendsHandle)

	if resp.Code != http.StatusNotFound {
		t.Error("Expected 404 response code, got " +
			strconv.Itoa(resp.Code))
	}

	// get friends from a user that does exist
	resp = makeRequest("GET",
		"/friend/:uuid",
		"/friend/thisuserdoesexist",
		bytes.NewBuffer([]byte("")),
		getFriendsHandle)

	if resp.Code != http.StatusOK {
		t.Error("Expected 200 response code, got " +
			strconv.Itoa(resp.Code))
	}
}

func TestGetAllUsers(t *testing.T) {
	getAllUsersHandle := httprouter.Handle(testController.GetAllUsers)
	
	// get all users
	resp := makeRequest("GET",
		"/user",
		"/user",
		bytes.NewBuffer([]byte("")),
		getAllUsersHandle)

	if resp.Code != http.StatusOK {
		t.Error("Expected 200 response code, got " +
			strconv.Itoa(resp.Code))
	}
}

func makeRequest(method string, uri string, req_uri string, body *bytes.Buffer, handle httprouter.Handle) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(method, req_uri, body)
	testRouter := httprouter.New()
	testRouter.Handle(method, uri, handle)
	testRouter.ServeHTTP(recorder, req)
	return recorder
}
