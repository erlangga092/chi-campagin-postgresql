package campaign

import "time"

type Campaign struct {
	ID               string
	UserID           string
	Name             string
	ShortDescription string
	Description      string
	Slug             string
	Perks            string
	GoalAmount       int
	CurrentAmount    int
	BackerCount      int
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type CampaignImage struct {
	ID        string
	Name      string
	IsPrimary bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
