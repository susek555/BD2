package image

import (
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type ImageBucketInterface interface {
	Upload(folder string, file *multipart.FileHeader) (string, string, error)
	Delete(publicID string) error
	DeleteByFolderName(folder string) error
}

type ImageBucket struct {
	CloudinaryClient *cloudinary.Cloudinary
}

func NewImageBucket(cld *cloudinary.Cloudinary) ImageBucketInterface {
	return &ImageBucket{CloudinaryClient: cld}
}

func (b *ImageBucket) Upload(folder string, file *multipart.FileHeader) (string, string, error) {
	ctx := context.Background()
	openedFile, err := file.Open()
	if err != nil {
		return "", "", err
	}
	resp, err := b.CloudinaryClient.Upload.Upload(ctx, openedFile, uploader.UploadParams{Folder: folder})
	if err != nil {
		return "", "", err
	}
	return resp.PublicID, resp.SecureURL, nil
}

func (b *ImageBucket) Delete(publicID string) error {
	ctx := context.Background()
	_, err := b.CloudinaryClient.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: publicID})
	return err
}

func (b *ImageBucket) DeleteByFolderName(folder string) error {
	ctx := context.Background()
	selectResp, err := b.CloudinaryClient.Admin.AssetsByAssetFolder(ctx, admin.AssetsByAssetFolderParams{AssetFolder: folder})
	if err != nil {
		return err
	}
	var publicIDs []string
	for _, asset := range selectResp.Assets {
		publicIDs = append(publicIDs, asset.PublicID)
	}
	if _, err = b.CloudinaryClient.Admin.DeleteAssets(ctx, admin.DeleteAssetsParams{PublicIDs: publicIDs}); err != nil {
		return err
	}
	_, err = b.CloudinaryClient.Admin.DeleteFolder(ctx, admin.DeleteFolderParams{Folder: folder})
	return err
}
