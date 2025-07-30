package imagekit

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/imagekit-developer/imagekit-go/api/media"
)

// DeleteImageByUrl deletes an image from ImageKit using its URL
func DeleteImageByUrl(url string) error {
	if url == "" {
		return errors.New("empty image URL")
	}
	urlEndpoint := os.Getenv("IMAGEKIT_URL_ENDPOINT")
	// Extract the image path relative to the URL endpoint
	trimmedPath := strings.TrimPrefix(url, urlEndpoint)
	trimmedPath = strings.TrimPrefix(trimmedPath, "/")

	files, err := ImageKit.Media.Files(context.Background(), media.FilesParam{
		SearchQuery: "name=\"" + trimmedPath + "\"",
	})

	if err != nil {
		return err
	}

	if len(files.Data) == 0 {
		return errors.New("file not found on ImageKit")
	}

	fileID := files.Data[0].FileId

	// Step 2: Delete the file by file ID
	_, err = ImageKit.Media.DeleteFile(context.Background(), fileID)
	if err != nil {
		return err
	}

	return nil
}
