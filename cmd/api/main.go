package main

import (
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/tjasha/Recipe_manager/internal/config"
	"github.com/tjasha/Recipe_manager/internal/handler"
	"github.com/tjasha/Recipe_manager/internal/repository"
	"github.com/tjasha/Recipe_manager/internal/router"

	//"github.com/golang-migrate/migrate/v4"
	//"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()

	pool, err := repository.NewPool(cfg.DB_URL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// Initialize session manager
	session := scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	// for cross-origin requests (different ports for FE and BE)
	// needs to be used SameSite=None in Secure=true.
	session.Cookie.SameSite = http.SameSiteNoneMode // Allows sending cookies between different domains/ports
	session.Cookie.Secure = true
	//session.Cookie.Secure = cfg.InProduction

	// Create the repository implementation
	dbRepo := repository.NewPostgresRepository(pool)

	app := &handler.Application{
		Config:  cfg,
		Session: session,
		DB:      dbRepo, // concrete implementation of DB
	}

	r := router.New(app)

	log.Println("Server running on :8080")
	// To-Do: in production ListenAndServeTLS for HTTPS needs to be used.
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		return
	}
}
