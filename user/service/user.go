package service

import (
	"log"
	"fmt"
	"net/http"

	// "os"
	"strings"
	"time"

	"uangbaik-account-microservice/database/models"
	_http "uangbaik-account-microservice/http"
	_request "uangbaik-account-microservice/user/http/request"
	_u "uangbaik-account-microservice/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

//CreateUser for register new user
func (idb *InDB) CreateUser(c *gin.Context) {
	var (
		user models.User
	)

	accountKit, errExchange := _u.ExchangeAccountKit(c.PostForm("code_access_token"))
	if errExchange != nil {
		c.JSON(http.StatusBadRequest, _u.Message(true, errExchange.Error(), nil, nil))
		return
	}

	hashVariabel := time.Now().String() + accountKit.Phone.NationalNumber
	hashedPhone, errHash := bcrypt.GenerateFromPassword([]byte(hashVariabel), bcrypt.MinCost)
	if errHash != nil {
		log.Fatal(errHash)
		c.JSON(http.StatusBadRequest, _u.Message(false, errHash.Error(), nil, nil))
		return
	}

	if idb.DB.Where("phone = ?", accountKit.Phone.NationalNumber).First(&user).RecordNotFound() {
		user.Hash = string(hashedPhone)
		user.Phone = accountKit.Phone.NationalNumber
		user.IsPhoneVerified = 1

		if err := idb.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusBadRequest, _u.Message(false, err.Error(), nil, nil))
			return
		} //Create DB User

		c.JSON(http.StatusCreated, _u.Message(true, "Registration Complete", &user, nil))
	} else {
		c.JSON(http.StatusBadRequest, _u.Message(false, "Number Phone Already Exist", nil, nil))
	}
}

// LoginUser for login
func (idb *InDB) LoginUser(c *gin.Context) {
	var user models.User

	userValidator := _request.UserValidator()
	if err := userValidator.Bind(c); err != nil {
		c.JSON(http.StatusBadRequest, _u.NewValidatorError(err))
		return
	}

	cv := &userValidator

	checkPhone := idb.DB.Model(&user).Where(&models.User{Phone: cv.Data.Phone}).First(&user).Error

	if checkPhone != nil {
		c.JSON(http.StatusBadRequest, _u.Message(false, "Phone number not found", checkPhone, cv.Data.Phone))
	} else {
		c.JSON(http.StatusCreated, _u.Message(true, "Phone is correct", &user, nil))
	}
}

// RefreshToken for refresh access token
func (idb *InDB) RefreshToken(c *gin.Context) {
	var (
		user        models.User
		accessToken models.AccessToken
	)

	refreshTokenValidator := _request.RefreshTokenValidator()
	if err := refreshTokenValidator.BindRefreshToken(c); err != nil {
		c.JSON(http.StatusBadRequest, _u.NewValidatorError(err))
		return
	}

	cv := &refreshTokenValidator

	if err := idb.DB.Where(models.AccessToken{RefreshToken: cv.Data.RefreshToken}).First(&accessToken).Error; err != nil {
		c.JSON(http.StatusBadRequest, _u.Message(false, "Access Token not found", nil, nil))
		return
	}

	if err := idb.DB.Where(models.User{ID: accessToken.UserID}).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, _u.Message(false, "User not found", nil, nil))
		return
	}

	// split user agent
	splitted := strings.Split(c.Request.UserAgent(), " ")

	// Create the Claims
	claims := _u.CreateClaims("users", user.ID, splitted[0], 72)
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims) //Create Token
	signed, err := token.SignedString([]byte("secret"))

	if err != nil {
		c.JSON(http.StatusBadRequest, _u.Message(false, err.Error(), nil, nil))
		c.Abort()
		return
	}

	accessToken.AccessToken = signed
	accessToken.RefreshToken = _u.TokenGenerator(user.Phone)
	accessToken.UserAgent = splitted[0]
	accessToken.ClientIP = c.ClientIP()
	accessToken.ExpiresAt = _u.GetExpiryTime(72)

	idb.DB.Save(&accessToken)
	idb.DB.Model(&user).Update(map[string]interface{}{"last_login": time.Now()})

	idb.DB.Where(&models.AccessToken{UserID: user.ID}).First(&accessToken)

	c.JSON(http.StatusAccepted, _u.Message(true, "Refresh token success", &accessToken, nil))
}

