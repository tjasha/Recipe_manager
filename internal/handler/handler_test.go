package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tjasha/Recipe_manager/internal/handler"
	"github.com/tjasha/Recipe_manager/internal/model"
	"github.com/tjasha/Recipe_manager/internal/repository"
)

// ---- Unit tests for handler.go ----

// testing logout handler
func TestHandler_Logout(t *testing.T) {
	// set up testing application (from helpers_test.go)
	app := NewTestApplication()

	// prepare http request, recorder, mock handler
	req, _ := http.NewRequest("POST", "/api/auth/logout", nil)
	ctx, _ := app.Session.Load(req.Context(), "")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	// create handler
	h := handler.New(app)
	// add test userID to session, so that we can check if it's gonna be deleted
	h.App.Session.Put(req.Context(), "userID", uint(1))
	//call handler
	h.Logout(rr, req)

	// checking results
	if rr.Code != http.StatusOK {
		t.Errorf("Logout handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}
	userInSession := h.App.Session.Exists(req.Context(), "userID")
	if userInSession {
		t.Errorf("userID should not exist in session after logout")
	}
}

// testing return all users handler
func TestUnit_ReturnAllUsers(t *testing.T) {

	testCases := []struct {
		name             string
		accessLevel      int
		expectedStatus   int
		expectedLen      int
		expectedUserName string
	}{
		{
			name:             "Admin requests all users",
			accessLevel:      0,
			expectedStatus:   http.StatusOK,
			expectedLen:      1,
			expectedUserName: "Test User",
		}, {
			name:           "Chef requests all users",
			accessLevel:    1,
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tc := range testCases {
		ts := NewTestServer()
		defer ts.Close()

		// set up testing application, data, DB, handler
		app := NewTestApplication()

		expectedUsers := []model.User{{ID: 1, UserName: "Test User"}}
		app.DB.(*repository.MockRepository).Users = expectedUsers

		h := handler.New(app)

		// simulating middleware chain - create middleware that save data in the session
		//it's done after LoadAndSave and before testing handler
		injector := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// here context is already ready by LoadAndSave
			app.Session.Put(r.Context(), "userID", uint(1))
			app.Session.Put(r.Context(), "accessLevel", tc.accessLevel)
			// now really call handler
			h.ReturnAllUsers(w, r)
		})
		// Wrap injector with middleware
		handlerWithSession := app.Session.LoadAndSave(injector)

		// create http request, recorder, mock handler
		req, _ := http.NewRequest("GET", "/api/admin/users", nil)
		rr := httptest.NewRecorder()
		handlerWithSession.ServeHTTP(rr, req)

		// checking results
		if rr.Code != tc.expectedStatus {
			t.Errorf("expected status %d; got %d", tc.expectedStatus, rr.Code)
		}

		if tc.expectedStatus == http.StatusOK {
			// checking correct response JSON
			var actualUsers []model.User
			if err := json.NewDecoder(rr.Body).Decode(&actualUsers); err != nil {
				t.Fatalf("Failed to decode response body: %v", err)
			}
			if len(actualUsers) != tc.expectedLen {
				t.Errorf("unexpected number of users returned: want %v got %v", tc.expectedLen, len(actualUsers))
			}
			if actualUsers[0].UserName != tc.expectedUserName {
				t.Errorf("unexpected user saved in session: want %v got %v", tc.expectedUserName, actualUsers[0].UserName)
			}
		}
	}
}
