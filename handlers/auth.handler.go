package handlers

import (
	"example/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type accountData struct {
	EMAIL    string `json:"email" binding:"required"`
	PASSWORD string `json:"password" binding:"required"`
}

func AuthAccountHandler(ctx *gin.Context) {
	var accountData accountData

	if err := ctx.ShouldBindJSON(&accountData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !utils.AlphanumericCheck(accountData.EMAIL) {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "email must be alphanumeric",
		})
		return
	}

	if !utils.NumericCheck(accountData.PASSWORD) {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "password must be numeric",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Succes",
	})

}

type AuthInterface interface {
	AuthLogin(*gin.Context)
	AuthSignUp(*gin.Context)
}

type authImplement struct{}

func NewAuth() AuthInterface {
	return &authImplement{}
}

type authLoginPayload struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (a *authImplement) AuthLogin(ctx *gin.Context) {
	bodyPayload := authLoginPayload{}

	if err := ctx.ShouldBindJSON(&bodyPayload); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	if !utils.AlphanumericCheck(bodyPayload.Email) {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "email must be alphanumeric",
		})
		return
	}

	if !utils.NumericCheck(bodyPayload.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "password must be numeric",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Succes",
	})
}

func (a *authImplement) AuthSignUp(ctx *gin.Context) {
	bodyPayload := authLoginPayload{}

	if err := ctx.ShouldBindJSON(&bodyPayload); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
}
