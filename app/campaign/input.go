package campaign

import "funding-app/app/user"

type (
	CreateCampaignInput struct {
		Name             string `json:"name" validate:"required"`
		ShortDescription string `json:"short_description" validate:"required"`
		Description      string `json:"description" validate:"required"`
		Perks            string `json:"perks" validate:"required"`
		GoalAmount       int    `json:"goal_amount" validate:"required"`
		User             user.User
	}

	CreateCampaignImageInput struct {
		CampaignID string `form:"campaign_id"`
		IsPrimary  bool   `form:"is_primary"`
		User       user.User
	}
)
