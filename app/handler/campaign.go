package handler

import (
	"encoding/json"
	"funding-app/app/campaign"
	"funding-app/app/helper"
	"funding-app/app/key"
	"funding-app/app/user"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type campaignHandler struct {
	campaignService campaign.Service
}

func NewCampaignHandler(campaignService campaign.Service) *campaignHandler {
	return &campaignHandler{campaignService}
}

func (h *campaignHandler) GetCampaigns(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")

	campaigns, err := h.campaignService.GetCampaigns(userID)
	if err != nil {
		response := helper.APIResponse("Failed to get campaigns", http.StatusBadRequest, "error", err.Error())
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	formatter := campaign.FormatCampaigns(campaigns)
	response := helper.APIResponse("List of campaigns", http.StatusOK, "success", formatter)
	helper.JSON(w, response, http.StatusOK)
}

func (h *campaignHandler) GetCampaignDetail(w http.ResponseWriter, r *http.Request) {
	campaignID := chi.URLParam(r, "id")

	detailCampaign, err := h.campaignService.GetCampaignDetail(campaignID)
	if err != nil {
		response := helper.APIResponse("Failed to get campaign", http.StatusBadRequest, "error", err.Error())
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	formatter := campaign.FormatCampaign(detailCampaign)
	response := helper.APIResponse("Detail of campaigns", http.StatusOK, "success", formatter)
	helper.JSON(w, response, http.StatusOK)
}

func (h *campaignHandler) CreateCampaign(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		errorMessage := "Content must be application/json"

		response := helper.APIResponse("Failed to create campaign", http.StatusBadRequest, "error", errorMessage)
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	v := validator.New()
	input := campaign.CreateCampaignInput{}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response := helper.APIResponse("Failed to create campaign", http.StatusBadRequest, "error", err.Error())
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

		response := helper.APIResponse("Failed to create campaign", http.StatusUnprocessableEntity, "error", errors)
		helper.JSON(w, response, http.StatusUnprocessableEntity)
		return
	}

	// get data user from middleware
	user := r.Context().Value(key.CtxAuthKey{}).(user.User)
	input.User = user

	newCampaign, err := h.campaignService.CreateCampaign(input)
	if err != nil {
		response := helper.APIResponse("Failed to create campaign", http.StatusBadRequest, "error", err.Error())
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	formatter := campaign.FormatCampaign(newCampaign)
	response := helper.APIResponse("Detail of campaigns", http.StatusOK, "success", formatter)
	helper.JSON(w, response, http.StatusOK)
}

func (h *campaignHandler) UploadCampaignImage(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		errorMessage := "Content must be multipart/form-data"

		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", errorMessage)
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	err := r.ParseMultipartForm(1024)
	if err != nil {
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", err.Error())
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	uploadedFile, _, err := r.FormFile("image")
	if err != nil {
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", err.Error())
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	defer uploadedFile.Close()

	// get user data from middleware
	user := r.Context().Value(key.CtxAuthKey{}).(user.User)

	isPrimary := false
	isPrimaryInput := r.FormValue("is_primary")
	isPrimaryInputBool, _ := strconv.ParseBool(isPrimaryInput)

	if isPrimaryInput != "" && isPrimaryInputBool {
		isPrimary = true
	}

	input := campaign.CreateCampaignImageInput{}
	input.CampaignID = r.FormValue("campaign_id")
	input.User = user
	input.IsPrimary = isPrimary

	_, err = h.campaignService.UploadCampaignImage(input, uploadedFile)
	if err != nil {
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", err.Error())
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	data := M{
		"is_uploaded": true,
	}

	response := helper.APIResponse("Success upload campaign image", http.StatusCreated, "success", data)
	helper.JSON(w, response, http.StatusCreated)
}
