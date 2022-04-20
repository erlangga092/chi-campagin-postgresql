package handler

import (
	"encoding/json"
	"fmt"
	"funding-app/app/campaign"
	"funding-app/app/helper"
	"funding-app/app/key"
	"funding-app/app/user"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
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

	uploadedFile, handler, err := r.FormFile("image")
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

	filename := fmt.Sprintf("%s-%s", input.CampaignID, handler.Filename)
	fileLocation := fmt.Sprintf("files/%s", filename)

	_, err = h.campaignService.UploadCampaignImage(input, fileLocation)
	if err != nil {
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", err.Error())
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 06666)
	if err != nil {
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", err.Error())
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	defer targetFile.Close()

	_, err = io.Copy(targetFile, uploadedFile)
	if err != nil {
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", err.Error())
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	log.Info("Success upload campaign image!")

	data := M{
		"is_uploaded": true,
	}

	response := helper.APIResponse("Success upload campaign image", http.StatusCreated, "success", data)
	helper.JSON(w, response, http.StatusCreated)
}
