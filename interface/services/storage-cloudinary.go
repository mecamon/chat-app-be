package services

import (
	"context"
	"errors"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/mecamon/chat-app-be/config"
)

var storage *Cloudinary

type Cloudinary struct {
	service *cloudinary.Cloudinary
	app     *config.App
}

func InitStorage() (Storage, error) {
	app := config.GetConfig()
	cld, err := cloudinary.NewFromParams(app.StorageCloudName, app.StorageAPIKey, app.StorageAPISecret)
	if err != nil {
		return nil, err
	}
	storage = &Cloudinary{
		app:     app,
		service: cld,
	}
	return storage, nil
}

func GetStorage() (Storage, error) {
	if storage == nil {
		return nil, errors.New("instance not created yet. Make sure to call InitStorage() before calling this method")
	}
	return storage, nil
}

func (s *Cloudinary) UploadImage(file interface{}, filename string) (string, error) {
	uniqueFileName := true
	overwrite := true

	ctx := context.Background()
	uploadResult, err := storage.service.Upload.Upload(
		ctx,
		file,
		uploader.UploadParams{
			PublicID:       filename,
			UniqueFilename: &uniqueFileName,
			Folder:         s.app.StorageDirectory,
			Overwrite:      &overwrite,
		})
	if err != nil {
		return "", err
	}
	return uploadResult.SecureURL, nil
}

func (s *Cloudinary) DeleteImage(publicID string) (string, error) {
	ctx := context.Background()
	resp, err := storage.service.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID:     publicID,
		ResourceType: "image"})

	if err != nil {
		return "", err
	}

	return resp.Result, nil
}
