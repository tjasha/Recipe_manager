//go:build integration

package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/tjasha/Recipe_manager/internal/model"
	"github.com/tjasha/Recipe_manager/internal/repository"
)

// ---- Integration tests for handler.go ----

// helper function to simulate login
func performLogin(t *testing.T, ts *TestServer, userToLogIn *model.User) *http.Client {

	// access mock repository
	db := ts.Application.DB
	// type assertion to make sure that we're using mock repository
	mockRepo, ok := db.(*repository.MockRepository)
	if !ok {
		t.Fatal("Could not assert DB to *repository.MockRepository")
	}

	// set up test user
	mockRepo.User = userToLogIn
	mockRepo.SkipGoogleTokenValidation = true

	// creating face google token
	tokenPayload := map[string]string{"credential": "mock-google-token"}
	body, _ := json.Marshal(tokenPayload)

	// calling real login endpoint
	res, err := ts.Client.Post(ts.URL+"/api/auth/google/verify", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Login request failed: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Login failed with status: %d", res.StatusCode)
	}

	// return client that now has saved session cookie
	return ts.Client
}

// Test for log out
func TestIntegration_Logout(t *testing.T) {

	// start test server (from helpers_test.go)
	ts := NewTestServer()
	defer ts.Close()

	// we should first call login endpoint to create session,
	//for this test we'll just set session manually on the server

	// send POST request
	res, err := ts.Client.Post(ts.URL+"/api/auth/logout", "application/json", nil)
	if err != nil {
		t.Fatal(err)
	}

	// check results
	if res.StatusCode != http.StatusOK {
		t.Errorf("Logout integration test returned wrong status code: got %v want %v", res.StatusCode, http.StatusOK)
	}
	// Here we should check if cookies in response are deleted or changed
}

// Test for getting all users
func TestIntegration_ReturnAllUsers(t *testing.T) {

	testCases := []struct {
		name           string
		expectedStatus int
		userToLogIn    *model.User
	}{
		{
			name:           "Admin requests all users",
			expectedStatus: http.StatusOK,
			userToLogIn:    &model.User{ID: 1, UserName: "Admin User", Email: "admin@example.com", AccessLevel: 0},
		},
		{
			name:           "Chef requests all users",
			expectedStatus: http.StatusForbidden,
			userToLogIn:    &model.User{ID: 1, UserName: "Chef User", Email: "chef@example.com", AccessLevel: 1},
		},
		{
			name:           "Guest requests all users",
			expectedStatus: http.StatusUnauthorized,
			userToLogIn:    nil, // No user logged in
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ts := NewTestServer()
			defer ts.Close()

			loggedInUser := &http.Client{}
			if tc.userToLogIn != nil {
				loggedInUser = performLogin(t, ts, tc.userToLogIn)
			}

			res, err := loggedInUser.Get(ts.URL + "/api/admin/users")
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()

			// check results
			if res.StatusCode != tc.expectedStatus {
				t.Errorf("expected status %d; got %d", tc.expectedStatus, res.StatusCode)
			}
		})
	}
}

// Test for deleting a user
func TestIntegration_DeleteUser_Authorization(t *testing.T) {

	testCases := []struct {
		name           string
		expectedStatus int
		userToLogIn    *model.User
		targetUserID   string
	}{
		{
			name:           "Admin deletes another user",
			expectedStatus: http.StatusNoContent,
			userToLogIn:    &model.User{ID: 1, UserName: "Admin User", Email: "admin@example.com", AccessLevel: 0},
			targetUserID:   "2",
		},
		{
			name:           "Admin tries to delete self",
			expectedStatus: http.StatusForbidden,
			userToLogIn:    &model.User{ID: 1, UserName: "Admin User", Email: "admin@example.com", AccessLevel: 0},
			targetUserID:   "1",
		},
		{
			name:           "Chef tries to delete a user",
			expectedStatus: http.StatusForbidden,
			userToLogIn:    &model.User{ID: 2, UserName: "Chef User", Email: "chef@example.com", AccessLevel: 1},
			targetUserID:   "1",
		},
		{
			name:           "Unauthorized user tries to delete a user",
			expectedStatus: http.StatusUnauthorized,
			userToLogIn:    nil, // No user logged in
			targetUserID:   "1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ts := NewTestServer()
			defer ts.Close()

			// log in user if required
			loggedInUser := &http.Client{}
			if tc.userToLogIn != nil {
				loggedInUser = performLogin(t, ts, tc.userToLogIn)
			}

			// prepare delete request
			req, _ := http.NewRequest("DELETE", ts.URL+"/api/admin/users/"+tc.targetUserID, nil)
			res, err := loggedInUser.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()

			if res.StatusCode != tc.expectedStatus {
				t.Errorf("expected status %d; got %d", tc.expectedStatus, res.StatusCode)
			}
		})
	}
}
