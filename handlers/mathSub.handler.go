package handlers

import (
	"example/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type mathSubData struct {
	A int `json:"a" binding:"required"` // annotation
	B int `json:"b" binding:"required"` // annotation
}

func MathSubHandler(ctx *gin.Context) {
	var mathData mathSubData

	if err := ctx.ShouldBindJSON(&mathData); err != nil { // penggunaan pointer karena function ShouldBindJSON tidak me-return data
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": utils.MagicSub(mathData.A, mathData.B),
	})
}
