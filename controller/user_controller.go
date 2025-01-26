package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rest-auth/DB"
	"rest-auth/cache"
	"rest-auth/datamodel"
	"rest-auth/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUsers(ctx *gin.Context) {
	emailId := ctx.Param("emailId")
	if emailId != "" {
		value, err := cache.GetValues(emailId)
		fmt.Println("------------------------")
		switch err {
		case nil:
			fmt.Println("----entering nil")
			user := &datamodel.User{}
			err := json.Unmarshal([]byte(value), user)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(http.StatusOK, gin.H{"user": user})
			return
		case redis.Nil:
			fmt.Println("----entering redis nil")

			user, err := DB.Db.GetUserByEmail(emailId)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					ctx.JSON(http.StatusNotFound, gin.H{"error": "No users found"})
					return

				}
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			userByte, err := json.Marshal(user)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			cache.SetValues(emailId, string(userByte), 3*time.Second)
			ctx.JSON(http.StatusOK, gin.H{"user": user})
			return
		default:
			fmt.Println("----entering default")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	users, err := DB.Db.GetUsers()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "No users found"})
			return

		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"users": users})

}

func RegisterUser(c *gin.Context) {
	user := &datamodel.User{}
	err := c.ShouldBind(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.Email == "" || user.Password == "" || user.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email, Username and Password are required"})
		return
	}

	_, err = DB.Db.GetUserByEmail(user.Email)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

	if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// pwd, err := utils.Decrypt(user.Password)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }
	user.Password, err = utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = DB.Db.InsertUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func LoginUser(c *gin.Context) {
	type Login struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	login := &Login{}

	if err := c.ShouldBind(login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := DB.Db.GetUserByEmail(login.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !utils.ComparePasswordHash(user.Password, login.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	sessionId := utils.GenerateSessionId()
	// session := &datamodel.Session{SessionId: sessionId, Email: login.Email}
	// err = DB.Db.InsertSession(session)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }
	err = cache.SetValues(sessionId, login.Email, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.SetCookie("session_id", sessionId, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "User logged in successfully", "session_id": sessionId})

}

func UpdateUser(ctx *gin.Context) {
	email := ctx.Param("emailId")
	if email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email is required in the path"})
		return
	}

	user := &datamodel.User{}
	err := ctx.ShouldBind(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	count, err := DB.Db.UpdateUser(email, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if count == 0 {
		err := DB.Db.InsertUser(user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return

		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func PatchUser(ctx *gin.Context) {
	email := ctx.Param("emailId")
	if email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email is required in the path"})
	}
	user := &datamodel.User{}
	err := ctx.ShouldBind(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if user.Username != "" {
		count, err := DB.Db.PatchUser(email, "username", user.Username)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else if count == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
	} else if user.Profession != "" {
		count, err := DB.Db.PatchUser(email, "profession", user.Profession)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else if count == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})

}

func DeleteUser(ctx *gin.Context) {
	sessionId, err := ctx.Cookie("session_id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "sessionId is required"})
		return
	}
	email := ctx.Param("emailId")
	if email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email is required in the path"})
	}

	count, err := DB.Db.DeleteUser(email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if count == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// if err := DB.Db.DeleteSession(email); err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }
	if err := cache.DeleteValues(sessionId); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})

}

func LogoutUser(c *gin.Context) {
	sessionId, err := c.Cookie("session_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sessionId is required"})
		return
	}

	// if err := DB.Db.DeleteSession(email); err != nil && err != mongo.ErrNoDocuments {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }
	if err := cache.DeleteValues(sessionId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User logged out successfully"})
}
