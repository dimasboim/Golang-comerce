package core

import (
	"Day15/config"
	"Day15/models"
	"net/http"
	"net/mail"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) {
	userData := []models.User{}
	config.Db.Find(&userData)
	c.JSON(200, gin.H{
		"status": "Success",
		"data":   userData,
	})
}

func InsertUser(c *gin.Context) {
	if c.PostForm("username") == "" || c.PostForm("fullname") == "" || c.PostForm("email") == "" || c.PostForm("address") == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": "Missing required parameter",
		})
		return
	}
	_, errEm := mail.ParseAddress(c.PostForm("email"))
	if errEm != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": "Invalid Email Format",
		})
		return
	}

	role, _ := strconv.ParseInt(c.PostForm("role"), 10, 64)

	userData := models.User{
		Username:  c.PostForm("username"),
		Firstname: c.PostForm("firstname"),
		Lastname:  c.PostForm("lastname"),
		Email:     c.PostForm("email"),
		Address:   c.PostForm("address"),
		Role:      role,
	}
	err := config.Db.Create(&userData)
	if err.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"Message": err.Error.Error(),
			"data":    userData,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  "Ok",
			"Message": "Data Insert Blog",
			"data":    userData,
		})

	}
}

func GetUserDetail(c *gin.Context) {
	userData := []models.User{}
	username := c.Param("username")
	config.Db.First(&userData, "username = ?", username)
	//	Db.Find(&blogData)
	c.JSON(200, gin.H{
		"status": "Success",
		"data":   userData,
	})
}
