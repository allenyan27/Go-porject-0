package imagekit

import (
	"bytes"
	"context"
	"mime/multipart"

	"github.com/imagekit-developer/imagekit-go/api/uploader"
)

func boolPtr(b bool) *bool {
	return &b
}

func UploadImageFromForm(file multipart.File, filename string) (string, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(file)
	if err != nil {
		return "", err
	}

	// Directly pass the bytes buffer (io.Reader)
	res, err := ImageKit.Uploader.Upload(context.Background(), buf, uploader.UploadParam{
		FileName:          filename,
		UseUniqueFileName: boolPtr(true),
	})
	if err != nil {
		return "", err
	}

	return res.Data.Url, nil
}

func UploadGalleryFromForm(files []*multipart.FileHeader) ([]string, error) {
	var urls []string

	for _, fileHeader := range files {
		// Open the file
		file, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}
		defer file.Close()

		// Read into buffer
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(file)
		if err != nil {
			return nil, err
		}

		// Upload to ImageKit
		res, err := ImageKit.Uploader.Upload(context.Background(), buf, uploader.UploadParam{
			FileName:          fileHeader.Filename,
			UseUniqueFileName: boolPtr(true),
		})
		if err != nil {
			return nil, err
		}

		// Append URL to result
		urls = append(urls, res.Data.Url)
	}

	return urls, nil
}
