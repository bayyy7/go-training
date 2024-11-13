package handlers

import (
	"example/model"
	"example/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AuthInterface interface {
	AuthLogin(*gin.Context)
	AuthSignUp(*gin.Context)
	Upsert(*gin.Context)
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

	if err := bcrypt.CompareHashAndPassword([]byte(auth.Password), []byte(payload.Password)); err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Password salah",
		})
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
	payload := authPayload{}

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	if !utils.CharacterCheck(payload.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "password must be numeric",
		})
		return
	}

	existingUser := model.Auth{}
	if result := a.db.Where("username = ?", payload.Username).First(&existingUser); result.RowsAffected > 0 {
		ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"error": "username already exist",
		})
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}

	newUser := model.Auth{
		Username: payload.Username,
		Password: string(hashPassword),
	}

	result := a.db.Create(&newUser)
	if result.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User register succesfully",
	})
}

type authUpsertPayload struct {
	AccountID int64  `json:"account_id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

func (a *authImplement) Upsert(c *gin.Context) {
	payload := authUpsertPayload{}

	err := c.BindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	var account model.Account
	if err := a.db.First(&account, payload.AccountID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Account Not found",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	auth := model.Auth{
		AccountID: payload.AccountID,
		Username:  payload.Username,
		Password:  string(hashed),
	}

	result := a.db.Clauses(
		clause.OnConflict{
			DoUpdates: clause.AssignmentColumns([]string{"username", "password"}),
			Columns:   []clause.Column{{Name: "account_id"}},
		}).Create(&auth)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Create success",
		"data":    payload.Username,
	})
}
