package image

import (
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type ImageBucketInterface interface {
	UploadImage(prefix string, file *multipart.FileHeader) (string, error)
	DeleteImage(publicID string) error
	DeleteImagesByPrefix(prefix string) error
}

type ImageBucket struct {
	CloudinaryClient *cloudinary.Cloudinary
}

func NewImageBucket(cld *cloudinary.Cloudinary) ImageBucketInterface {
	return &ImageBucket{CloudinaryClient: cld}
}

func (b *ImageBucket) UploadImage(prefix string, file *multipart.FileHeader) (string, error) {
	ctx := context.Background()
	resp, err := b.CloudinaryClient.Upload.Upload(ctx, prefix+file.Filename, uploader.UploadParams{})
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

func (b *ImageBucket) DeleteImagesByPrefix(prefix string) error {
	prefixes := []string{prefix}
	ctx := context.Background()
	_, err := b.CloudinaryClient.Admin.DeleteAssetsByPrefix(ctx, admin.DeleteAssetsByPrefixParams{Prefix: prefixes})
	return err
}
