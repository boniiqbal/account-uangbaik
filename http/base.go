package http

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// HttpCallService struct
type HttpCallService struct {
	TransactionServ  TransactionData
	MerchantServ     MerchantDetail
	NotificationServ Notification
	ZiswafServ       ZiswafData
}

func NewHttpService() *HttpCallService {
	e := godotenv.Load()
	if e != nil {
		log.Print(e)
	}

	urlTransaction := os.Getenv("TRANSACTION_MICROSERVICE")
	urlNotification := os.Getenv("NOTIFICATION_MICROSERVICE")
	urlPos := os.Getenv("POS_MICROSERVICE")
	urlZiswaf := os.Getenv("ZISWAF_MICROSERVICE")

	nTransactionService := NewWalletRequest(urlTransaction)
	nMerchantService := NewMerchantDetail(urlPos)
	nNotificationService := NewNotification(urlNotification)
	nZiswafService := NewZiswaf(urlZiswaf)

	return &HttpCallService{
		TransactionServ:  nTransactionService,
		MerchantServ:     nMerchantService,
		NotificationServ: nNotificationService,
		ZiswafServ:       nZiswafService,
	}
}
