package imagekit

import (
	"os"

	"github.com/imagekit-developer/imagekit-go"
)

var ImageKit *imagekit.ImageKit

func Init() {
	privateKey := os.Getenv("IMAGEKIT_PRIVATE_KEY")
	publicKey := os.Getenv("IMAGEKIT_PUBLIC_KEY")
	urlEndpoint := os.Getenv("IMAGEKIT_URL_ENDPOINT")

	ImageKit = imagekit.NewFromParams(imagekit.NewParams{
		PrivateKey:  privateKey,
		PublicKey:   publicKey,
		UrlEndpoint: urlEndpoint,
	})
}
