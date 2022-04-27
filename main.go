package main

import (
	"fmt"
	"funding-app/app/auth"
	"funding-app/app/campaign"
	"funding-app/app/handler"
	cm "funding-app/app/middleware"
	"funding-app/app/user"
	"funding-app/database"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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
	campaignRepository := campaign.NewCampaignRepository(db)

	// service
	userService := user.NewService(userRepository)
	authService := auth.NewJwtService()
	campaignService := campaign.NewCampaignService(campaignRepository)

	// handler
	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)

	// initial route
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// list of route
	r.Route("/api/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Welcome!"))
			})
		})

		r.Group(func(r chi.Router) {
			r.Post("/users", userHandler.RegisterUser)
			r.Post("/sessions", userHandler.LoginUser)
			r.Post("/email_checkers", userHandler.IsEmailAvailable)

			r.With(func(h http.Handler) http.Handler {
				return cm.AuthMiddleware(h, authService, userService)
			}).Post("/avatars", userHandler.UploadAvatar)

			r.With(func(h http.Handler) http.Handler {
				return cm.AuthMiddleware(h, authService, userService)
			}).Post("/refresh-token", userHandler.RefreshToken)
		})

		r.Group(func(r chi.Router) {
			r.Get("/campaigns", campaignHandler.GetCampaigns)
			r.Get("/campaigns/{id}", campaignHandler.GetCampaignDetail)

			r.With(func(h http.Handler) http.Handler {
				return cm.AuthMiddleware(h, authService, userService)
			}).Post("/campaigns", campaignHandler.CreateCampaign)

			r.With(func(h http.Handler) http.Handler {
				return cm.AuthMiddleware(h, authService, userService)
			}).Post("/campaign-images", campaignHandler.UploadCampaignImage)
		})
	})

	fmt.Println("Server running on port - 9000")
	err = http.ListenAndServe(":9000", r)
	if err != nil {
		log.Fatal(err)
	}
}
