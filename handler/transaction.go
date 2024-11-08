package handler

import (
	"godb/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Define the interface for the transaction handler
type TransactionHandlerInterface interface {
	NewTransaction(*gin.Context)
	TransactionList(*gin.Context)
}

type transactionHandler struct {
	db *gorm.DB
}

// Constructor function to create a new transaction handler
func NewTransactionHandler(db *gorm.DB) TransactionHandlerInterface {
	return &transactionHandler{
		db: db,
	}
}

// NewTransaction - Create a new transaction record
func (t *transactionHandler) NewTransaction(c *gin.Context) {
	var payload model.Transaction

	// Bind JSON to payload
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Set the transaction date to the current time
	payload.TransactionDate = time.Now()

	// Insert the transaction record into the database
	if err := t.db.Create(&payload).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction created successfully",
		"data":    payload,
	})
}

// TransactionList - Get the latest transactions for a given account ID
func (t *transactionHandler) TransactionList(c *gin.Context) {
	accountID := c.Param("account_id")

	// Prepare a slice to hold the results
	var transactions []model.Transaction

	// Query the database to get the latest transactions for the given account ID
	if err := t.db.Where("account_id = ?", accountID).Order("transaction_date DESC").Limit(5).Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"data": transactions,
	})
}
