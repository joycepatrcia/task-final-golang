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

	// bind JSON Request to payload
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

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Create success",
		"data":    payload,
	})
}

func (a *accountImplement) Read(c *gin.Context) {
	var account model.Account

	// get id from url account/read/5, 5 will be the id
	id := c.Param("id")

	// Find first data based on id and put to account model
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

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"data": account,
	})
}

func (a *accountImplement) Update(c *gin.Context) {
	payload := model.Account{}

	// bind JSON Request to payload
	err := c.BindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	// get id from url account/update/5, 5 will be the id
	id := c.Param("id")

	// Find first data based on id and put to account model
	account := model.Account{}
	result := a.db.First(&account, "account_id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Not found",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	// Update data
	account.Name = payload.Name
	a.db.Save(account)

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Update success",
	})
}

func (a *accountImplement) Delete(c *gin.Context) {
	// get id from url account/delete/5, 5 will be the id
	id := c.Param("id")

	// Find first data based on id and delete it
	if err := a.db.Where("account_id = ?", id).Delete(&model.Account{}).Error; err != nil {
		// No data found and deleted
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

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Delete success",
		"data": map[string]string{
			"account_id": id,
		},
	})
}

func (a *accountImplement) List(c *gin.Context) {
	// Prepare empty result
	var accounts []model.Account

	// Find and get all accounts data and put to &accounts
	if err := a.db.Find(&accounts).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"data": accounts,
	})
}

// Handler for "POST /account/topup"
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

	// Validasi amount
	if payload.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be greater than 0"})
		return
	}

	// Cek apakah account_id ada di database (gunakan nama kolom yang benar, misalnya "account_id" bukan "id")
	var account model.Account
	result := h.db.First(&account, "account_id = ?", payload.AccountID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
			return
		}
		// Log error lainnya
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve account", "details": result.Error.Error()})
		return
	}

	// Debugging log untuk melihat saldo sebelum top-up
	fmt.Printf("Account balance before top-up: %v\n", account.Balance)

	// Update saldo akun dengan menambah jumlah top-up
	updateResult := h.db.Model(&model.Account{}).Where("account_id = ?", payload.AccountID).
		Update("balance", gorm.Expr("balance + ?", payload.Amount))

	// Cek apakah update berhasil
	if updateResult.Error != nil {
		// Menampilkan error dari GORM
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to top up balance", "details": updateResult.Error.Error()})
		return
	}

	// Menampilkan saldo setelah top-up
	fmt.Printf("Account balance after top-up: %v\n", account.Balance+int64(payload.Amount))

	c.JSON(http.StatusOK, gin.H{"message": "Balance topped up successfully"})
}

// Handler for "GET /account/balance"
func (h *accountImplement) Balance(c *gin.Context) {
	accountID := c.GetInt64("account_id") // Assume account_id comes from middleware

	var account model.Account
	err := h.db.Select("balance").Where("account_id = ?", accountID).First(&account).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve balance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": account.Balance})
}

// Handler for "GET /account/my"
func (h *accountImplement) My(c *gin.Context) {
	accountID := c.GetInt64("account_id") // Assume account_id comes from middleware

	var account model.Account
	err := h.db.First(&account, "account_id = ?", accountID).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve account information"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": account})
}

// Handler for "POST /account/transfer"
func (h *accountImplement) Transfer(c *gin.Context) {
    var payload struct {
        ToAccountID     int64  `json:"to_account_id"`
        Amount              int    `json:"amount"`
        TransactionCategoryID *int64 `json:"transaction_category_id"`
    }

	// Bind request body to struct
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Validate amount
	if payload.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be greater than 0"})
		return
	}

	// Get current account info
	currentAccountID := c.GetInt64("account_id")
	var currentAccount model.Account
	err = h.db.First(&currentAccount, "account_id = ?", currentAccountID).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve current account"})
		return
	}

	// Check if the balance is sufficient
	if currentAccount.Balance < int64(payload.Amount) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
		return
	}

	// Retrieve the target account
	var targetAccount model.Account
	err = h.db.First(&targetAccount, "account_id = ?", payload.ToAccountID).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Target account not found"})
		return
	}

	// Deduct balance from current account and add balance to target account
	err = h.db.Transaction(func(tx *gorm.DB) error {
		// Deduct from current account
		err := tx.Model(&currentAccount).Update("balance", gorm.Expr("balance - ?", payload.Amount)).Error
		if err != nil {
			return err
		}

		// Add to target account
		err = tx.Model(&targetAccount).Update("balance", gorm.Expr("balance + ?", payload.Amount)).Error
		if err != nil {
			return err
		}

		// Create a transaction record (use pointers to account IDs)
		transaction := model.Transaction{
			FromAccountID:        &currentAccountID,
			ToAccountID:          &payload.ToAccountID,
			TransactionCategoryID: payload.TransactionCategoryID,  // Set category
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

// Handler for "GET /account/mutation"
func (h *accountImplement) MutationList(c *gin.Context) {
	accountID := c.GetInt64("account_id") // Assume account_id comes from middleware

	startDate := c.DefaultQuery("start_date", "")
	endDate := c.DefaultQuery("end_date", "")

	var mutations []model.Transaction

	query := h.db.Where("from_account_id = ? OR to_account_id = ?", accountID, accountID)

	// Add date filters if provided
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
	
	if err := query.Order("transaction_date desc").Find(&mutations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve mutations"})
		return
	}
	
	

	c.JSON(http.StatusOK, gin.H{"data": mutations})
}
