package auth

import (
	"crypto/sha256"
	"fmt"
	"net/http"

	sessionmanager "github.com/carter4299/gin_auth/session_manager"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	log.Info("Entered /login -> Login() in auth / login.go ...")
	userid := c.PostForm("user_id")
	if userid == "" {
		log.Error("user_id not found in context")
		c.JSON(400, gin.H{"error": "user_id not found in context"})
		return
	}
	log.Debug("user_id found in context")

	decodedPasswordInterface, exists := c.Get("decoded_password")
	if !exists {
		log.Error("decoded_password not found in context")
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	if decodedPassword, ok := decodedPasswordInterface.([]byte); ok {
		log.Debug("decoded_password found in context")
		auth_token, err := foo_login(convert_to_array(sha256.Sum256([]byte(userid))), decodedPassword)
		clear_bytes(decodedPassword)
		if err != nil {
			log.Error(err)
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		} else {
			cookie := set_cookie(auth_token)
			http.SetCookie(c.Writer, &cookie)

			c.JSON(200, gin.H{"status": "User Auth Successfull"})
		}
	} else {
		log.Error("failed type assertion for decoded_password")
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
}

func foo_login(userid []byte, hashed_password []byte) (string, error) {
	log.Info("Entered foo_login() in auth / login.go ...")

	x, y, err := check_login(userid, hashed_password)
	if err != nil {
		log.Error(err)
		return "", fmt.Errorf("failed to add user to the database")
	}

	log.Info("Cleared Locals ... Exited foo_login() ...")
	return sessionmanager.AddSession(x, y)
}
