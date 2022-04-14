package handler

import (
	"encoding/json"
	"funding-app/app/auth"
	"funding-app/app/helper"
	"funding-app/app/user"
	"net/http"
	"runtime"

	log "github.com/sirupsen/logrus"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{
		userService: userService,
		authService: authService,
	}
}

func (h *userHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content must be application/json", http.StatusBadRequest)
		return
	}

	log.Info("goroutine-start-register-handler : ", runtime.NumGoroutine())

	input := user.RegisterUserInput{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		helper.JSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		helper.JSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.authService.GenerateToken(newUser.ID)
	if err != nil {
		helper.JSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check goroutine
	log.Info("goroutine-end-register-handler : ", runtime.NumGoroutine())
	response := user.FormatUser(newUser, token)
	helper.JSON(w, response, http.StatusCreated)
}

func (h *userHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content must be application/json", http.StatusBadRequest)
		return
	}

	log.Info("goroutine-start-login-handler : ", runtime.NumGoroutine())

	input := user.LoginUserInput{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		helper.JSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	loggedInUser, err := h.userService.LoginUser(input)
	if err != nil {
		helper.JSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.authService.GenerateToken(loggedInUser.ID)
	if err != nil {
		helper.JSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check goroutine
	log.Info("goroutine-end-login-handler : ", runtime.NumGoroutine())
	response := user.FormatUser(loggedInUser, token)
	helper.JSON(w, response, http.StatusOK)
}
