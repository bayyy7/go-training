package handlers

import (
	"example/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type transactionInterface interface {
	LastTransaction(*gin.Context)
}

type transactionImplement struct {
	db *gorm.DB
}

func NewTransaction(db *gorm.DB) transactionInterface {
	return &transactionImplement{
		db: db,
	}
}

func (a *transactionImplement) LastTransaction(ctx *gin.Context) {
	var lastTransaction []model.Transaction

	id := ctx.Param("id")
	if err := a.db.Where("account_id= ?", id).Last(&lastTransaction).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "ID not found",
			})
			return
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":     "success",
		"transaction": lastTransaction,
	})
}
