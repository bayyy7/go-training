package handlers

import (
	"example/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AccountInterface interface {
	Create(*gin.Context)
	Read(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
	List(*gin.Context)
	My(*gin.Context)
	TopUp(*gin.Context)
	Balance(*gin.Context)
	Transfer(*gin.Context)
}

type accountImplement struct {
	db *gorm.DB
}

func NewAccount(db *gorm.DB) AccountInterface {
	return &accountImplement{
		db: db,
	}
}

type transferPayload struct {
	TargetID int64 `json:"target_account_id" binding:"required"`
	Amount   int64 `json:"balance" binding:"required"`
}

func (a *accountImplement) Create(ctx *gin.Context) {
	payload := model.Account{}

	// bind JSON Request to payload
	err := ctx.BindJSON(&payload)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	// Create data
	result := a.db.Create(&payload)
	if result.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	// Success response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Create success",
		"data":    payload,
	})
}

func (a *accountImplement) Read(ctx *gin.Context) {
	var account model.Account

	// get id from url account/read/5, 5 will be the id
	id := ctx.Param("id")

	// Find first data based on id and put to account model
	if err := a.db.First(&account, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Not found",
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Success response
	ctx.JSON(http.StatusOK, gin.H{
		"data": account,
	})
}

func (a *accountImplement) Update(ctx *gin.Context) {
	payload := model.Account{}

	// bind JSON Request to payload
	err := ctx.BindJSON(&payload)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	// get id from url account/update/5, 5 will be the id
	id := ctx.Param("id")

	// Find first data based on id and put to account model
	account := model.Account{}
	result := a.db.First(&account, "account_id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Not found",
			})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	// Update data
	account.Name = payload.Name
	a.db.Save(account)

	// Success response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Update success",
	})
}

func (a *accountImplement) Delete(ctx *gin.Context) {
	// get id from url account/delete/5, 5 will be the id
	id := ctx.Param("id")

	// Find first data based on id and delete it
	if err := a.db.Where("account_id = ?", id).Delete(&model.Account{}).Error; err != nil {
		// No data found and deleted
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Not found",
			})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Success response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Delete success",
		"data": map[string]string{
			"account_id": id,
		},
	})
}

func (a *accountImplement) List(ctx *gin.Context) {
	// Prepare empty result
	var accounts []model.Account

	// Find and get all accounts data and put to &accounts
	if err := a.db.Find(&accounts).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Success response
	ctx.JSON(http.StatusOK, gin.H{
		"data": accounts,
	})
}

func (a *accountImplement) My(ctx *gin.Context) {
	var account model.Account
	// get account_id from middleware auth
	accountID := ctx.GetInt64("account_id")

	// Find first data based on account_id given
	if err := a.db.First(&account, accountID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Not found",
			})
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Success response
	ctx.JSON(http.StatusOK, gin.H{
		"data": account,
	})
}

func (a *accountImplement) TopUp(ctx *gin.Context) {
	payload := model.Account{}
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})

	}

	account := model.Account{}
	id := ctx.Param("id")
	if result := a.db.First(&account, "account_id=?", id); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": result.Error.Error(),
			})
			return
		}
		return
	}

	account.Balance = account.Balance + payload.Balance
	a.db.Save(account)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Update success",
	})
}

func (a *accountImplement) Balance(ctx *gin.Context) {
	var account model.Account
	accountID := ctx.GetInt64("account_id")

	if err := a.db.First(&account, accountID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Data not found",
			})
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"balance": account.Balance,
	})
}

func (a *accountImplement) Transfer(ctx *gin.Context) {
	payload := transferPayload{}
	var senderAccount, recepientAccount model.Account
	accountID := ctx.GetInt64("account_id")

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Fetch both sender and recipient accounts in a single query
	accounts := []model.Account{senderAccount, recepientAccount}
	if err := a.db.Where("account_id IN (?)", []int64{accountID, payload.TargetID}).Find(&accounts).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Data not found",
			})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	senderAccount = accounts[0]
	recepientAccount = accounts[1]

	if senderAccount.Balance < payload.Amount {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{
			"error": "Balance not enough",
		})
		return
	}

	// Start a database transaction
	tx := a.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update account balances
	senderAccount.Balance -= payload.Amount
	if err := tx.Save(&senderAccount).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	recepientAccount.Balance += payload.Amount
	if err := tx.Save(&recepientAccount).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":           "Update success",
		"amount":            payload.Amount,
		"sender_balance":    senderAccount.Balance,
		"recepient_balance": recepientAccount.Balance,
	})
}
