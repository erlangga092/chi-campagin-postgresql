package helper

import (
	"context"
	"funding-app/app/key"
	"os"
	"sync"
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

func ImageUploadAvatarHandler(wg *sync.WaitGroup, input interface{}, fileResponse chan key.FileUploadResponse) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	defer wg.Done()

	// create cloudinary instance
	cld, err := cloudinary.NewFromParams(envCloudName, envAPIKey, envAPISecret)
	if err != nil {
		fileResponse <- key.FileUploadResponse{
			SecureURL: "",
			Err:       err,
		}
	}

	uploadParam, err := cld.Upload.Upload(ctx, input, uploader.UploadParams{Folder: envUploadFolderAvatar})
	if err != nil {
		fileResponse <- key.FileUploadResponse{
			SecureURL: "",
			Err:       err,
		}
	}

	fileResponse <- key.FileUploadResponse{
		SecureURL: uploadParam.SecureURL,
		Err:       nil,
	}
}

func ImageUploadCampaignImageHandler(wg *sync.WaitGroup, input interface{}, fileResponse chan key.FileUploadResponse) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	defer wg.Done()

	// create cloudinary instance
	cld, err := cloudinary.NewFromParams(envCloudName, envAPIKey, envAPISecret)
	if err != nil {
		fileResponse <- key.FileUploadResponse{
			SecureURL: "",
			Err:       err,
		}
	}

	uploadParam, err := cld.Upload.Upload(ctx, input, uploader.UploadParams{Folder: envUploadFolderCampaignImage})
	if err != nil {
		fileResponse <- key.FileUploadResponse{
			SecureURL: "",
			Err:       err,
		}
	}

	fileResponse <- key.FileUploadResponse{
		SecureURL: uploadParam.SecureURL,
		Err:       nil,
	}
}
