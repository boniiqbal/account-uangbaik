package service

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	_request "uangbaik-account-microservice/admin/http/request"
	"uangbaik-account-microservice/database/models"
	_http "uangbaik-account-microservice/http"
	_u "uangbaik-account-microservice/utils"

	"github.com/google/uuid"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// LoginHandler for login dashboard
func (idb *InDB) LoginHandler(c *gin.Context) {
	var (
		admin           models.Admin
		adminAcessToken models.AdminAccessToken
	)

	loginValidator := _request.LoginValidator()
	if err := loginValidator.BindLoginAdmin(c); err != nil {
		c.JSON(http.StatusBadRequest, _u.NewValidatorError(err))
		return
	}

	cv := &loginValidator

	if idb.DB.Where(models.Admin{Username: cv.Data.Username}).First(&admin).RecordNotFound() {
		c.JSON(http.StatusBadRequest, _u.Message(false, "Wrong Username", nil, nil))
		return
	}

	byteDBPass := []byte(admin.Password)
	byteReqPass := []byte(cv.Data.Password)

	if error := bcrypt.CompareHashAndPassword(byteDBPass, byteReqPass); error != nil {
		c.JSON(http.StatusBadRequest, _u.Message(false, "Wrong Password", error.Error(), nil))
		return
	}

	type MyCustomClaims struct {
		UserID uint
		jwt.StandardClaims
	}

	// split user agent
	splitted := strings.Split(c.Request.UserAgent(), " ")

	// Create the Claims
	claims := _u.CreateClaims("admins", admin.ID, splitted[0], 6)
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims) //Create Token
	signed, err := token.SignedString([]byte("secret"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, _u.Message(false, "Problem create token", err.Error(), nil))
		return
	}

	checkAccessToken := idb.DB.Where(models.AdminAccessToken{AdminID: admin.ID}).First(&adminAcessToken).RecordNotFound()

	adminAcessToken.AdminID = admin.ID
	adminAcessToken.AccessToken = signed
	adminAcessToken.RefreshToken = uuid.New().String()
	adminAcessToken.TokenType = "bearer"
	adminAcessToken.Scope = "admin"
	adminAcessToken.UserAgent = splitted[0]
	adminAcessToken.ClientIP = c.ClientIP()
	adminAcessToken.ExpiresAt = _u.GetExpiryTime(6)

	if checkAccessToken {
		idb.DB.Create(&adminAcessToken)
	} else {
		idb.DB.Save(&adminAcessToken)
	}

	idb.DB.Model(&admin).Select("last_login").Update(map[string]interface{}{"last_login": time.Now()})

	idb.DB.Where(&models.Admin{ID: admin.ID}).First(&admin)
	idb.DB.Where(&models.AdminAccessToken{AdminID: admin.ID}).First(&adminAcessToken)

	c.Writer.Header().Set("Content-Type", "application/json")
	c.JSON(http.StatusCreated, _u.Message(true, "Login success", &adminAcessToken, &admin))
}

// GetUserDetail for detail user
func (idb *InDB) GetUserDetail(c *gin.Context) {
	var (
		user models.User
	)

	code := c.Query("uuid")
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusUnprocessableEntity, _u.Message(false, "User Id is undefined", nil, nil))
		return
	}

	userIDUint, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, _u.Message(false, err.Error(), nil, nil))
		return
	}

	if code == "" {
		if err := idb.DB.Preload("BankAccount").Preload("UploadAsset").Where(models.User{ID: uint(userIDUint)}).First(&user).Error; err != nil {
			c.JSON(http.StatusBadRequest, _u.Message(false, "User Not Found", nil, nil))
			return
		}
	} else {
		actorUUIDParsed, errParseUUID := uuid.Parse(code)
		if errParseUUID != nil {
			c.JSON(http.StatusUnprocessableEntity, _u.Message(false, "error actor code undefined", errParseUUID.Error(), nil))
			return
		}
		if err := idb.DB.Preload("BankAccount").Preload("UploadAsset").Where(models.User{UUID: actorUUIDParsed}).First(&user).Error; err != nil {
			c.JSON(http.StatusBadRequest, _u.Message(false, err.Error(), nil, nil))
			return
		}
	}

	// Get Wallet User
	req := _http.NewHttpService()
	req.TransactionServ.SetToken(c.GetHeader("Authorization"))
	walletUser, err := req.TransactionServ.GetWallet(user.ID, "users")
	if err != nil {
		c.JSON(http.StatusBadRequest, _u.Message(false, "Failed to get wallet", nil, nil))
		return
	}

	userWallet := walletUser.Data[0]

	if walletUser.Status == false {
		c.JSON(http.StatusBadRequest, _u.Message(false, walletUser.Message, nil, nil))
		return
	}

	type data struct {
		*models.User
		EWallet interface{} `json:"e_wallet"`
	}

	c.JSON(http.StatusAccepted, _u.Message(true, "User fetched", data{
		User:    &user,
		EWallet: userWallet,
	}, nil))
}

