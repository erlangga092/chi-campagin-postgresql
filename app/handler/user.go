package handler

import (
	"encoding/json"
	"funding-app/app/auth"
	"funding-app/app/helper"
	"funding-app/app/user"
	"net/http"
	"runtime"
	"strings"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

// alias map
type M map[string]interface{}

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
		errorMessage := "Content must be application/json"

		response := helper.APIResponse("Failed to register user", http.StatusBadRequest, "error", errorMessage)
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	log.Info("goroutine-start-register-handler : ", runtime.NumGoroutine())

	v := validator.New()
	input := user.RegisterUserInput{}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response := helper.APIResponse("Failed to register user", http.StatusBadRequest, "error", err.Error())
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	// validate input
	err = v.Struct(input)
	if err != nil {
		var errors []string

		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, e.Error())
		}

		response := helper.APIResponse("Failed to register user", http.StatusUnprocessableEntity, "error", errors)
		helper.JSON(w, response, http.StatusUnprocessableEntity)
		return
	}

	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		response := helper.APIResponse("Failed to register user", http.StatusBadRequest, "error", err.Error())
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	token, err := h.authService.GenerateToken(newUser.ID)
	if err != nil {
		response := helper.APIResponse("Failed to register user", http.StatusBadRequest, "error", err.Error())
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	// check goroutine
	log.Info("goroutine-end-register-handler : ", runtime.NumGoroutine())

	formatter := user.FormatUser(newUser, token)
	response := helper.APIResponse("Account has been created", http.StatusCreated, "success", formatter)
	helper.JSON(w, response, http.StatusCreated)
}

func (h *userHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		errorMessage := "Content must be application/json"

		response := helper.APIResponse("Login user failed", http.StatusBadRequest, "error", errorMessage)
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	log.Info("goroutine-start-login-handler : ", runtime.NumGoroutine())

	v := validator.New()
	input := user.LoginUserInput{}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response := helper.APIResponse("Login user failed", http.StatusBadRequest, "error", err.Error())
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	// validate input
	err = v.Struct(input)
	if err != nil {
		var errors []string

		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, e.Error())
		}

		response := helper.APIResponse("Login user failed", http.StatusUnprocessableEntity, "error", errors)
		helper.JSON(w, response, http.StatusUnprocessableEntity)
		return
	}

	loggedInUser, err := h.userService.LoginUser(input)
	if err != nil {
		response := helper.APIResponse("Login user failed", http.StatusBadRequest, "error", err.Error())
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	token, err := h.authService.GenerateToken(loggedInUser.ID)
	if err != nil {
		response := helper.APIResponse("Login user failed", http.StatusBadRequest, "error", err.Error())
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	// check goroutine
	log.Info("goroutine-end-login-handler : ", runtime.NumGoroutine())

	formatter := user.FormatUser(loggedInUser, token)
	response := helper.APIResponse("Login successfully", http.StatusOK, "success", formatter)
	helper.JSON(w, response, http.StatusOK)
}

func (h *userHandler) IsEmailAvailable(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		errorMessage := "Content type must be application/json"

		response := helper.APIResponse("Checking email failed", http.StatusBadRequest, "error", errorMessage)
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	log.Info("goroutine-start-login-handler : ", runtime.NumGoroutine())

	v := validator.New()
	input := user.CheckEmailInput{}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response := helper.APIResponse("Checking email failed", http.StatusBadRequest, "error", err.Error())
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	// validate input
	err = v.Struct(input)
	if err != nil {
		var errors []string

		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, e.Error())
		}

		response := helper.APIResponse("Checking email failed", http.StatusUnprocessableEntity, "error", errors)
		helper.JSON(w, response, http.StatusUnprocessableEntity)
		return
	}

	isAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		response := helper.APIResponse("Checking email failed", http.StatusUnprocessableEntity, "error", err.Error())
		helper.JSON(w, response, http.StatusUnprocessableEntity)
		return
	}

	data := M{
		"is_available": isAvailable,
	}

	// check goroutine
	log.Info("goroutine-end-login-handler : ", runtime.NumGoroutine())

	response := helper.APIResponse("Success checking email", http.StatusOK, "success", data)
	helper.JSON(w, response, http.StatusOK)
}

func (h *userHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		errorMessage := "Content-Type must be multipart/form-data"

		response := helper.APIResponse("Failed to upload avatar", http.StatusBadRequest, "error", errorMessage)
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	log.Info("OK!")
}
