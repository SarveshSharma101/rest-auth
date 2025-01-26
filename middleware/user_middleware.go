package middleware

import (
	"net/http"
	"rest-auth/cache"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func AuthMiddleware(c *gin.Context) {

	session, err := c.Cookie("session_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session cookie not found"})
		c.Abort()
		return
	}

	// _, err = DB.Db.GetSession(session)
	_, err = cache.GetValues(session)
	if err != nil {
		if err == redis.Nil {
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

	// session, err := DB.Db.GetSession(sessionID)
	emailId, err := cache.GetValues(sessionID)

	if err != nil {
		if err == redis.Nil {
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
	if email != emailId {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized, email does not match"})
		c.Abort()
		return
	}
	c.Next()
}
