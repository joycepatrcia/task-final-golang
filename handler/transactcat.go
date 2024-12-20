package handler

import (
	"godb/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TransactionCategoriesInterface interface {
	Create(*gin.Context)
	Read(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
	List(*gin.Context)

	My(*gin.Context)
}

type transactionCatImplement struct {
	db *gorm.DB
}

func NewTransactionCategories(db *gorm.DB) TransactionCategoriesInterface {
	return &transactionCatImplement{
		db: db,
	}
}

func (a *transactionCatImplement) Create(c *gin.Context) {
	payload := model.TransactionCategories{}

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

func (a *transactionCatImplement) Read(c *gin.Context) {
	var transcat model.TransactionCategories

	// get id from url account/read/5, 5 will be the id
	id := c.Param("id")

	// Find first data based on id and put to account model
	if err := a.db.First(&transcat, id).Error; err != nil {
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
		"data": transcat,
	})
}

func (a *transactionCatImplement) Update(c *gin.Context) {
    payload := model.TransactionCategories{}

    // bind JSON Request to payload
    err := c.BindJSON(&payload)
    if err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    id := c.Param("id")

    transactcat := model.TransactionCategories{}
    result := a.db.First(&transactcat, "transaction_category_id = ?", id)
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
    transactcat.Name = payload.Name
    updateResult := a.db.Save(&transactcat) 
    if updateResult.Error != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
            "error": updateResult.Error.Error(),
        })
        return
    }

    // Success response
    c.JSON(http.StatusOK, gin.H{
        "message": "Update success",
        "data":    transactcat, 
    })
}


func (a *transactionCatImplement) Delete(c *gin.Context) {
	id := c.Param("id")

	// Delete the data based on id
	if err := a.db.Where("transaction_category_id = ?", id).Delete(&model.TransactionCategories{}).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Delete success",
		"data": map[string]string{
			"transaction_category_id": id,
		},
	})
}

func (a *transactionCatImplement) List(c *gin.Context) {
	// Prepare empty result
	var transactcats []model.TransactionCategories

	// Find and get all accounts data and put to &accounts
	if err := a.db.Find(&transactcats).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"data": transactcats,
	})
}

func (a *transactionCatImplement) My(c *gin.Context) {
	var transactcat model.TransactionCategories
	transactcatID := c.GetInt64("transaction_category_id")

	// Find first data based on transaction_category_id given
	if err := a.db.First(&transactcat, transactcatID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Not found",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"data": transactcat,
	})
}
