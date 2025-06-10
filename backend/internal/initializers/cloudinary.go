package initializers

import (
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/image"
)

var CloudinaryClient *cloudinary.Cloudinary
var ImageBucket image.ImageBucketInterface

func ConnectToCloudinary() {
	cloudinaryURL := os.Getenv("CLOUDINARY_URL")
	cld, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		log.Fatal("Connection with cloudinary couldn't be established.")
	}
	CloudinaryClient = cld
	ImageBucket = image.NewImageBucket(cld)
}
