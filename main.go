package main

import (
	"fmt"
	"funding-app/app/auth"
	"funding-app/app/handler"
	"funding-app/app/user"
	"funding-app/database"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error when load .env file", err)
	}
}

func main() {
	db, err := database.GetConnection()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	status := "up"
	eb := db.Ping()
	if eb != nil {
		status = "down"
		panic(err.Error())
	}

	log.Println(status)
	fmt.Println("postgreSQL connected!")

	// repository
	userRepository := user.NewUserRepository(db)

	// service
	userService := user.NewService(userRepository)
	authService := auth.NewJwtService()

	// handler
	userHandler := handler.NewUserHandler(userService, authService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/api/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Welcome!"))
			})
		})

		r.Group(func(r chi.Router) {
			r.Post("/users", userHandler.RegisterUser)
			r.Post("/sessions", userHandler.LoginUser)
		})
	})

	fmt.Println("Server running on port - 9000")
	err = http.ListenAndServe(":9000", r)
	if err != nil {
		log.Fatal(err)
	}
}
