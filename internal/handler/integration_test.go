//go:build integration

package handler_test

import (
	"net/http"
	"testing"
)

// ---- Integration tests for handler.go ----

func TestIntegration_Logout(t *testing.T) {

	// start test server (from helpers_test.go)
	ts := NewTestServer()
	defer ts.Close()

	// we should first call login endpoint to create session,
	//for this test we'll just set sesssion manually on the server

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
