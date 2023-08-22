package middleware

import (
	"encoding/base64"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	devlogger "github.com/carter4299/gin_auth/my_server_config"
	session_manager "github.com/carter4299/gin_auth/session_manager"
)

var log = devlogger.DevLogger()

func InitMiddleware(r *gin.Engine) *gin.Engine {
	r.Use(decode_base_64_middleware())
	r.Use(check_valid_path())
	r.Use(tokenAuthMiddleware())
	r.Use(setStaticFileMimeType())
	return r
}

func check_valid_path() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

/*
var ongoingRequests = make(map[string]bool)

	func isRequestOngoing(hash string) bool {
		_, exists := ongoingRequests[hash]
		return exists
	}

	func markRequestAsOngoing(hash string) {
		ongoingRequests[hash] = true
	}

	func unmarkRequest(hash string) {
		delete(ongoingRequests, hash)
	}

	func check_valid_path() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	func computeRequestHash(c *gin.Context) string {
		body, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		return fmt.Sprintf("%s:%s:%s", c.Request.Method, c.Request.URL.Path, string(body))
	}

	func request_guard_middleware() gin.HandlerFunc {
		return func(c *gin.Context) {
			requestHash := computeRequestHash(c)

			if isRequestOngoing(requestHash) {
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Duplicate request"})
				return
			}

			markRequestAsOngoing(requestHash)
			c.Next()
			unmarkRequest(requestHash)
		}
	}
*/
func decode_base_64_middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" {
			c.Request.ParseMultipartForm(10 << 20) // size limit 1024 bytes
			for key, values := range c.Request.PostForm {
				for i, value := range values {
					decodedValue, err := base64.StdEncoding.DecodeString(value)
					if err == nil {
						c.Request.PostForm[key][i] = string(decodedValue)
					}
				}
			}
			if file, _, err := c.Request.FormFile("password"); err == nil {
				bytes, err := io.ReadAll(file)
				if err != nil {
					c.Abort()
					return
				}

				decodedBytes, err := base64.StdEncoding.DecodeString(string(bytes))
				if err == nil {
					c.Set("decoded_password", decodedBytes)
				}
			}
		}

		c.Next()
	}
}
func setStaticFileMimeType() gin.HandlerFunc {
    return func(c *gin.Context) {
        if strings.HasSuffix(c.Request.URL.Path, ".js") {
            c.Header("Content-Type", "application/javascript")
        } else if strings.HasSuffix(c.Request.URL.Path, ".css") {
            c.Header("Content-Type", "text/css")
        }
        c.Next()
    }
}
func tokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip middleware for authentication route
		if c.FullPath() == "/" ||
			c.FullPath() == "/login" ||
			c.FullPath() == "/signup" ||
			c.FullPath() == "/assets/" ||
			c.FullPath() == "/src/" ||
			strings.HasPrefix(c.FullPath(), "/assets/") || 
			c.FullPath() == "/api/go/auth/login" ||
			c.FullPath() == "/api/go/auth/signup" {
			c.Next()
			return
		}
		cookie, err := c.Request.Cookie("auth_token")
		if err != nil {
			log.Error(err)
			log.Debug(c.Request.Cookies())
			c.JSON(http.StatusUnauthorized, "auth token required")
			c.Abort()
			return
		}


		authToken := cookie.Value

		if c.FullPath() == "/api/go/other/logout" {
			err := session_manager.RemoveSession(authToken)
			if err != nil {
				log.Error(err)
				c.JSON(http.StatusInternalServerError, "Internal server error")
				c.Abort()
				return
			} 
			invalidCookie := remove_cookie()
			c.SetCookie(invalidCookie.Name, invalidCookie.Value, invalidCookie.MaxAge, invalidCookie.Path, invalidCookie.Domain, invalidCookie.Secure, invalidCookie.HttpOnly)

			c.JSON(200, gin.H{"message": "Logged out successfully"})
			c.Abort()
			return
		}

		isit, err := session_manager.ValidateAuthToken(authToken)
		if err != nil || !isit {
			c.JSON(http.StatusUnauthorized, "invalid auth token")
			c.Abort()
			return
		}

		c.Next()
	}
}

func remove_cookie() http.Cookie {
	maxAge := -1
	cookie := http.Cookie{
		Name:     "auth_token",
		Value:    "",
		MaxAge:   maxAge,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		// Secure:   true,

	}
	return cookie
}
