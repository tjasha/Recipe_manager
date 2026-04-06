package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tjasha/Recipe_manager/internal/handler"
)

// ---- Unit tests for handler.go ----

// testing logout handler
func TestHandler_Logout(t *testing.T) {
	// set up testing application (from helpers_test.go)
	app := newTestApplication()

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
