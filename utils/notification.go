package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

/** Notification helpers **/

// NotificationFilter filter
type NotificationFilter struct {
	Field    string `json:"field"`
	Key      string `json:"key"`
	Relation string `json:"relation"`
	Value    string `json:"value"`
}

// Notification struct
type Notification struct {
	Filters  []NotificationFilter   `json:"filters"`
	Headings map[string]interface{} `json:"headings"`
	Contents map[string]interface{} `json:"contents"`
	Data     interface{}            `json:"data"`
	AppID    string                 `json:"app_id"`
}

// Send notification
func (notification *Notification) Send() error {
	e := godotenv.Load()
	if e != nil {
		log.Fatalln(e)
		return e
	}
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	if notif, err := json.Marshal(notification); err == nil {
		req, errBuff := http.NewRequest("POST", "https://onesignal.com/api/v1/notifications", bytes.NewBuffer(notif))
		if errBuff != nil {
			log.Fatalln(errBuff.Error())
			return errBuff
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Basic %s", os.Getenv("ONE_SIGNAL_TOKEN")))
		log.Println(string(notif))
		resps, errReq := client.Do(req)
		if errReq != nil {
			log.Fatalln(errReq.Error())
			return errReq
		}

		body, errBody := ioutil.ReadAll(resps.Body)
		log.Println(string(body))
		if errBody != nil {
			log.Fatalln(errBody.Error())
			return errBody
		}
	} else {
		log.Fatalln(err)
		return err
	}

	return nil
}

// NewNotification send notification
func NewNotification(senderID string, title string, message string, data map[string]interface{}) Notification {
	payload := Notification{}
	e := godotenv.Load()
	if e != nil {
		log.Fatal(e)
	}
	var onesignalEnv string
	if os.Getenv("ONE_SIGNAL_ENV") == "production" {
		onesignalEnv = "production"
	} else {
		onesignalEnv = "development"
	}

	// Set client id
	payload.AppID = os.Getenv("ONE_SIGNAL_APPID")

	// Set filters
	defaultFilter := NotificationFilter{}
	defaultFilter.Field = "tag"
	defaultFilter.Key = "userId"
	defaultFilter.Relation = "="
	defaultFilter.Value = senderID

	// ENV filters
	envFilter := NotificationFilter{}
	envFilter.Field = "tag"
	envFilter.Key = "env"
	envFilter.Relation = "="
	envFilter.Value = onesignalEnv

	payload.Filters = append(payload.Filters, defaultFilter, envFilter)

	payload.Headings = make(map[string]interface{})
	payload.Contents = make(map[string]interface{})

	if data != nil {
		payload.Data = data
	} else {
		payload.Data = make(map[string]interface{})
	}

	payload.Headings["en"] = title
	payload.Contents["en"] = message

	return payload
}
