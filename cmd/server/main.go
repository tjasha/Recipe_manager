package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tjasha/Recipe_manager/internal/config"
	"github.com/tjasha/Recipe_manager/internal/database"
	"github.com/tjasha/Recipe_manager/internal/handlers"

	//"github.com/golang-migrate/migrate/v4"
	//"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var app config.AppConfig // now we can also use it in routes

func main() {

	pool, err := run()
	if err != nil {
		log.Fatal(err)
	}
	//we're closing connection here not in the run()
	defer pool.Close()

	// Recipes API
	http.HandleFunc("/", handlers.RecipesHandler)

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func run() (*pgxpool.Pool, *config.AppConfig) {
	//
	// defining is application is in production. default is true
	inProduction := flag.Bool("production", true, "Application is in production")
	// define cash, default is true
	useCache := flag.Bool("cache", true, "Use template cache")
	// DB info
	dbHost := flag.String("dbhost", "localhost", "Database host")
	dbName := flag.String("dbname", "", "Database name")
	dbUser := flag.String("dbuser", "", "Database user")
	dbPass := flag.String("dbpass", "", "Database password")
	dbPort := flag.String("dbport", "5445", "Database port")
	dbSSL := flag.String("dbssl", "disable", "Database ssl settings (disable, prefer, require)")

	// this is actually parsing the flags and make it usable
	flag.Parse()
	if *dbName == "" || *dbUser == "" || *dbHost == "" || *dbPort == "" {
		fmt.Println("Missing required flags")
		os.Exit(1)
	}
	//
	// this is read from the flag
	app.InProduction = *inProduction
	app.UseCache = *useCache

	dbURL := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL)

	pool, err := database.NewPool(dbURL)
	if err != nil {
		log.Fatal("Cannot create DB pool:", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}

	log.Println("Connected to PostgreSQL")

	return pool, nil
}
