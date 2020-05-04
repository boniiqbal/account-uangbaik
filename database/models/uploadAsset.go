package models

import (
	"errors"
	"time"

	_config "uangbaik-account-microservice/config"
)

//UploadAsset for models Asset
type UploadAsset struct {
	ID         uint       `gorm:"primary_key" json:"id"`
	UploadID   uint       `json:"upload_id"`
	UploadType string     `gorm:"size:255" json:"upload_type"`
	Path       string     `gorm:"size:255" json:"path"`
	FileName   string     `gorm:"size:255" json:"file_name"`
	URL        string     `gorm:"size:255" json:"url"`
	MimeType   string     `gorm:"size:255" json:"mime_type"`
	IsVerified int        `gorm:"default:0" json:"is_verified"`
	Type       string     `gorm:"size:255" json:"type"`
	Status     int        `gorm:"default:1" json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
}

// Validate incoming Upload Asset
func (uploadAsset *UploadAsset) Validate() (string, bool) {

	if uploadAsset.Path == "" {
		return "UploadAsset Path is required", false
	}
	if uploadAsset.FileName == "" {
		return "UploadAsset FileName is required", false
	}
	if uploadAsset.URL == "" {
		return "UploadAsset URL is required", false
	}
	if uploadAsset.MimeType == "" {
		return "UploadAsset MimeType is required", false
	}
	if uploadAsset.Type == "" {
		return "UploadAsset Type is required", false
	}
	if uploadAsset.UploadID == 0 {
		return "UploadAsset UploadID is required", false
	}
	if uploadAsset.UploadType == "" {
		return "UploadAsset UploadType is required", false
	}

	return "Requirement passed", true
}

// BeforeCreate . 
func (uploadAsset *UploadAsset) BeforeCreate() (err error) {
	if resp, ok := uploadAsset.Validate(); !ok {
		return errors.New(resp)
	}
	return
}

// Create upload asset
func (uploadAsset *UploadAsset) Create() (*UploadAsset, error) {

	if err := _config.GetDB().Create(uploadAsset).Error; err != nil {
		return &UploadAsset{}, err
	}

	return uploadAsset, nil
}
