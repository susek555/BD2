package initializers

import (
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
)

var CloudinaryClient *cloudinary.Cloudinary

func ConnectToCloudinary() {
	cloudinaryURL := os.Getenv("CLOUDINARY_URL")
	cld, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		log.Fatal("Connection with cloudinary couldn't be established.")
	}
	CloudinaryClient = cld
}
