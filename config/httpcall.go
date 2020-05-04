package config

import (
	"github.com/imroc/req"
)

type request struct {
	baseURL string
}

//Request methods to make http request
type Request interface {
	APICall(string, string, string, interface{}) (*req.Resp, error)
}

//NewAPICall create instance of HTTPCall
func NewAPICall(nBaseURL string) Request {
	return &request{
		baseURL: nBaseURL,
	}
}

// APICall middleware to make call http request
func (h *request) APICall(nMethod string, nURL string, nAccessToken string, nBody interface{}) (*req.Resp, error) {

	APIURL := h.baseURL + nURL

	req.Debug = true

	reqHeader := req.Header{
		"Accept":  "application/json",
		"Authorization": nAccessToken,
	}

	switch nMethod {
	case "GET":
		resp, err := req.Get(APIURL, reqHeader, nBody)
		return resp, err
	case "POST":
		resp, err := req.Post(APIURL, reqHeader, nBody)
		return resp, err
	}
	return nil, nil
}
