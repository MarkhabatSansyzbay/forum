package service

import (
	"errors"
	"fmt"
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

func SaveImages(images []*multipart.FileHeader) error {
	for _, fileHeader := range images {
		if fileHeader.Size > imgMaxSize {
			return ErrImgSize
		}

		content, err := os.ReadFile(fileHeader.Filename)
		if err != nil {
			return err
		}

		filetype := http.DetectContentType(content)
		if filetype != "image/jpeg" && filetype != "image/png" && filetype != "image/gif" && filetype != "image/svg" {
			return ErrImgFormat
		}

		f, err := os.Create(fmt.Sprintf("./uploads/%d%s", time.Now().UnixNano(), fileHeader.Filename))
		if err != nil {
			return err
		}

		defer f.Close()

		img, err := fileHeader.Open()
		if err != nil {
			return err
		}
		defer img.Close()

		_, err = io.Copy(f, img)
		if err != nil {
			return err
		}
	}
	return nil
}
