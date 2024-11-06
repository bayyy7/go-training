package handlers

import (
	"example/model"
	"example/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AuthInterface interface {
	AuthLogin(*gin.Context)
	AuthSignUp(*gin.Context)
}

type authImplement struct {
	db     *gorm.DB
	jwtKey []byte
}

func NewAuth(db *gorm.DB, jwtKey []byte) AuthInterface {
	return &authImplement{
		db,
		jwtKey,
	}
}

type authPayload struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (a *authImplement) createJWT(auth *model.Auth) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["auth_id"] = auth.AuthID
	claims["account_id"] = auth.AccountID
	claims["username"] = auth.Username
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix()

	tokenString, err := token.SignedString(a.jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *authImplement) AuthLogin(ctx *gin.Context) {
	payload := authPayload{}

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	auth := model.Auth{}
	if err := a.db.Where("username = ?", payload.Username).First(&auth).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Username tidak ditemukan",
			})
			return
		}
		return
	}

	if err := a.db.Where("password = ?", payload.Password).First(&auth).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Password Salah",
			})
			return
		}
		return
	}

	token, err := a.createJWT(&auth)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"token":   token,
	})
}

func (a *authImplement) AuthSignUp(ctx *gin.Context) {
	payload := model.Auth{}

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	if !utils.AlphanumericCheck(payload.Username) {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "email must be alphanumeric",
		})
		return
	}

	if !utils.NumericCheck(payload.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "password must be numeric",
		})
		return
	}

	result := a.db.Create(&payload)
	if result.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Succes",
	})
}
