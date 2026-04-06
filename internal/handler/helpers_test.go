package handler_test

import (
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"

	"github.com/alexedwards/scs/v2"
	"github.com/tjasha/Recipe_manager/internal/handler"
	"github.com/tjasha/Recipe_manager/internal/repository"
	"github.com/tjasha/Recipe_manager/internal/router"
)

// ---- this is a helper file that creates all instances that are
//      repeatedly needed in unit and integration tests -------

// newTestApplication creates and return new testing instance.
func newTestApplication() *handler.Application {
	session := scs.New()
	return &handler.Application{
		Session: session,
		DB:      &repository.MockRepository{},
	}
}

// testServer struct including test server and client.
type testServer struct {
	*httptest.Server
	Client *http.Client
}

// newTestServer starts new test server with our application and returns it.
func NewTestServer() *testServer {
	app := newTestApplication()
	mux := router.New(app)
	ts := httptest.NewServer(mux)

	// create client with cookie jar-om
	jar, _ := cookiejar.New(nil)
	client := ts.Client()
	client.Jar = jar

	return &testServer{
		Server: ts,
		Client: client,
	}
}
