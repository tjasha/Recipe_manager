package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/tjasha/Recipe_manager/internal/database"
	apphttp "github.com/tjasha/Recipe_manager/internal/handlers"
	"github.com/tjasha/Recipe_manager/internal/repository"
	"github.com/tjasha/Recipe_manager/internal/service"

	//"github.com/golang-migrate/migrate/v4"
	//"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {

	dbURL := os.Getenv("DB_URL")
	fmt.Println(dbURL)
	if dbURL == "" {
		log.Fatal("DB_URL not set")
	}

	pool, err := database.NewPool(dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	repo := repository.NewPostgresRepository(pool)
	userService := service.NewService(repo)
	handler := apphttp.NewHandler(userService)

	http.HandleFunc("/", handler.CreateUserHandler)

	//// Recipes API - reading from json file
	//http.HandleFunc("/", handlers.Handler)

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
