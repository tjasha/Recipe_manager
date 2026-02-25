package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/tjasha/Recipe_manager/internal/config"
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

	config := config.LoadConfig()

	pool, err := database.NewPool(config.DB_URL)
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

	port := config.PORT
	fmt.Println("port:", port)
	if port == "" {
		log.Fatal("PORT not set")
	}

	// Start the server
	fmt.Println("Server started at http://localhost:" + port)
	err = http.ListenAndServe(fmt.Sprintf(":"+port), nil)
	if err != nil {
		log.Fatal(err)
	}

}
