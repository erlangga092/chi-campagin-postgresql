package handler

import (
	"funding-app/app/campaign"
	"funding-app/app/helper"
	"net/http"

	"github.com/go-chi/chi/v5"
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

	response := helper.APIResponse("List of campaigns", http.StatusOK, "success", campaigns)
	helper.JSON(w, response, http.StatusOK)
}

func (h *campaignHandler) GetCampaignByID(w http.ResponseWriter, r *http.Request) {
	campaignID := chi.URLParam(r, "id")

	campaign, err := h.campaignService.GetCampaignByID(campaignID)
	if err != nil {
		response := helper.APIResponse("Failed to get campaign", http.StatusBadRequest, "error", err.Error())
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	response := helper.APIResponse("Detail of campaigns", http.StatusOK, "success", campaign)
	helper.JSON(w, response, http.StatusOK)
}
