# Recipe-manager
This is a repository that contains code for recipe manager application 

## Goal

The goal of this project is to:
- Focus primarily on learning and developing the Go backend
- Keep the frontend simple

Technologies and tools used
- Golang standard library for backend (using HTTP server)
- React frontend built with Vite


## How to run the project 

Prerequisites:
- Go installed
- Node.js + npm installed

1. Build the frontend
   - cd frontend
   - npm install
   - npm run build (This generates the production-ready frontend in frontend/dist.)
   - in developement run npm run dev
   - frontend will run on http://localhost:5173/

2. Run the backend
   - from the root directory
   - go run cmd/server/main.go
   - backend is running on http://localhost:8080

## Database

Project is running Postgres database <br>
My local database is connected with 
`export DB_URL=postgres://tjasaspes:@localhost:5445/recipe-manager?sslmode=disable`

For migrations [migrate](https://github.com/golang-migrate/migrate) is used

- To install it run `brew install golang-migrate`
- Install migrate driver
```aiexclude
go get -u github.com/golang-migrate/migrate/v4
go get -u github.com/golang-migrate/migrate/v4/database/postgres
go get -u github.com/golang-migrate/migrate/v4/source/file
```
- create migrations files `migrate create -ext sql -dir migrations -seq <file_name>`         

- run migrations 
```aiexclude
migrate -path migrations -database "$DB_URL" up
migrate -path migrations -database "$DB_URL" down
```

### Driver 

Project is using pgxpool driver.

## Session management

For session management I am using SCS - https://github.com/alexedwards/scs

## Google OAuth

For managing users, Google OAuth is used