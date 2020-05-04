package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"uangbaik-account-microservice/database/models"
	_http "uangbaik-account-microservice/http"
	_request "uangbaik-account-microservice/setting/http/request"
	_u "uangbaik-account-microservice/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserInfo struct {
	UserID   uint   `json:"user_id"`
	FullName string `json:"full_name"`
}

type data struct {
	*models.User
	EWallet interface{} `json:"e_wallet"`
}

// GetMe get user data
func (idb *InDB) GetMe(c *gin.Context) {
	var (
		user     models.User
		userInfo UserInfo
	)

	userData := c.GetHeader("User-Info")

	if userData == "" {
		c.JSON(http.StatusBadRequest, _u.Message(false, "User not found", nil, nil))
		return
	}

	err := json.Unmarshal([]byte(userData), &userInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, _u.Message(false, err.Error(), nil, nil))
		return
	}

	if err := idb.DB.Preload("BankAccount").Preload("UploadAsset").Where(models.User{ID: userInfo.UserID}).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, _u.Message(false, err.Error(), nil, nil))
		return
	}

	// Check and get wallet user
	req := _http.NewHttpService()
	req.TransactionServ.SetToken(c.GetHeader("Authorization"))
	walletUser, err := req.TransactionServ.GetWallet(user.ID, "users")
	if err != nil {
		c.JSON(http.StatusBadRequest, _u.Message(false, "Failed to get wallet user", nil, nil))
		return
	}

	if walletUser.Status == false {
		c.JSON(http.StatusBadRequest, _u.Message(false, walletUser.Message, nil, nil))
		return
	}

	c.JSON(http.StatusAccepted, _u.Message(true, "User fetched", data{
		User:    &user,
		EWallet: walletUser.Data[0],
	}, nil))
}

// EditProfile for edit profile
func (idb *InDB) EditProfile(c *gin.Context) {
	var (
		user     models.User
		userInfo UserInfo
		bank     models.BankAccount
	)

	userData := c.GetHeader("User-Info")

	if userData == "" {
		c.JSON(http.StatusBadRequest, _u.Message(false, "User not found", nil, nil))
		return
	}

	err := json.Unmarshal([]byte(userData), &userInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, _u.Message(false, err.Error(), nil, nil))
		return
	}

	concern := c.Query("concern")
	if concern == "" {
		c.JSON(http.StatusUnprocessableEntity, _u.Message(false, "Concern is undefined", nil, nil))
		return
	}

	if concern == "change_pin" {
		accountKit, errExchange := _u.ExchangeAccountKit(c.Query("code_access_token"))
		if errExchange != nil {
			c.JSON(http.StatusBadRequest, _u.Message(true, errExchange.Error(), nil, nil))
			return
		}

		if err := idb.DB.Preload("BankAccount").Where(models.User{Phone: accountKit.Phone.NationalNumber, ID: userInfo.UserID}).First(&user).Error; err != nil {
			c.JSON(http.StatusBadRequest, _u.Message(false, err.Error(), nil, nil))
			return
		}
	} else if concern == "bank_account" {
		changeBankValidator := _request.ChangeBankValidator()
		if err := changeBankValidator.BindChangeBank(c); err != nil {
			c.JSON(http.StatusBadRequest, _u.NewValidatorError(err))
			return
		}

		cv := &changeBankValidator.Data

		bankID, _ := strconv.ParseUint(c.Query("id"), 10, 32)

		if err := idb.DB.Where(models.BankAccount{ID: uint(bankID)}).First(&bank).Error; err != nil {
			c.JSON(http.StatusBadRequest, _u.Message(false, "Bank Not Found", nil, nil))
			return
		}

		if err := idb.DB.Model(&bank).Where(models.BankAccount{ID: uint(bankID)}).Update(&cv).Error; err != nil {
			c.JSON(http.StatusBadRequest, _u.Message(false, err.Error(), &cv, nil))
			return
		}

		if err := idb.DB.Model(&bank).Where(models.BankAccount{ID: uint(bankID)}).Update(&cv.BankName).Error; err != nil {
			c.JSON(http.StatusBadRequest, _u.Message(false, err.Error(), &cv, nil))
			return
		}

		idb.DB.Where(&models.BankAccount{ID: uint(bankID)}).First(&bank)

		c.JSON(http.StatusAccepted, _u.Message(true, "Bank has been updated", &bank, nil))
		return
	}

	// Separate by concern type
	if concern == "profile" {
		profileValidator := _request.EditUserValidator()
		if err := profileValidator.BindEditUser(c); err != nil {
			c.JSON(http.StatusBadRequest, _u.NewValidatorError(err))
			return
		}

		cv := &profileValidator

		if err := idb.DB.Model(&user).Where(models.User{ID: userInfo.UserID}).Update(&cv.Data).Error; err != nil {
			c.JSON(http.StatusBadRequest, _u.Message(false, err.Error(), &cv, nil))
			return
		}

		idb.DB.Where(&models.User{ID: userInfo.UserID}).First(&user)
		c.JSON(http.StatusAccepted, _u.Message(true, "User has been updated", user, nil))
		return
	} else if concern == "change_pin" {
		changePinValidator := _request.ChangePinValidator()
		if err := changePinValidator.BindChangePin(c); err != nil {
			c.JSON(http.StatusBadRequest, _u.NewValidatorError(err))
			return
		}

		cv := &changePinValidator

		byteDBPin := []byte(user.Pin)
		byteReqPin := []byte(cv.Data.OldPin)

		if errorPin := bcrypt.CompareHashAndPassword(byteDBPin, byteReqPin); errorPin != nil {
			c.JSON(http.StatusBadRequest, _u.Message(false, "Wrong Old Pin Number", nil, nil))
			return
		}

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(cv.Data.NewPin), bcrypt.DefaultCost)
		user.Pin = string(hashedPassword)

		if err := idb.DB.Save(&user).Error; err != nil {
			c.JSON(http.StatusBadRequest, _u.Message(false, err.Error(), nil, nil))
			return
		}

		c.JSON(http.StatusAccepted, _u.Message(true, "Pin updated", nil, nil))
		return
	} else {
		c.JSON(http.StatusUnprocessableEntity, _u.Message(false, "concern type unrecognised", nil, nil))
	}
}

