package campaign

type (
	CampaignFormatter struct {
		ID               string `json:"id"`
		UserID           string `json:"user_id"`
		Name             string `json:"name"`
		ShortDescription string `json:"short_description"`
		ImageURL         string `json:"image_url"`
		CurrentAmount    int    `json:"current_amount"`
		GoalAmount       int    `json:"goal_amount"`
	}
)

func FormatCampaign(campaign Campaign) CampaignFormatter {
	formatter := CampaignFormatter{}
	formatter.ID = campaign.ID
	formatter.UserID = campaign.UserID
	formatter.Name = campaign.Name
	formatter.ShortDescription = campaign.ShortDescription
	formatter.ImageURL = ""
	formatter.CurrentAmount = campaign.CurrentAmount
	formatter.GoalAmount = campaign.GoalAmount

	return formatter
}

func FormatCampaigns(campaigns []Campaign) []CampaignFormatter {
	formatter := []CampaignFormatter{}

	for _, campaign := range campaigns {
		formatter = append(formatter, FormatCampaign(campaign))
	}

	return formatter
}
