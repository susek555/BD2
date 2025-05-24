package image

import (
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type ImageBucketInterface interface {
	UploadImage(folder string, file *multipart.FileHeader) (string, error)
	DeleteImage(publicID string) error
	DeleteImagesByFolder(folder string) error
}

type ImageBucket struct {
	CloudinaryClient *cloudinary.Cloudinary
}

func NewImageBucket(cld *cloudinary.Cloudinary) ImageBucketInterface {
	return &ImageBucket{CloudinaryClient: cld}
}

func (b *ImageBucket) UploadImage(folder string, file *multipart.FileHeader) (string, error) {
	ctx := context.Background()
	openedFile, err := file.Open()
	if err != nil {
		return "", err
	}
	resp, err := b.CloudinaryClient.Upload.Upload(ctx, openedFile, uploader.UploadParams{Folder: folder})
	if err != nil {
		return "", err
	}
	return resp.SecureURL, nil
}

func (b *ImageBucket) DeleteImage(publicID string) error {
	ctx := context.Background()
	_, err := b.CloudinaryClient.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: publicID})
	return err
}

func (b *ImageBucket) DeleteImagesByFolder(folder string) error {
	ctx := context.Background()
	_, err := b.CloudinaryClient.Admin.DeleteFolder(ctx, admin.DeleteFolderParams{Folder: folder})
	return err
}
