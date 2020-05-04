package utils

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SavedFile struct {
	FileName string
	FilePath string
	FileURL  string
	MimeType string
	Concern  string
}

// SaveFile to save image get from request
func SaveFile(c *gin.Context, param string, concern string) (*SavedFile, error) {
	paramFile := param
	if param == "" {
		paramFile = "attachment_files"
	}
	file, errFiles := c.FormFile(paramFile)

	if file == nil {
		return &SavedFile{}, fmt.Errorf("Empty File: %s is empty", concern)
	}

	if errFiles != nil {
		return &SavedFile{}, errFiles
	}
	mimeType := file.Header.Get("Content-Type")
	if !strings.HasPrefix(mimeType, "image/") {
		return &SavedFile{}, errors.New("Only image file allowed")
	}

	folderPath := "../attachments/images/"
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		os.MkdirAll(folderPath, 0755)
	}

	extfile := strings.Split(file.Filename, ".")
	randFileName := uuid.New()
	file.Filename = randFileName.String() + concern + "." + extfile[len(extfile)-1]

	path := folderPath + file.Filename

	var host string
	if os.Getenv("GO_ENV") == "production" {
		host = "https://api.yabaik.id"
	} else {
		host = "http://" + c.Request.Host
	}
	url := host + "/api/v1/images/" + file.Filename

	if err := c.SaveUploadedFile(file, path); err != nil {
		return &SavedFile{}, err
	}

	return &SavedFile{
		FileName: file.Filename,
		FilePath: path,
		FileURL:  url,
		MimeType: file.Header.Get("Content-Type"),
		Concern:  concern,
	}, nil
}
