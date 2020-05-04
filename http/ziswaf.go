package http

import (
	_APICall "uangbaik-account-microservice/config"
)

type campaignModel struct {
	ID             uint   `gorm:"primary_key" json:"id"`
	OrganizationID uint   `gorm:"not_null" json:"organization_id"`
	Title          string `gorm:"size:255" json:"title"`
	Description    string `gorm:"size:255" json:"description"`
	DonationGoals  int64  `gorm:"default:0" json:"donation_goals"`
	Status         int    `gorm:"default:1" json:"status"`
}

type organizationModel struct {
	ID     uint   `gorm:"primary_key" json:"id"`
	UserID uint   `gorm:"not_null" json:"user_id"`
	Name   string `gorm:"size:255" json:"name"`
	Status int    `gorm:"default:1" json:"status"`
}

type campaignResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		TotalCount     int             `json:"total_count"`
		TotalSuspended int             `json:"total_suspended"`
		Attributes     []campaignModel `json:"attributes"`
	} `json:"data"`
}

type organizationResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		TotalCount     int                 `json:"total_count"`
		TotalSuspended int                 `json:"total_suspended"`
		Attributes     []organizationModel `json:"attributes"`
	} `json:"data"`
}

type Ziswaf struct {
	AccessToken string
	ZiswafAPI   _APICall.Request
}

// ZiswafData methods
type ZiswafData interface {
	GetCampaign() (*campaignResponse, error)
	GetOrganization() (*organizationResponse, error)
	GetToken() string
	SetToken(string)
}

// NewZiswaf baseurl
func NewZiswaf(nBaseURL string) ZiswafData {
	return &Ziswaf{
		ZiswafAPI: _APICall.NewAPICall(nBaseURL),
	}
}

func (r *Ziswaf) SetToken(AuthToken string) {
	r.AccessToken = AuthToken
}
func (r *Ziswaf) GetToken() string {
	return r.AccessToken
}

func (r *Ziswaf) GetCampaign() (*campaignResponse, error) {
	var (
		APIEndpoint  string
		CampaignResp campaignResponse
	)

	APIEndpoint = "admin/campaign"

	resp, err := r.ZiswafAPI.APICall("GET", APIEndpoint, r.GetToken(), nil)
	if err != nil {
		panic(err)
	}
	resp.ToJSON(&CampaignResp)

	return &CampaignResp, err
}

func (r *Ziswaf) GetOrganization() (*organizationResponse, error) {
	var (
		APIEndpoint      string
		OrganizationResp organizationResponse
	)

	APIEndpoint = "/admin/organizations"

	resp, err := r.ZiswafAPI.APICall("GET", APIEndpoint, r.GetToken(), nil)
	if err != nil {
		panic(err)
	}
	resp.ToJSON(&OrganizationResp)

	return &OrganizationResp, err
}
