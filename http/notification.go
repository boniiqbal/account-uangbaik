package http

import (
	"bytes"
	"fmt"
	_APICall "uangbaik-account-microservice/config"

	"github.com/imroc/req"
)

type notification struct {
	Data map[string]interface{} `json:"data"`
}

type notificationResponse struct {
	Status  string       `json:"status"`
	Message string       `json:"message"`
	Data    notification `json:"data"`
}

type notificationDetail struct {
	AccessToken     string
	NotificationAPI _APICall.Request
}

// NotificationDetail methods to interact with account service
type Notification interface {
	PostNotification(string, string, map[string]interface{}) (*notificationResponse, error)
	GetToken() string
	SetToken(string)
}

// NewNotification baseurl
func NewNotification(nBaseURL string) Notification {
	return &notificationDetail{
		NotificationAPI: _APICall.NewAPICall(nBaseURL),
	}
}

func (r *notificationDetail) SetToken(AuthToken string) {
	r.AccessToken = AuthToken
}
func (r *notificationDetail) GetToken() string {
	return r.AccessToken
}

func (r *notificationDetail) PostNotification(nUserID string, message string, data map[string]interface{}) (*notificationResponse, error) {
	var (
		APIEndpoint      string
		notificationResp notificationResponse
	)

	APIEndpoint = "/notification"
	dataNotif := CreateKeyValuePairs(data)

	param := req.QueryParam{
		"user_id": nUserID,
		"message": message,
		"data":    dataNotif,
	}

	resp, err := r.NotificationAPI.APICall("POST", APIEndpoint, r.GetToken(), param)
	if err != nil {
		panic(err)
	}
	resp.ToJSON(&notificationResp)

	return &notificationResp, err
}

// CreateKeyValuePairs .
func CreateKeyValuePairs(m map[string]interface{}) string {
	if m != nil {
		b := new(bytes.Buffer)
		for key, value := range m {
			fmt.Fprintf(b, "%s=\"%s\"###", key, value)
		}
		toString := string(b.String()[:len(b.String())-3])
		return toString
	}
	return ""
}
