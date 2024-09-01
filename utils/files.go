package utils

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

func UploadSingleFileLocally(fileHeader *multipart.FileHeader) (string, error) {

	saveDir := "./uploads"
	if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
		return "", err
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}

	defer file.Close()
	nanoID, err := gonanoid.New()
	if err != nil {
		return "", err
	}

	newFileName := nanoID + "_" + fileHeader.Filename
	filePath := filepath.Join(saveDir, newFileName)

	outFile, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		return "", err
	}

	return newFileName, nil
}

func UploadMultipleFilesLocally(files []*multipart.FileHeader) ([]string, error) {
	var slice []string

	for _, fileHeader := range files {
		uploadedFileName, err := UploadSingleFileLocally(fileHeader)
		if err != nil {
			return slice, err
		}
		slice = append(slice, uploadedFileName)
	}

	return slice, nil
}
