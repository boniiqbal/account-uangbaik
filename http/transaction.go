package http

import (
	"net/url"
	"strconv"
	"time"
	_APICall "uangbaik-account-microservice/config"

	"github.com/imroc/req"
)

type walletDetail struct {
	ID              uint      `gorm:"primary_key" json:"id"`
	ActorID         uint      `json:"actor_id"`
	ActorType       string    `gorm:"size:45" json:"actor_type"`
	Saldo           int64     `json:"saldo"`
	LastTransaction time.Time `json:"last_transaction"`
	PendingWithdraw int64     `gorm:"-" json:"pending_withdraw"`
}

type walletResponse struct {
	Status  bool           `json:"status"`
	Message string         `json:"message"`
	Data    []walletDetail `json:"data"`
}

type WithdrawData struct {
	ID        uint   `gorm:"primary_key" json:"id"`
	ActorID   uint   `json:"actor_id"`
	ActorType string `gorm:"size:45" json:"actor_type"`
	Amount    int64  `json:"amount"`
	Status    string `gorm:"size:255" json:"status"`
}

type withdrawResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Withdraw struct {
			Total     int `json:"total"`
			Rejected  int `json:"rejected"`
			Pending   int `json:"pending"`
			Approved  int `json:"approved"`
			Disbursed int `json:"disbursed"`
		} `json:"withdraw"`
		Attributes []WithdrawData `json:"attributes"`
	} `json:"data"`
}

type feeSettingModel struct {
	ID     uint   `gorm:"primary_key" json:"id"`
	Amount int64  `json:"amount"`
	Type   string `json:"type"`
	Status int    `json:"status"`
}

type feeSettingResponse struct {
	Status  bool          `json:"status"`
	Message string          `json:"message"`
	Data    feeSettingModel `json:"data"`
}

type wallet struct {
	AccessToken string
	WalletAPI   _APICall.Request
}

// TransactionData methods to interact with wallet service
type TransactionData interface {
	GetFeeSetting(string) (*feeSettingResponse, error)
	CreateWallet(uint, string, int64, string) (*walletResponse, error)
	GetWallet(uint, string) (*walletResponse, error)
	GetWithdraw() (*withdrawResponse, error)
	GetToken() string
	SetToken(string)
}

// NewWalletRequest baseurl
func NewWalletRequest(nBaseURL string) TransactionData {
	return &wallet{
		WalletAPI: _APICall.NewAPICall(nBaseURL),
	}
}

func (r *wallet) SetToken(AuthToken string) {
	r.AccessToken = AuthToken
}
func (r *wallet) GetToken() string {
	return r.AccessToken
}

func (r *wallet) GetFeeSetting(nFeeType string) (*feeSettingResponse, error) {
	var (
		APIEndpoint    string
		FeeSettingResp feeSettingResponse
	)

	APIEndpoint = "/fee"

	param := req.Param{
		"fee_type": nFeeType,
	}

	resp, err := r.WalletAPI.APICall("GET", APIEndpoint, r.GetToken(), param)
	if err != nil {
		panic(err)
	}
	resp.ToJSON(&FeeSettingResp)

	return &FeeSettingResp, err
}

// GetWithdraw
func (r *wallet) GetWithdraw() (*withdrawResponse, error) {
	var (
		APIEndpoint  string
		withdrawResp withdrawResponse
	)

	APIEndpoint = "/admin/withdraw"

	resp, err := r.WalletAPI.APICall("GET", APIEndpoint, r.GetToken(), nil)
	if err != nil {
		panic(err)
	}
	resp.ToJSON(&withdrawResp)

	return &withdrawResp, err
}

func (r *wallet) GetWallet(nActorID uint, nActorType string) (*walletResponse, error) {
	var (
		APIEndpoint string
		walletResp  walletResponse
	)

	APIEndpoint = "/wallet"

	query := url.Values{}
	query.Add("actor_id", strconv.FormatUint(uint64(nActorID), 10))
	query.Add("actor_type", nActorType)

	if nActorID == 0 {
		query.Set("actor_id", "0")
	}

	resp, err := r.WalletAPI.APICall("GET", APIEndpoint, r.GetToken(), query)
	if err != nil {
		panic(err)
	}
	resp.ToJSON(&walletResp)

	return &walletResp, err
}

func (r *wallet) CreateWallet(nActor uint, nActorType string, nSaldo int64, method string) (*walletResponse, error) {
	var (
		walletResp  walletResponse
		APIEndpoint string
		param       interface{}
	)

	switch method {
	case "Post":
		param = req.Param{
			"actor_id":   nActor,
			"actor_type": nActorType,
			"saldo":      0,
		}

		APIEndpoint = "/wallet"
	case "Update":
		param = req.Param{
			"actor_id":   nActor,
			"actor_type": nActorType,
			"saldo":      nSaldo,
		}

		APIEndpoint = "/wallet-edit"
	}

	resp, err := r.WalletAPI.APICall("POST", APIEndpoint, r.GetToken(), param)
	if err != nil {
		panic(err)
	}

	resp.ToJSON(&walletResp)

	return &walletResp, err
}
