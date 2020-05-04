package http

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	_APICall "uangbaik-account-microservice/config"

	"github.com/imroc/req"
)

type merchantDetailModel struct {
	ID               uint            `gorm:"primary_key" json:"id"`
	UUID             string          `json:"uuid"`
	UserID           uint            `gorm:"not null" json:"user_id"`
	Name             string          `gorm:"size:255" json:"name"`
	Phone            string          `gorm:"size:255;unique_index;not null" json:"phone"`
	Email            string          `gorm:"size:255" json:"email"`
	Address          string          `gorm:"size:255" json:"address"`
	AddressMap       string          `gorm:"size:255" json:"address_map"`
	BusinessCategory string          `gorm:"size:255" json:"business_category"`
	Hash             string          `gorm:"size:255" json:"hash"`
	Status           int             `gorm:"default:0" json:"status"`
	MerchantImage    string          `gorm:"size:255" json:"merchant_image"`
	Ewallet          json.RawMessage `json:"e_wallet"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
	DeletedAt        *time.Time      `json:"deleted_at"`
}

type merchantEditResponse struct {
	Status  bool                `json:"status"`
	Message string              `json:"message"`
	Data    merchantDetailModel `json:"data"`
}

type merchantDetailResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Attributes []merchantDetailModel `json:"attributes"`
		Total      struct {
			Count           int `json:"count"`
			UnverifiedCount int `json:"unverified_count"`
			BannedCount     int `json:"banned_count"`
		} `json:"total"`
	} `json:"data"`
}

type merchantDetail struct {
	AccessToken       string
	MerchantDetailAPI _APICall.Request
}

// MerchantDetail methods to interact with account service
type MerchantDetail interface {
	GetMerchantDetail(uint, uint, string, int) (*merchantDetailResponse, error)
	UpdateMerchant(uint) (*merchantEditResponse, error)
	GetToken() string
	SetToken(string)
}

// NewMerchantDetail baseurl
func NewMerchantDetail(nBaseURL string) MerchantDetail {
	return &merchantDetail{
		MerchantDetailAPI: _APICall.NewAPICall(nBaseURL),
	}
}

func (r *merchantDetail) SetToken(AuthToken string) {
	r.AccessToken = AuthToken
}
func (r *merchantDetail) GetToken() string {
	return r.AccessToken
}

func (r *merchantDetail) UpdateMerchant(nID uint) (*merchantEditResponse, error) {
	var (
		APIEndpoint        string
		merchantUpdateResp merchantEditResponse
	)

	APIEndpoint = "merchants/" + strconv.FormatUint(uint64(nID), 10)

	resp, err := r.MerchantDetailAPI.APICall("POST", APIEndpoint, r.GetToken(), nil)
	if err != nil {
		panic(err)
	}
	resp.ToJSON(&merchantUpdateResp)

	return &merchantUpdateResp, err
}

func (r *merchantDetail) GetMerchantDetail(nID uint, nUserID uint, uuidMerchant string, nStatus int) (*merchantDetailResponse, error) {
	var (
		APIEndpoint        string
		merchantDetailResp merchantDetailResponse
		param              interface{}
	)

	APIEndpoint = "admin/merchants/"

	if nUserID != 0 {
		param = req.Param{
			"user_id": nUserID,
		}
	} else if uuidMerchant != "" {
		param = req.Param{
			"uuid": uuidMerchant,
		}
	} else if nID != 0 && nStatus == 0 {
		param = req.Param{
			"id": nID,
		}
	} else if nID != 0 && nStatus != 0 {
		fmt.Println("masuk1")
		param = req.Param{
			"id":     nID,
			"status": nStatus,
		}
	} else if nID == 0 {
		fmt.Println("masuk1", nID)
		param = req.Param{
			//empty
		}
	}

	resp, err := r.MerchantDetailAPI.APICall("GET", APIEndpoint, r.GetToken(), param)
	if err != nil {
		panic(err)
	}
	resp.ToJSON(&merchantDetailResp)

	return &merchantDetailResp, err
}
