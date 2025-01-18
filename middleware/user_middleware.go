package middleware

import (
	"net/http"
	"rest-auth/DB"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func AuthMiddleware(c *gin.Context) {

	session, err := c.Cookie("session_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session cookie not found"})
		return
	}

	_, err = DB.Db.GetSession(session)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		}
		return
	}
	c.Next()
}