// Upload for verified document
func (idb *InDB) Upload(c *gin.Context) {
	var (
		Asset    models.UploadAsset
		userInfo UserInfo
	)

	userData := c.GetHeader("User-Info")

	if userData == "" {
		c.JSON(http.StatusBadRequest, _u.Message(false, "User not found", nil, nil))
		return
	}

	err := json.Unmarshal([]byte(userData), &userInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, _u.Message(false, err.Error(), nil, nil))
		return
	}

	concern := c.Query("concern")
	if concern == "" {
		c.JSON(http.StatusUnprocessableEntity, _u.Message(false, "concern is undefined", nil, nil))
		return
	}

	saveFile, errFile := _u.SaveFile(c, "files", "user_"+concern)
	if errFile != nil {
		c.JSON(http.StatusUnprocessableEntity, _u.Message(false, errFile.Error(), nil, nil))
		return
	}

	Asset.Path = saveFile.FilePath
	Asset.FileName = saveFile.FileName
	Asset.URL = saveFile.FileURL
	Asset.MimeType = saveFile.MimeType
	Asset.UploadID = userInfo.UserID
	Asset.UploadType = "users"
	Asset.Type = concern

	saveDB := idb.DB.Create(&Asset).Error
	if saveDB != nil {
		c.JSON(http.StatusBadRequest, _u.Message(false, saveDB.Error(), nil, nil))
		return
	}
	c.JSON(http.StatusCreated, _u.Message(true, "Image saved success", Asset, nil))
}

// CreateBankAccount for create bank account
func (idb *InDB) CreateBankAccount(c *gin.Context) {
	var (
		user        models.User
		bankAccount models.BankAccount
		bank        models.Bank
		userInfo    UserInfo
	)

	userData := c.GetHeader("User-Info")

	if userData == "" {
		c.JSON(http.StatusBadRequest, _u.Message(false, "User not found", nil, nil))
		return
	}

	err := json.Unmarshal([]byte(userData), &userInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, _u.Message(false, err.Error(), nil, nil))
		return
	}

	concern := c.Query("concern")

	bankValidator := _request.BankValidator()
	if err := bankValidator.BindCreateBank(c); err != nil {
		c.JSON(http.StatusBadRequest, _u.NewValidatorError(err))
		return
	}

	cv := &bankValidator

	if err := idb.DB.Where(models.User{ID: userInfo.UserID}).First(&user).Error; err != nil {
		c.JSON(http.StatusUnprocessableEntity, _u.Message(false, err.Error(), nil, nil))
		return
	}

	if concern == "users" {
		bankAccount.ActorID = user.ID
	} else if concern == "merchants" {

		// Check and get Merchant
		req := _http.NewHttpService()
		req.MerchantServ.SetToken(c.GetHeader("Authorization"))
		merchantResponse, err := req.MerchantServ.GetMerchantDetail(0, user.ID, "", 0)
		if err != nil {
			c.JSON(http.StatusBadRequest, _u.Message(false, err.Error(), nil, nil))
			return
		}

		if len(merchantResponse.Data.Attributes) == 0 {
			c.JSON(http.StatusBadRequest, _u.Message(false, "Merchant Not Found", nil, nil))
			return
		}

		if merchantResponse.Status == false {
			c.JSON(http.StatusBadRequest, _u.Message(false, merchantResponse.Message, nil, nil))
			return
		} else {
			bankAccount.ActorID = merchantResponse.Data.Attributes[0].ID
		}
	}

	if err := idb.DB.Where(models.Bank{ID: cv.Data.BankID}).First(&bank).Error; err != nil {
		c.JSON(http.StatusBadRequest, _u.Message(false, "Bank Not Found", nil, nil))
		return
	}

	bankAccount.ActorType = concern
	bankAccount.BankName = bank.Name
	bankAccount.AccountNumber = cv.Data.AccountNumber
	bankAccount.Status = 1
	idb.DB.Create(&bankAccount)

	c.JSON(http.StatusCreated, _u.Message(true, "Bank Account berhasil dibuat", &bankAccount, nil))
	return
}

