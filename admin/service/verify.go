package service

import (
	"fmt"
	"strconv"
	"uangbaik-account-microservice/database/models"
	_http "uangbaik-account-microservice/http"
	_u "uangbaik-account-microservice/utils"

	"net/http"

	"github.com/gin-gonic/gin"
)

// VerifyStatus for processing request from client
func (idb *InDB) VerifyStatus(c *gin.Context) {
	var (
		// merchant models.Merchant
		user models.User
	)
	concern := c.Request.FormValue("concern")
	request := struct {
		Data struct {
			ID     uint   `json:"id"`
			Status string `json:"status"`
		} `json:"data"`
	}{}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusUnprocessableEntity, _u.Message(false, "Error while decoding request body", err.Error(), nil))
		return
	}

	req := _http.NewHttpService()

	switch concern {
	case "merchant":
		// if err := idb.DB.Where(models.Merchant{ID: request.Data.ID}).First(&merchant).Error; err != nil {
		// 	c.JSON(http.StatusBadRequest, _u.Message(false, err.Error(), nil, nil))
		// 	return
		// }

		// Get merchant data
		req.MerchantServ.SetToken(c.GetHeader("Authorization"))
		merchant, err := req.MerchantServ.GetMerchantDetail(request.Data.ID, 0, "", 0)
		if err != nil {
			c.JSON(http.StatusBadRequest, _u.Message(false, "Failed to get merchant", nil, nil))
			return
		}

		if merchant.Status == false {
			c.JSON(http.StatusBadRequest, _u.Message(false, merchant.Message, nil, nil))
			return
		}

		req.MerchantServ.SetToken(c.GetHeader("Authorization"))
		merchantUpdate, err1 := req.MerchantServ.UpdateMerchant(request.Data.ID)
		if err1 != nil {
			c.JSON(http.StatusBadRequest, _u.Message(false, "Failed to update merchant", nil, nil))
			return
		}

		strID := strconv.FormatUint(uint64(user.ID), 10)
		notifMsg := fmt.Sprintf("Merchant %s berhasil diaktivasi", merchant.Data.Attributes[0].Name)
		notif := _u.NewNotification(strID, "Ya Baik", notifMsg, nil)
		notif.Send()

		// Post Notification
		req.NotificationServ.SetToken(c.GetHeader("Authorization"))
		_, err2 := req.NotificationServ.PostNotification(strID, notifMsg, nil)
		if err2 != nil {
			c.JSON(http.StatusBadRequest, _u.Message(false, "Failed to post notification", nil, nil))
			return
		}

		c.JSON(http.StatusCreated, _u.Message(true, "Status merchant updated", &merchantUpdate, nil))
		return
	case "user":
		if err := idb.DB.Where(models.User{ID: request.Data.ID}).First(&user).Error; err != nil {
			c.JSON(http.StatusBadRequest, _u.Message(false, err.Error(), nil, nil))
			return
		}

		user.Status = 2 // 2 means user profile verified by admin
		idb.DB.Save(&user)

		strID := strconv.FormatUint(uint64(user.ID), 10)
		notifMsg := fmt.Sprintf("Akun %s berhasil diaktivasi", user.Phone)
		notif := _u.NewNotification(strID, "Ya Baik", notifMsg, nil)
		notif.Send()

		// Post Notification
		req.NotificationServ.SetToken(c.GetHeader("Authorization"))
		_, err2 := req.NotificationServ.PostNotification(strID, notifMsg, nil)
		if err2 != nil {
			c.JSON(http.StatusBadRequest, _u.Message(false, "Failed to post notification", nil, nil))
			return
		}
		// models.NewNotification(strID, notifMsg, nil)

		c.JSON(http.StatusCreated, _u.Message(true, "Status user updated", &user, nil))
		return
	default:
		c.JSON(http.StatusUnprocessableEntity, _u.Message(true, "Wrong concern type", nil, nil))
		return
	}
}
