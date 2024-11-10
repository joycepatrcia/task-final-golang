package handler

import (
	"fmt"
	"godb/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AccountInterface interface {
	Create(*gin.Context)
	Read(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
	List(*gin.Context)
	TopUp(*gin.Context)
	Balance(*gin.Context)
	My(*gin.Context)
	Transfer(*gin.Context)
	MutationList(*gin.Context)
}

type accountImplement struct {
	db *gorm.DB
}

func NewAccount(db *gorm.DB) AccountInterface {
	return &accountImplement{
		db: db,
	}
}

func (a *accountImplement) Create(c *gin.Context) {
	payload := model.Account{}

	err := c.BindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	// Create data
	result := a.db.Create(&payload)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Create success",
		"data":    payload,
	})
}

func (a *accountImplement) Read(c *gin.Context) {
	var account model.Account

	id := c.Param("id")

	if err := a.db.First(&account, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Not found",
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": account,
	})
}

func (a *accountImplement) Update(c *gin.Context) {
	payload := model.Account{}

	// Bind the incoming JSON to the payload
	err := c.BindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	id := c.Param("id")

	// Retrieve the existing account by ID
	account := model.Account{}
	result := a.db.First(&account, "account_id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Account not found",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	// Update fields only if they are present in the payload
	if payload.Name != "" {
		account.Name = payload.Name
	}

	// Check if Balance is provided and is a valid number
	if payload.Balance != 0 {
		account.Balance = payload.Balance
	}

	// Save the updated account information in the database
	if err := a.db.Save(&account).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update account",
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Update successful",
	})
}


func (a *accountImplement) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := a.db.Where("account_id = ?", id).Delete(&model.Account{}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Not found",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Delete success",
		"data": map[string]string{
			"account_id": id,
		},
	})
}

func (a *accountImplement) List(c *gin.Context) {
	var accounts []model.Account

	// Find and get all accounts data
	if err := a.db.Find(&accounts).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": accounts,
	})
}

func (h *accountImplement) TopUp(c *gin.Context) {
	var payload struct {
		AccountID int64 `json:"account_id"`
		Amount    int   `json:"amount"`
	}

	// Bind request body ke struct
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	if payload.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be greater than 0"})
		return
	}

	var account model.Account
	result := h.db.First(&account, "account_id = ?", payload.AccountID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve account", "details": result.Error.Error()})
		return
	}

	fmt.Printf("Account balance before top-up: %v\n", account.Balance)
	updateResult := h.db.Model(&model.Account{}).Where("account_id = ?", payload.AccountID).
		Update("balance", gorm.Expr("balance + ?", payload.Amount))

	if updateResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to top up balance", "details": updateResult.Error.Error()})
		return
	}

	fmt.Printf("Account balance after top-up: %v\n", account.Balance+int64(payload.Amount))
	c.JSON(http.StatusOK, gin.H{"message": "Balance topped up successfully"})
}


func (h *accountImplement) Balance(c *gin.Context) {
	accountID := c.GetInt64("account_id") 

	var account model.Account
	err := h.db.Select("balance").Where("account_id = ?", accountID).First(&account).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve balance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": account.Balance})
}

func (h *accountImplement) My(c *gin.Context) {
	accountID := c.GetInt64("account_id") 

	var account model.Account
	err := h.db.First(&account, "account_id = ?", accountID).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve account information"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": account})
}

func (h *accountImplement) Transfer(c *gin.Context) {
    var payload struct {
        ToAccountID     int64  `json:"to_account_id"`
        Amount              int    `json:"amount"`
        TransactionCategoryID *int64 `json:"transaction_category_id"`
    }

	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	if payload.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be greater than 0"})
		return
	}

	currentAccountID := c.GetInt64("account_id")
	var currentAccount model.Account
	err = h.db.First(&currentAccount, "account_id = ?", currentAccountID).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve current account"})
		return
	}

	if currentAccount.Balance < int64(payload.Amount) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
		return
	}

	var targetAccount model.Account
	err = h.db.First(&targetAccount, "account_id = ?", payload.ToAccountID).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Target account not found"})
		return
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&currentAccount).Update("balance", gorm.Expr("balance - ?", payload.Amount)).Error //minus from
		if err != nil {
			return err
		}

		// Add target account
		err = tx.Model(&targetAccount).Update("balance", gorm.Expr("balance + ?", payload.Amount)).Error
		if err != nil {
			return err
		}

		// masukin ke table transaction
		transaction := model.Transaction{
			AccountID: 			  &currentAccountID,		
			FromAccountID:        &currentAccountID,
			ToAccountID:          &payload.ToAccountID,
			TransactionCategoryID: payload.TransactionCategoryID,  
			Amount:               int64(payload.Amount),
			TransactionDate:      time.Now(),
		}
		err = tx.Create(&transaction).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete transfer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transfer successful"})
}

func (h *accountImplement) MutationList(c *gin.Context) {
	accountID := c.GetInt64("account_id") 

	startDate := c.DefaultQuery("start_date", "")
	endDate := c.DefaultQuery("end_date", "")

	var mutations []model.Transaction

	query := h.db.Where("from_account_id = ? OR to_account_id = ?", accountID, accountID)

	if startDate != "" {
		startTime, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format"})
			return
		}
		query = query.Where("transaction_date >= ?", startTime)
	}

	if endDate != "" {
		endTime, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format"})
			return
		}
		query = query.Where("transaction_date <= ?", endTime)
	}

	if err := query.Order("transaction_date desc").Limit(10).Find(&mutations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve mutations"})
		return
	}

	c.JSON(http.StatusOK, mutations)
}

