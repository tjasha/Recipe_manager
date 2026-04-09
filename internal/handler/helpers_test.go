package handler_test

import (
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"

	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	"github.com/tjasha/Recipe_manager/internal/config"
	"github.com/tjasha/Recipe_manager/internal/handler"
	"github.com/tjasha/Recipe_manager/internal/repository"
	"github.com/tjasha/Recipe_manager/internal/router"
)

// ---- this is a helper file that creates all instances that are
//      repeatedly needed in unit and integration tests -------

// newTestApplication creates and return new testing instance.
func NewTestApplication() *handler.Application {
	store := memstore.New()
	session := scs.New()
	session.Store = store

	testConfig := &config.Config{
		GoogleOauthClientID: "test-client-id",
	}

	return &handler.Application{
		Config:  testConfig,
		Session: session,
		DB:      &repository.MockRepository{},
	}
}

// testServer struct including test server and client.
type TestServer struct {
	*httptest.Server
	Client      *http.Client
	Application *handler.Application
}

// newTestServer starts new test server with our application and returns it.
func NewTestServer() *TestServer {
	app := NewTestApplication()
	mux := router.New(app)
	ts := httptest.NewServer(mux)

	// create client with cookie jar-om
	jar, _ := cookiejar.New(nil)
	client := ts.Client()
	client.Jar = jar

	return &TestServer{
		Server:      ts,
		Client:      client,
		Application: app,
	}
}