//GetListBankAccount for
func (idb *InDB) GetListBankAccount(c *gin.Context) {
	var (
		bank      []models.BankAccount
		user      models.User
		actorType string
		actorID   uint
		userID    uint
		userInfo  UserInfo
	)

	concern := c.Query("concern")
	userData := c.GetHeader("User-Info")

	if userData == "" {
		c.JSON(http.StatusBadRequest, _u.Message(false, "User Not Found", nil, nil))
		return
	}

	if concern == "" {
		c.JSON(http.StatusBadRequest, _u.Message(false, "concern must be fill", nil, nil))
		return
	}

	err := json.Unmarshal([]byte(userData), &userInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, _u.Message(false, err.Error(), nil, nil))
		return
	}
	userID = userInfo.UserID

	if err := idb.DB.Where(models.User{ID: userID}).First(&user).Error; err != nil {
		c.JSON(http.StatusUnprocessableEntity, _u.Message(false, "User Not Found", err.Error(), nil))
		return
	}

	if concern == "merchants" {
		// Check and get Merchant
		req := _http.NewHttpService()
		req.MerchantServ.SetToken(c.GetHeader("Authorization"))
		merchant, err := req.MerchantServ.GetMerchantDetail(0, user.ID, "", 0)
		if err != nil {
			c.JSON(http.StatusBadRequest, _u.Message(false, "Failed to get merchant", nil, nil))
			return
		}

		actorType = concern

		if merchant.Status == false {
			c.JSON(http.StatusBadRequest, _u.Message(false, merchant.Message, nil, nil))
			return
		}
		actorID = merchant.Data.Attributes[0].ID
	} else if concern == "users" {
		actorType = concern
		actorID = user.ID
	}

	if err := idb.DB.Where("actor_id = ? AND actor_type = ?", actorID, actorType).Find(&bank).Error; err != nil {
		c.JSON(http.StatusBadRequest, _u.Message(false, err.Error(), nil, nil))
		return
	}
	c.JSON(http.StatusOK, _u.Message(true, "Bank Account successfully retrieved", &bank, nil))
}

// GetDetailBankAccount for
func (idb *InDB) GetDetailBankAccount(c *gin.Context) {
	var bank models.BankAccount

	if c.Param("id") == "" {
		c.JSON(http.StatusUnprocessableEntity, _u.Message(false, "Bank Account Id is undefined", nil, nil))
		return
	}

	bankID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, _u.Message(false, err.Error(), nil, nil))
		return
	}

	if c.Query("user_id") == "" {
		if err := idb.DB.Where("id = ?", bankID).First(&bank).Error; err != nil {
			c.JSON(http.StatusBadRequest, _u.Message(false, err.Error(), nil, nil))
			return
		}
	} else if c.Query("user_id") != "" {
		if err := idb.DB.Where("id = ? AND actor_id = ? AND actor_type = ?", bankID, c.Query("user_id"), "users").First(&bank).Error; err != nil {
			c.JSON(http.StatusBadRequest, _u.Message(false, err.Error(), nil, nil))
			return
		}
	}

	c.JSON(http.StatusOK, _u.Message(true, "Bank Account successfully retrived", &bank, nil))
}

//GetListBank
func (idb *InDB) GetListBank(c *gin.Context) {
	var bank []models.Bank

	if err := idb.DB.Find(&bank).Error; err != nil {
		c.JSON(http.StatusBadRequest, _u.Message(false, "Bank not found", nil, nil))
		return
	}

	c.JSON(http.StatusOK, _u.Message(true, "Bank data sucessfully retrived", &bank, nil))
}

// DeleteBankAccount
func (idb *InDB) DeleteBankAccount(c *gin.Context) {
	var (
		userInfo UserInfo
		bank     models.BankAccount
	)

	bankAccountID := c.Param("id")

	if bankAccountID == "" {
		c.JSON(http.StatusUnprocessableEntity, _u.Message(false, "Param must be required", nil, nil))
		return
	}

	if err := idb.DB.Where("id = ? AND actor_id = ? AND actor_type = ?", bankAccountID, userInfo.UserID, "users").First(&bank).Error; err != nil {
		c.JSON(http.StatusBadRequest, _u.Message(false, "Bank Account Not Found", nil, nil))
		return
	}

	if err := idb.DB.Where("id = ?", bankAccountID).Delete(&bank).Error; err != nil {
		c.JSON(http.StatusBadRequest, _u.Message(false, err.Error(), nil, nil))
		return
	}

	c.JSON(http.StatusOK, _u.Message(true, "Bank Account Successfully deleted", nil, nil))
}
