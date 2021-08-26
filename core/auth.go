package core

import (
	"Day15/config"
	"Day15/models"
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/mail"

	"github.com/gin-gonic/gin"

	 
	"encoding/base64"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func createToken(user *models.User) string {
	key := []byte(os.Getenv("JWT_SECRET_KEY"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   user.ID,
		"role_id":   user.Role,
		"exp":       time.Now().AddDate(0, 0, 1).Unix(),
		"issued_at": time.Now().Unix(),
	})

	tokenString, _ := token.SignedString(key)
	fmt.Printf("%v\n", tokenString)
	return tokenString
}

func Register(c *gin.Context) {
	var (
		user models.User
		err  error
	)
	if c.PostForm("username") == "" || c.PostForm("firstname") == "" || c.PostForm("email") == "" || c.PostForm("password") == "" || c.PostForm("address") == "" || c.PostForm("lastname") == "" || c.PostForm("role") == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": "Missing required parameter",
		})
		c.Abort()
		return
	}
	if !validFormatEmail(c.PostForm("email")) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": "Invalid Email format",
		})
		c.Abort()
		return
	}
	user.Username = c.PostForm("username")
	user.Firstname = c.PostForm("firstname")
	user.Email = c.PostForm("email")
	user.Password = encryptPass(c.PostForm("password"))
	user.Address = c.PostForm("address")
	user.Lastname = c.PostForm("lastname")
	user.Role, err = strconv.ParseInt(c.PostForm("role"), 10, 64)

	var usercheck models.User
	config.Db.First(&usercheck, "username = ? ", user.Username)
	if usercheck.Username == user.Username {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": "username already exist",
		})
		c.Abort()
		return
	}

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":   "Unauthorized",
			"messages": "Error while parsing role",
		})
		c.Abort()
		return
	}

	result := config.Db.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": "Error saving to database",
		})
		c.Abort()
		return
	}
	token := createToken(&user)

	var userToken models.User
	result = config.Db.First(&userToken, "id = ?", user.ID)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": "Error when getting user data ",
		})
		c.Abort()
		return
	}
	userToken.Token = token
	result = config.Db.Save(&userToken)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": "Error when creating token and saving to database",
		})
		c.Abort()
		return
	}
	c.Set("user_token", token)
	c.JSON(http.StatusOK, gin.H{
		"status": "Success Registered",
		"user":   user,
		"token":  token,
	})
}

func Login(c *gin.Context) {
	var (
		user models.User
	)
	if c.PostForm("username") == "" || c.PostForm("password") == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": "Missing required parameter",
		})
		c.Abort()
		return
	}
	user.Username = c.PostForm("username")
	user.Password = encryptPass(c.PostForm("password"))

	config.Db.First(&user, "username = ? And Password =?", user.Username, user.Password)

	token := createToken(&user)

	user.Token = token

	result := config.Db.Save(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": "Error when creating token and saving to database",
		})
		c.Abort()
		return
	}
	c.Set("user_token", token)
	c.JSON(http.StatusOK, gin.H{
		"status": "Success Login",
		"user":   user,
		"token":  token,
	})
}
func Logout(c *gin.Context) {
	var (
		user models.User
	)
	user.Username = c.PostForm("username")

	config.Db.First(&user, "username = ?  ", user.Username)

	//token := createToken(&user)

	user.Token = ""

	result := config.Db.Save(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": "Error when creating token and saving to database",
		})
		c.Abort()
		return
	}
	c.Set("user_token", "")
	c.JSON(http.StatusOK, gin.H{
		"status": "Success Logout",
		"user":   user,
		"token":  "",
	})
}

func encryptPass(password string) string {
	h := sha256.Sum256([]byte(password))
	return base64.StdEncoding.EncodeToString(h[:])
}
func validFormatEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
func IsAuthenticated(token string) bool {
	var user models.User
	result := config.Db.Find(&user, "token = ?", token)
	if result.Error != nil {
		panic("error when checking token")
	} else if user.ID == 0 {
		return false
	}
	return true
}
