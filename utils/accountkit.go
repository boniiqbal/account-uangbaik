package utils

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type AccountKit struct {
	ID    string `json:"id"`
	Phone struct {
		Number         string `json:"number"`
		CountryPrefix  string `json:"country_prefix"`
		NationalNumber string `json:"national_number"`
	} `json:"phone"`
	Application struct {
		ID string `json:"id"`
	} `json:"application"`
}

// ExchangeAccountKit to check account kit truth
func ExchangeAccountKit(codeAccessToken string) (*AccountKit, error) {
	e := godotenv.Load()
	if e != nil {
		log.Fatal(e)
		return &AccountKit{}, e
	}
	var accountKit AccountKit

	if codeAccessToken == "" {
		return &accountKit, errors.New("Account Kit Code is empty")
	}

	appID := os.Getenv("FACEBOOK_APP_ID")
	appSecret := os.Getenv("ACCOUNT_KIT_APP_SECRET")
	tokenExchangeBaseURL := "https://graph.accountkit.com/v1.3/access_token?grant_type=authorization_code&code=" + codeAccessToken + "&access_token=AA|" + appID + "|" + appSecret
	meEndpointBaseURL := "https://graph.accountkit.com/v1.3/me?access_token="

	resp, _ := http.Get(tokenExchangeBaseURL)
	if resp.StatusCode != 200 {
		return &accountKit, errors.New("Access token code invalid or expired")
	}

	requestBody := struct {
		ID                      string `json:"id"`
		AccessToken             string `json:"access_token"`
		TokenRefreshIntervalSec int    `json:"token_refresh_interval_sec"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&requestBody); err != nil {
		log.Fatal(err)
		return &accountKit, err
	}

	// If the code is valid and not expired, the Account Kit API will respond with a User Access Token.
	if requestBody.AccessToken != "" {
		resp, errGet := http.Get(meEndpointBaseURL + requestBody.AccessToken)
		if errGet != nil {
			log.Fatal(errGet)
			return &accountKit, errGet
		}

		if err := json.NewDecoder(resp.Body).Decode(&accountKit); err != nil {
			log.Fatal(err)
			return &accountKit, err
		}

		return &accountKit, nil
	}

	return &accountKit, errors.New("Access token code invalid or expired")
}