// GetUserList for all user
func (idb *InDB) GetUserList(c *gin.Context) {
	type Flow struct {
		Count int
	}

	type Filtered struct {
		Count int `json:"count"`
	}

	type Total struct {
		Count int `json:"count"`
	}

	type Pagination struct {
		PrevOffset int `json:"prev_offset"`
		NextOffset int `json:"next_offset"`
		Offset     int `json:"offset"`
		Limit      int `json:"limit"`
	}

	type response struct {
		Pagination      Pagination  `json:"pagination"`
		Total           Total       `json:"total"`
		Filtered        Filtered    `json:"filtered"`
		TotalCount      int         `json:"total_count"`
		TotalUnverified int         `json:"total_unverified"`
		Attributes      interface{} `json:"attributes"`
	}

	type DataResponse struct {
		*models.User
		EWallet interface{} `json:"e_wallet"`
	}

	var (
		user            []models.User
		unverifiedCount int
		offset          int
		limit           int
		flow            Flow
		filteredFlow    Flow
		dataArrayUser   []interface{}
	)

	search := c.Query("search")
	userID := c.Query("user_id")
	status := c.Query("status")

	offset, _ = strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ = strconv.Atoi(c.DefaultQuery("limit", "10"))
	sort := c.DefaultQuery("sort", "-updated_at")
	sortString := _u.SortQuery(sort)

	idb.DB.Table("users").Select("count(*) as count").Scan(&flow)

	tx := idb.DB.Model(&user)
	if search == "" {
		tx = tx.Where("id > ?", 0)
	}

	if userID != "" {
		tx = tx.Where("id = ?", userID)
	}

	if status != "" {
		tx.Where("status = ?", status)
	}

	if search != "" {
		tx = tx.Where("full_name LIKE ?", "%"+search+"%")
	}

	tx.Select("count(*) as count").Scan(&filteredFlow)

	tx.Offset(offset).Limit(limit).Preload("BankAccount").Order(sortString).Find(&user)

	idb.DB.Table("users").Where("status < 2").Count(&unverifiedCount)

	var prevOffset int
	if offset < limit {
		prevOffset = 0
	} else {
		prevOffset = offset - limit
	}

	var nextOffset int
	if offset+limit < flow.Count {
		nextOffset = offset + limit
	} else {
		nextOffset = offset
	}

	// Check and get Wallet user
	req := _http.NewHttpService()
	req.TransactionServ.SetToken(c.GetHeader("Authorization"))
	walletUser, err := req.TransactionServ.GetWallet(0, "users") 
	if err != nil {
		c.JSON(http.StatusBadRequest, _u.Message(false, "Failed to get wallet", nil, nil))
		return
	}

	if len(walletUser.Data) == 0 {
		c.JSON(http.StatusBadRequest, _u.Message(false, "Wallet user not found", nil, nil))
		return
	}

	dataResponse := &DataResponse{}
	for i := 0; i < len(user); i++ {
		for j := 0; j < len(walletUser.Data); j++ {
			if walletUser.Data[j].ActorID == user[i].ID {
				dataResponse = &DataResponse{
					User:    &user[i],
					EWallet: walletUser.Data[j],
				}
				dataArrayUser = append(dataArrayUser, dataResponse)
			}
			continue
		}
	}

	data := response{
		Pagination: Pagination{
			PrevOffset: prevOffset,
			NextOffset: nextOffset,
			Offset:     offset,
			Limit:      limit,
		},
		Filtered: Filtered{
			Count: filteredFlow.Count,
		},
		Total: Total{
			Count: flow.Count,
		},
		TotalCount:      flow.Count,
		TotalUnverified: unverifiedCount,
		Attributes:      dataArrayUser,
	}

	c.JSON(http.StatusOK, _u.Message(true, "Your user are successfully retrieved", &data, nil))
}

// EditUserData for all user
func (idb *InDB) EditUserData(c *gin.Context) {
	var (
		user     models.User
		notifMsg string
	)

	if c.Param("id") == "" {
		c.JSON(http.StatusUnprocessableEntity, _u.Message(false, "param ID can't be empty", nil, nil))
		return
	}

	editRequest := _request.EditUserValidator()
	if err := editRequest.BindEditUser(c); err != nil {
		c.JSON(http.StatusBadRequest, _u.NewValidatorError(err))
		return
	}

	if err := idb.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, _u.Message(false, err.Error(), nil, nil))
		return
	}

	if err := idb.DB.Model(&user).Where("id = ?", c.Param("id")).Update(&editRequest.Data).Error; err != nil {
		c.JSON(http.StatusBadRequest, _u.Message(false, err.Error(), &editRequest, nil))
		return
	}

	strID := strconv.FormatUint(uint64(user.ID), 10)

	if editRequest.Data.Status == 2 {
		notifMsg = fmt.Sprintf("Akun %s berhasil diaktivasi", user.Phone)
	} else if editRequest.Data.Status == -1 {
		notifMsg = fmt.Sprintf("Akun %s telah ditolak", user.Phone)
	}

	notif := _u.NewNotification(strID, "Ya Baik", notifMsg, nil)
	notif.Send()

	req := _http.NewHttpService()
	req.NotificationServ.SetToken(c.GetHeader("Authotization"))
	_, err := req.NotificationServ.PostNotification(strID, notifMsg, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, _u.Message(false, "Failed to push Notification", err.Error(), nil))
		return
	}

	idb.DB.Where("id = ?", c.Param("id")).First(&user)

	c.JSON(http.StatusOK, _u.Message(true, "Your user are successfully updated", &user, nil))
}

// GetUploadAsset .
func (idb *InDB) GetUploadAsset(c *gin.Context) {
	var upload models.UploadAsset

	UploadID := c.Query("upload_id")
	uploadType := c.Query("upload_type")
	Status := c.Query("status")
	nType := c.Query("type")

	uploadID, _ := strconv.ParseUint(UploadID, 10, 32)
	status, _ := strconv.ParseInt(Status, 10, 64)

	if err := idb.DB.Where(models.UploadAsset{UploadID: uint(uploadID), UploadType: uploadType, Type: nType, Status: int(status)}).Last(&upload).Error; err != nil {
		c.JSON(http.StatusBadRequest, _u.Message(false, "Image Not Found", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, _u.Message(true, "Image succesfully retrieved", &upload, nil))
}
