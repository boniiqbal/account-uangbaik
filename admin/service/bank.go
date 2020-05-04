package service

import (
	"strconv"
	"uangbaik-account-microservice/database/models"
	_u "uangbaik-account-microservice/utils"

	"net/http"

	"github.com/gin-gonic/gin"
)

//GetListBank
func (idb *InDB) GetListBank(c *gin.Context) {
	var bank []models.Bank

	if err := idb.DB.Find(&bank).Error; err != nil {
		c.JSON(http.StatusBadRequest, _u.Message(false, "Bank not found", nil, nil))
		return
	}

	c.JSON(http.StatusOK, _u.Message(true, "Bank data sucessfully retrived", &bank, nil))
}

//CreateBank for create Bank
func (idb *InDB) CreateBank(c *gin.Context) {
	var (
		bank models.Bank
	)

	request := struct {
		Data struct {
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"data"`
	}{}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusUnprocessableEntity, _u.Message(false, "Error while decoding request body", err.Error(), nil))
		return
	}

	bank.Code = request.Data.Code
	bank.Name = request.Data.Name

	if err := idb.DB.Create(&bank).Error; err != nil {
		c.JSON(http.StatusBadRequest, _u.Message(false, err.Error(), nil, nil))
		return
	}

	c.JSON(http.StatusOK, _u.Message(true, "Bank Successfully created", &bank, nil))
}

// UpdateBank .
func (idb *InDB) UpdateBank(c *gin.Context) {
	var bank models.Bank

	if c.Param("id") == "" {
		c.JSON(http.StatusUnprocessableEntity, _u.Message(false, "id must be fill", nil, nil))
		return
	}

	bankID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, _u.Message(false, err.Error(), nil, nil))
		return
	}

	request := struct {
		Data struct {
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"data"`
	}{}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusUnprocessableEntity, _u.Message(false, "Error while decoding request body", err.Error(), nil))
		return
	}

	if err := idb.DB.Where(&models.Bank{ID: uint(bankID)}).First(&bank).Error; err != nil {
		c.JSON(http.StatusBadRequest, _u.Message(false, "Bank not found", nil, nil))
		return
	}

	bank.Code = request.Data.Code
	bank.Name = request.Data.Name

	if err := idb.DB.Save(&bank).Error; err != nil {
		c.JSON(http.StatusNotModified, _u.Message(false, err.Error(), nil, nil))
		return
	}
	c.JSON(http.StatusOK, _u.Message(true, "Bank has been updated", &bank, nil))
}

// DeleteBank for delete bank
func (idb *InDB) DeleteBank(c *gin.Context) {
	var (
		bank models.Bank
	)

	if c.Param("id") == "" {
		c.JSON(http.StatusUnprocessableEntity, _u.Message(false, "Param id required", nil, nil))
		return
	}

	bankID, errContext := strconv.ParseUint(c.Param("id"), 10, 32)
	if errContext != nil {
		c.JSON(http.StatusUnprocessableEntity, _u.Message(false, errContext.Error(), nil, nil))
		return
	}

	if err := idb.DB.Where("id = ?", uint(bankID)).Delete(&bank).Error; err != nil {
		c.JSON(http.StatusBadRequest, _u.Message(false, err.Error(), nil, nil))
		return
	}

	c.JSON(http.StatusOK, _u.Message(true, "Deleted Successfully", nil, nil))
}

// GetDetailBank for
func (idb *InDB) GetDetailBank(c *gin.Context) {
	var (
		bank models.Bank
	)

	if c.Param("id") == "" {
		c.JSON(http.StatusBadRequest, _u.Message(false, "Param ID is required", nil, nil))
		return
	}

	bankID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := idb.DB.Where(models.Bank{ID: uint(bankID)}).First(&bank).Error; err != nil {
		c.JSON(http.StatusBadRequest, _u.Message(false, "Bank Not Found", nil, nil))
	}

	c.JSON(http.StatusAccepted, _u.Message(true, "User fetched", &bank, nil))
}
