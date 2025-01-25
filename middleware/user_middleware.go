package middleware

import (
	"fmt"
	"net/http"
	"rest-auth/DB"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func AuthMiddleware(c *gin.Context) {

	session, err := c.Cookie("session_id")
	if err != nil {
		fmt.Println("--->", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session cookie not found"})
		c.Abort()
		return
	}

	_, err = DB.Db.GetSession(session)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		}
		c.Abort()
		return
	}
	c.Next()
}

func UpdateAuthMiddleware(c *gin.Context) {

	sessionID, err := c.Cookie("session_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session cookie not found"})
		c.Abort()
		return
	}

	session, err := DB.Db.GetSession(sessionID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		}
		c.Abort()
		return
	}

	email := c.Param("emailId")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required in the path"})
		c.Abort()
		return
	}
	if email != session.Email {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized, email does not match"})
		c.Abort()
		return
	}
	c.Next()
}
