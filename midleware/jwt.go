
package midleware

import (
	"Day15/core"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func IsAdmin() gin.HandlerFunc {
	return CheckJWT(1)
}
func IsUser() gin.HandlerFunc {
	return CheckJWT(0)
}

func CheckJWT(role uint) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			token *jwt.Token
			err   error
		)

		authHeader := c.Request.Header.Get("Authorization")
		splitToken := strings.Split(authHeader, " ")

		if len(splitToken) == 2 {
			token, err = jwt.Parse(splitToken[1], func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET_KEY")), nil
			})

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "Internal Server Error",
					"message": "Error When Parsing JWT",
				})
				return
			}

		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":   "Unauthorized",
				"messages": "Token is Required",
			})
			c.Abort()
			return
		}

		if token.Valid {
			isAuth := core.IsAuthenticated(splitToken[1])
			if isAuth {
				claim, ok := token.Claims.(jwt.MapClaims)
				if !ok {
					c.JSON(http.StatusUnauthorized, gin.H{
						"status":   " Unauthorized",
						"messages": "Failed to claim this token",
					})
					c.Abort()
					return
				}
				fmt.Println(claim)
				user_id := uint(claim["user_id"].(float64))
				role_id := uint(claim["role_id"].(float64))

				if role == role_id {
					c.Set("user_id", user_id)
					c.Set("role_id", role_id)
				} else {
					c.JSON(http.StatusUnprocessableEntity, gin.H{
						"status":  "Unprocessable Entity",
						"message": "Youre not allowed to acces this endpoint",
					})
					c.Abort()
					return
				}
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"status":  "Login Success",
					"message": "Authorization required, please login/register first",
				})
				c.Abort()
				return
			}
		} else if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":   "Unauthorized",
					"messages": "That is not a token",
				})
				return
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":   "Unauthorized",
					"messages": "Token expired or not active yet",
				})
				return
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":   "Unauthorized",
					"messages": "Token Invalid",
				})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":   "Unauthorized",
				"messages": "Token Invalid",
			})
			return
		}

	}
}
