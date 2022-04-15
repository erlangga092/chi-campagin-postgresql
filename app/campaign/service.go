package campaign

import (
	"context"
	"errors"
)

type Service interface {
	GetCampaigns(userID string) ([]Campaign, error)
	GetCampaignDetail(ID string) (Campaign, error)
}

type service struct {
	campaignRepository Repository
}

func NewCampaignService(campaignRepository Repository) Service {
	return &service{campaignRepository}
}

func (s *service) GetCampaigns(userID string) ([]Campaign, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if userID != "" {
		campaigns, err := s.campaignRepository.FindByUserID(ctx, userID)
		if err != nil {
			return campaigns, err
		}

		return campaigns, nil
	}

	campaigns, err := s.campaignRepository.FindAll(ctx)
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (s *service) GetCampaignDetail(ID string) (Campaign, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	campaign, err := s.campaignRepository.FindByID(ctx, ID)
	if err != nil {
		return campaign, err
	}

	if campaign.ID == "" {
		return campaign, errors.New("no campaign found")
	}

	return campaign, nil
}
