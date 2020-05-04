package service

import (
	"encoding/json"
	"uangbaik-account-microservice/database/models"
	_http "uangbaik-account-microservice/http"
	_u "uangbaik-account-microservice/utils"

	"net/http"

	"github.com/gin-gonic/gin"
)

type UserInfo struct {
	UserID   uint   `json:"user_id"`
	FullName string `json:"full_name"`
}

//Summary for detail organization
func (idb *InDB) Summary(c *gin.Context) {
	var (
		admin                    models.Admin
		users                    models.User
		userInfo                 UserInfo
		countUser                int
		countUnverifUser         int
		countMerchant            int
		countUnverifMerchant     int
		countWithdraw            int
		countRejectWithdraw      int
		countPendingWithdraw     int
		countApproveWithdraw     int
		countDisburseWithdraw    int
		countCampaign            int
		countSuspendCampaign     int
		countOrganization        int
		countSuspendOrganization int
	)

	userData := c.GetHeader("User-Info")

	if userData == "" {
		c.JSON(http.StatusBadRequest, _u.Message(false, "User not found", nil, nil))
		return
	}

	if userData == "" {
		c.JSON(http.StatusBadRequest, _u.Message(false, "User Not Found", nil, nil))
		return
	}

	err := json.Unmarshal([]byte(userData), &userInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, _u.Message(false, err.Error(), nil, nil))
		return
	}

	type Users struct {
		Total      int `json:"total"`
		Unverified int `json:"unverified"`
	}

	type Merchants struct {
		Total      int `json:"total"`
		Unverified int `json:"unverified"`
	}

	type Withdraws struct {
		Total     int `json:"total"`
		Rejected  int `json:"rejected"`
		Pending   int `json:"pending"`
		Approved  int `json:"approved"`
		Disbursed int `json:"disbursed"`
	}

	type Organizations struct {
		Total     int `json:"total"`
		Suspended int `json:"suspended"`
	}

	type Campaigns struct {
		Total     int `json:"total"`
		Suspended int `json:"suspended"`
	}

	type Attributes struct {
		Admin         *models.Admin `json:"admin_info"`
		Users         Users         `json:"users"`
		Merchants     Merchants     `json:"merchants"`
		Withdraws     Withdraws     `json:"withdraws"`
		Organizations Organizations `json:"organizations"`
		Campaigns     Campaigns     `json:"campaigns"`
	}

	type response struct {
		Attributes Attributes `json:"attributes"`
	}

	req := _http.NewHttpService()

	idb.DB.Model(&admin).Where("id = ?", userInfo.UserID).First(&admin)

	idb.DB.Model(&users).Count(&countUser)
	idb.DB.Model(&users).Where("status < 2").Count(&countUnverifUser)

	// Check and Get Merchant Data
	req.MerchantServ.SetToken(c.GetHeader("Authorization"))
	merchant, err1 := req.MerchantServ.GetMerchantDetail(0, 0, "", 0)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, _u.Message(false, "Failed to get merchant data", nil, nil))
		return
	}

	//Check and get Withdraw Data
	req.TransactionServ.SetToken(c.GetHeader("Authorization"))
	withdraw, err2 := req.TransactionServ.GetWithdraw()
	if err2 != nil {
		c.JSON(http.StatusBadRequest, _u.Message(false, "Failed to get withdraw", nil, nil))
		return
	}

	// Check and get campaign
	req.ZiswafServ.SetToken(c.GetHeader("Authorization"))
	campaign, err := req.ZiswafServ.GetCampaign()
	if err != nil {
		c.JSON(http.StatusBadRequest, _u.Message(false, "Failed to get campaign", nil, nil))
		return
	}

	// Check and get Organization
	req.ZiswafServ.SetToken(c.GetHeader("Authorization"))
	organization, err := req.ZiswafServ.GetOrganization()
	if err != nil {
		c.JSON(http.StatusBadRequest, _u.Message(false, "Failed to get organization", nil, nil))
		return
	}

	countMerchant = merchant.Data.Total.Count
	countUnverifMerchant = merchant.Data.Total.UnverifiedCount

	countWithdraw = withdraw.Data.Withdraw.Total
	countDisburseWithdraw = withdraw.Data.Withdraw.Disbursed
	countApproveWithdraw = withdraw.Data.Withdraw.Approved
	countPendingWithdraw = withdraw.Data.Withdraw.Pending
	countRejectWithdraw = withdraw.Data.Withdraw.Rejected

	countCampaign = campaign.Data.TotalCount
	countSuspendCampaign = campaign.Data.TotalSuspended

	countOrganization = organization.Data.TotalCount
	countSuspendOrganization = organization.Data.TotalSuspended

	data := response{
		Attributes: Attributes{
			&admin,
			Users{
				Total:      countUser,
				Unverified: countUnverifUser,
			},
			Merchants{
				Total:      countMerchant,
				Unverified: countUnverifMerchant,
			},
			Withdraws{
				Total:     countWithdraw,
				Pending:   countPendingWithdraw,
				Approved:  countApproveWithdraw,
				Disbursed: countDisburseWithdraw,
				Rejected:  countRejectWithdraw,
			},
			Organizations{
				Total:     countOrganization,
				Suspended: countSuspendOrganization,
			},
			Campaigns{
				Total:     countCampaign,
				Suspended: countSuspendCampaign,
			},
		},
	}

	c.JSON(http.StatusOK, _u.Message(true, "Summary fetched", &data, nil))
}
