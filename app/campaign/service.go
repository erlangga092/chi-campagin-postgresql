package campaign

import (
	"context"
	"errors"
	"funding-app/app/helper"
	"funding-app/app/key"
	"mime/multipart"
	"strings"
	"sync"
)

type Service interface {
	GetCampaigns(userID string) ([]Campaign, error)
	GetCampaignDetail(ID string) (Campaign, error)
	CreateCampaign(input CreateCampaignInput) (Campaign, error)
	UploadCampaignImage(input CreateCampaignImageInput, uploadedFile multipart.File) (CampaignImage, error)
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

func (s *service) CreateCampaign(input CreateCampaignInput) (Campaign, error) {
	var campaign Campaign
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	campaign.ID = helper.GenerateID()
	campaign.UserID = input.User.ID
	campaign.Name = input.Name
	campaign.ShortDescription = input.ShortDescription
	campaign.Description = input.Description
	campaign.Perks = input.Perks
	campaign.GoalAmount = input.GoalAmount

	slugCandidate := strings.Join(strings.Split(strings.ToLower(campaign.Name), " "), "-")
	campaign.Slug = slugCandidate

	newCampaign, err := s.campaignRepository.Save(ctx, campaign)
	if err != nil {
		return newCampaign, err
	}

	return newCampaign, nil
}

func (s *service) UploadCampaignImage(input CreateCampaignImageInput, uploadedFile multipart.File) (CampaignImage, error) {
	var wg sync.WaitGroup
	campaignImage := CampaignImage{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan key.FileUploadResponse)
	defer close(ch)

	campaign, err := s.campaignRepository.FindByID(ctx, input.CampaignID)
	if err != nil {
		return campaignImage, err
	}

	if campaign.UserID != input.User.ID {
		return campaignImage, errors.New("not an owner of the campaign")
	}

	isPrimary := 0
	if input.IsPrimary {
		isPrimary = 1

		_, err := s.campaignRepository.MarkAllImageAsNonPrimary(ctx, input.CampaignID)
		if err != nil {
			return campaignImage, err
		}
	}

	campaignImage.ID = helper.GenerateID()
	campaignImage.CampaignID = input.CampaignID
	campaignImage.IsPrimary = isPrimary

	wg.Add(1)

	// make goroutine with passing channel
	go helper.ImageUploadCampaignImageHandler(&wg, uploadedFile, ch)
	fileResponse := <-ch

	wg.Wait()

	if fileResponse.Err != nil {
		return campaignImage, fileResponse.Err
	}

	campaignImage.FileName = fileResponse.SecureURL
	newCampaignImage, err := s.campaignRepository.SaveImage(ctx, campaignImage)
	if err != nil {
		return newCampaignImage, err
	}

	return newCampaignImage, nil
}