// VerifyLogin for check pin and login
func (idb *InDB) VerifyLogin(c *gin.Context) {
	var (
		user        models.User
		accessToken models.AccessToken
	)

	verifyValidator := _request.VerifyValidator()
	if err := verifyValidator.BindVerify(c); err != nil {
		c.JSON(http.StatusBadRequest, _u.NewValidatorError(err))
		return
	}

	cv := &verifyValidator

	err := idb.DB.Where(models.User{Phone: cv.Data.Phone}).First(&user).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, _u.Message(false, "User not found", nil, nil))
		return
	}

	byteDBPin := []byte(user.Pin)
	byteReqPin := []byte(cv.Data.Pin)

	if errorPin := bcrypt.CompareHashAndPassword(byteDBPin, byteReqPin); errorPin != nil {
		c.JSON(http.StatusBadRequest, _u.Message(false, "Wrong Pin Number", nil, nil))
		return
	}

	if c.Query("only_check") == "true" {
		c.JSON(http.StatusOK, _u.Message(true, "Verified", &user, nil))
	} else {
		// split user agent
		splitted := strings.Split(c.Request.UserAgent(), " ")

		// Create the Claims
		claims := _u.CreateClaims("users", user.ID, splitted[0], 72)
		token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims) //Create Token
		signed, err := token.SignedString([]byte("secret"))

		if err != nil {
			c.JSON(http.StatusBadRequest, _u.Message(false, err.Error(), nil, nil))
			c.Abort()
			return
		}

		checkAccessToken := idb.DB.Where("user_id = ?", user.ID).First(&accessToken).RecordNotFound()

		accessToken.UserID = user.ID
		accessToken.AccessToken = signed
		accessToken.RefreshToken = _u.TokenGenerator(user.Phone)
		accessToken.TokenType = "bearer"
		accessToken.Scope = "mobile"
		accessToken.UserAgent = splitted[0]
		accessToken.ClientIP = c.ClientIP()
		accessToken.ExpiresAt = _u.GetExpiryTime(72)

		if checkAccessToken {
			idb.DB.Create(&accessToken)
		} else {
			idb.DB.Model(&accessToken).Update(&accessToken)
		}

		idb.DB.Model(&user).Update(map[string]interface{}{"last_login": time.Now()})

		idb.DB.Where(&models.User{ID: user.ID}).First(&user)
		idb.DB.Where(&models.AccessToken{UserID: user.ID}).First(&accessToken)

		// Check and get wallet user
		req := _http.NewHttpService()
		req.TransactionServ.SetToken(c.GetHeader("Authorization"))
		walletUser, err := req.TransactionServ.GetWallet(user.ID, "users")
		if err != nil {
			c.JSON(http.StatusBadRequest, _u.Message(false, "Failed to Get wallet user", nil, nil))
			return
		}

		if walletUser.Status == false {
			c.JSON(http.StatusBadRequest, _u.Message(false, walletUser.Message, nil, nil))
			return
		}

		if len(walletUser.Data) == 0 {
			c.JSON(http.StatusBadRequest, _u.Message(false, "Wallet user not found", nil, nil))
			return
		}

		type includes struct {
			*models.User
			EWallet interface{} `json:"e_wallet"`
		}

		c.JSON(http.StatusCreated, _u.Message(true, "Login success", &accessToken, includes{
			User:    &user,
			EWallet: walletUser.Data[0],
		}))
	}
}

