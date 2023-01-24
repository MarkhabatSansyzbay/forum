package service

import (
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

var (
	ErrImgSize   = errors.New("image file size is too big")
	ErrImgFormat = errors.New("your image format is not provided") // provided image formats are JPEG, PNG, SVG and GIF
)

const imgMaxSize = 5 << 20

// SaveImages reads image files, saves it and return its base64 encoding as a slice
func SaveImages(images []*multipart.FileHeader) ([]template.URL, error) {
	paths := make([]template.URL, len(images))
	for i, fileHeader := range images {
		if fileHeader.Size > imgMaxSize {
			return nil, ErrImgSize
		}

		img, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}
		defer img.Close()

		f, err := os.Create(fmt.Sprintf("./uploads/%d%s", time.Now().UnixNano(), fileHeader.Filename))
		if err != nil {
			return nil, err
		}

		defer f.Close()

		_, err = io.Copy(f, img)
		if err != nil {
			return nil, err
		}

		content, err := os.ReadFile(f.Name())
		if err != nil {
			return nil, err
		}

		filetype := http.DetectContentType(content)
		if filetype != "image/jpeg" && filetype != "image/png" && filetype != "image/gif" && filetype != "image/svg" {
			return nil, ErrImgFormat
		}

		path, err := retrieveBase64(content, filetype)
		if err != nil {
			return nil, err
		}
		paths[i] = template.URL(path)
	}
	return paths, nil
}

func retrieveBase64(imgBytes []byte, mimeType string) (string, error) {
	var base64Encoding string

	switch mimeType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	}

	base64Encoding += toBase64(imgBytes)
	return base64Encoding, nil
}

func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
