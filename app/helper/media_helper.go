package helper

import (
	"context"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

var (
	envCloudName                 = os.Getenv("CLOUDINARY_CLOUD_NAME")
	envAPIKey                    = os.Getenv("CLOUDINARY_API_KEY")
	envAPISecret                 = os.Getenv("CLOUDINARY_API_SECRET")
	envUploadFolderAvatar        = os.Getenv("CLOUDINARY_UPLOAD_FOLDER_AVATAR")
	envUploadFolderCampaignImage = os.Getenv("CLOUDINARY_UPLOAD_FOLDER_CAMPAIGN_IMAGE")
)

func ImageUploadAvatarHandler(input interface{}) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// create cloudinary instance
	cld, err := cloudinary.NewFromParams(envCloudName, envAPIKey, envAPISecret)
	if err != nil {
		return "", err
	}

	uploadParam, err := cld.Upload.Upload(ctx, input, uploader.UploadParams{Folder: envUploadFolderAvatar})
	if err != nil {
		return "", err
	}

	return uploadParam.SecureURL, nil
}

func ImageUploadCampaignImageHandler(input interface{}) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// create cloudinary instance
	cld, err := cloudinary.NewFromParams(envCloudName, envAPIKey, envAPISecret)
	if err != nil {
		return "", err
	}

	uploadParam, err := cld.Upload.Upload(ctx, input, uploader.UploadParams{Folder: envUploadFolderCampaignImage})
	if err != nil {
		return "", err
	}

	return uploadParam.SecureURL, nil
}