// CreatePin for login with PIN after input number phone
func (idb *InDB) CreatePin(c *gin.Context) {
	var (
		accessToken models.AccessToken
		user        models.User
	)

	e := godotenv.Load()
	if e != nil {
		log.Print(e)
	}

	hashing := c.Query("hash")
	if hashing == "" {
		c.JSON(http.StatusUnprocessableEntity, _u.Message(false, "hash query not found", nil, nil))
		return
	}

	pinValidator := _request.PinValidator()
	if err := pinValidator.BindCreatePin(c); err != nil {
		fmt.Println("errr", err)
		c.JSON(http.StatusBadRequest, _u.NewValidatorError(err))
		return
	}

	cv := &pinValidator

	if idb.DB.Where(models.User{Hash: hashing}).First(&user).RecordNotFound() {
		c.JSON(http.StatusBadRequest, _u.Message(false, "Wrong hash matching", nil, nil))
	} else {
		hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(cv.Data.PinUser), bcrypt.DefaultCost)
		if errHash != nil {
			c.JSON(http.StatusBadRequest, _u.Message(false, errHash.Error(), nil, nil))
			return
		}
		user.Pin = string(hashedPassword)
		user.Status = 1
		if err := idb.DB.Save(&user).Error; err != nil {
			c.JSON(http.StatusBadRequest, _u.Message(false, err.Error(), nil, nil))
			return
		}

		// split user agent
		splitted := strings.Split(c.Request.UserAgent(), " ")

		// Create the Claims
		claims := _u.CreateClaims("users", user.ID, splitted[0], 72)
		token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims) //Create Token
		signed, err := token.SignedString([]byte("secret"))

		if err != nil {
			c.JSON(http.StatusBadRequest, _u.Message(false, err.Error(), nil, nil))
			return
		}

		accessToken.UserID = user.ID
		accessToken.AccessToken = signed
		accessToken.RefreshToken = _u.TokenGenerator(user.Phone)
		accessToken.TokenType = "bearer"
		accessToken.Scope = "mobile"
		accessToken.UserAgent = splitted[0]
		accessToken.ClientIP = c.ClientIP()
		accessToken.ExpiresAt = _u.GetExpiryTime(72)

		// Check and post wallet user
		req := _http.NewHttpService()
		req.TransactionServ.SetToken(c.GetHeader("Authorization"))
		walletUser, err := req.TransactionServ.CreateWallet(user.ID, "users", 0, "Post")
		if err != nil {
			c.JSON(http.StatusBadRequest, _u.Message(false, "Failed to Get wallet user", nil, nil))
			return
		}

		if walletUser.Status == false {
			c.JSON(http.StatusBadRequest, _u.Message(false, walletUser.Message, nil, nil))
			return
		}

		idb.DB.Create(&accessToken) //Create DB Access Token

		includes := map[string]interface{}{
			"access_token": &accessToken,
			"e_wallet":     &walletUser.Data[0],
		}
		c.JSON(http.StatusCreated, _u.Message(true, "Pin successfully created", &user, &includes))
	}
}

//ForgotPin for forgot pin
func (idb *InDB) ForgotPin(c *gin.Context) {
	var user models.User

	forgotPinValidator := _request.ForgotPinValidator()
	if err := forgotPinValidator.BindForgotPin(c); err != nil {
		c.JSON(http.StatusBadRequest, _u.NewValidatorError(err))
		return
	}

	cv := &forgotPinValidator

	accountKit, errExchange := _u.ExchangeAccountKit(c.Query("code_access_token"))
	if errExchange != nil {
		c.JSON(http.StatusBadRequest, _u.Message(true, errExchange.Error(), nil, nil))
		return
	}

	if idb.DB.Where("phone = ?", accountKit.Phone.NationalNumber).First(&user).RecordNotFound() {
		c.JSON(http.StatusBadRequest, _u.Message(true, "Phone not found", nil, nil))
		return
	}

	hashedPIN, errorNewPin := bcrypt.GenerateFromPassword([]byte(cv.Data.NewPin), bcrypt.DefaultCost)
	if errorNewPin != nil {
		c.JSON(http.StatusBadRequest, _u.Message(false, errorNewPin.Error(), nil, nil))
		return
	}

	user.Pin = string(hashedPIN)

	if err := idb.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusNotModified, _u.Message(false, err.Error(), nil, nil))
		return
	}

	c.JSON(http.StatusAccepted, _u.Message(true, "Change forgotten pin completed", &user, nil))
}
