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
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = cfg.InProduction

	// Create the repository implementation
	dbRepo := repository.NewPostgresRepository(pool)

	app := &handler.Application{
		Config:  cfg,
		Session: session,
		DB:      dbRepo, // concrete implementation of DB
	}

	//repo := repository.NewPostgresRepository(pool)
	//userService := service.NewService(repo)
	//handler := apphttp.NewHandler(userService)

	//http.HandleFunc("/", handler.CreateUserHandler)

	//// Recipes API - reading from json file
	//http.HandleFunc("/", handlers.Handler)

	//// Serve frontend files
	//fs := http.FileServer(http.Dir("/Users/tjasaspes/Documents/Learning Go/Recipe_manager/frontend/src/App.tsx"))
	//http.Handle("/", fs)
	//
	//port := config.PORT
	//fmt.Println("port:", port)
	//if port == "" {
	//	log.Fatal("PORT not set")
	//}
	//
	//// Start the server
	//fmt.Println("Server started at http://localhost:" + port)
	//err = http.ListenAndServe(fmt.Sprintf(":"+port), nil)
	//if err != nil {
	//	log.Fatal(err)
	//}

	r := router.New(app)

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)
}
