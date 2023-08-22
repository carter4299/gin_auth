package auth

import (
	"crypto/sha256"
	"fmt"
	"net/http"

	sessionmanager "github.com/carter4299/gin_auth/session_manager"
	util "github.com/carter4299/gin_auth/util"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	log.Info("Entered /signup -> signup() in auth / signup.go ...")
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

		hashed_password, err := bcrypt.GenerateFromPassword(decodedPassword, bcrypt.DefaultCost)
		if err != nil {
			clear_bytes(decodedPassword)
			log.Error(err)
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}
		auth_token, err := foo_signup(convert_to_array(sha256.Sum256([]byte(userid))), hashed_password)
		clear_bytes(hashed_password)
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

func foo_signup(userid []byte, hashed_password []byte) (string, error) {
	log.Info("Entered foo_signup() in auth / signup.go ...")

	initial_user_table, err := fill_initial_user_table(userid)
	if err != nil {
		log.Error(err)
		return "", err
	}

	user_password_indexs_table, err := fill_user_password_indexs_table(initial_user_table)
	if err != nil {
		log.Error(err)
		return "", err
	}

	usernames_table, err := fill_user(user_password_indexs_table, userid)
	if err != nil {
		log.Error(err)
		return "", err
	}

	passwords_table, err := fill_password(user_password_indexs_table, hashed_password)
	if err != nil {
		log.Error(err)
		return "", err
	}

	err = store_signup(initial_user_table, user_password_indexs_table, usernames_table, passwords_table)
	if err != nil {
		log.Error(err)
		return "", err
	}

	x := user_password_indexs_table.to_usernames_table_x
	y := user_password_indexs_table.to_usernames_table_y

	clear_initial_user_table(initial_user_table)
	clear_user_password_indexs_table(user_password_indexs_table)
	clear_usernames_table(usernames_table)
	clear_passwords_table(passwords_table)

	log.Info("Cleared Locals ... Exited foo_signup() ...")
	return sessionmanager.AddSession(x, y)
}

func fill_initial_user_table(userid []byte) (initial_user_table, error) {
	log.Info("Entered fill_initial_user_table() in auth / signup.go ...")
	var initial_user_table initial_user_table

	initial_user_table.hashed_username = userid
	if initial_user_table.hashed_username != nil {
		initial_user_table.to_user_password_indexs_table_x = util.GenerateLargeRand()
		if initial_user_table.to_user_password_indexs_table_x != nil {
			initial_user_table.to_user_password_indexs_table_y = util.GenerateSmallRand()
			if initial_user_table.to_user_password_indexs_table_y != 0 {
				log.Debug("Exited fill_initial_user_table() ...")
				return initial_user_table, nil
			} else {
				return initial_user_table, fmt.Errorf("failed to generate random key for user_password_indexs_table")
			}
		} else {
			return initial_user_table, fmt.Errorf("failed to generate random to_user_password_indexs_table")
		}
	} else {
		return initial_user_table, fmt.Errorf("failed to generate random user_id")
	}
}

func fill_user_password_indexs_table(initial_user_table initial_user_table) (user_password_indexs_table, error) {
	log.Info("Entered fill_user_password_indexs_table() in auth / signup.go ...")
	var user_password_indexs_table user_password_indexs_table

	user_password_indexs_table.i_d_X = util.DeriveBytesFromHMAC(util.GetKeyFromEnv(initial_user_table.to_user_password_indexs_table_y), initial_user_table.to_user_password_indexs_table_x)
	if user_password_indexs_table.i_d_X != nil {
		user_password_indexs_table.to_passwords_table_x = util.GenerateLargeRand()
		if user_password_indexs_table.to_passwords_table_x != nil {
			user_password_indexs_table.to_passwords_table_y = util.GenerateSmallRand()
			if user_password_indexs_table.to_passwords_table_y != 0 {
				user_password_indexs_table.to_usernames_table_x = util.GenerateLargeRand()
				if user_password_indexs_table.to_usernames_table_x != nil {
					user_password_indexs_table.to_usernames_table_y = util.GenerateSmallRand()
					if user_password_indexs_table.to_usernames_table_y != 0 {
						log.Debug("Exited fill_user_password_indexs_table() ...")
						return user_password_indexs_table, nil
					} else {
						return user_password_indexs_table, fmt.Errorf("failed to generate random key for user_password_indexs_table")
					}
				} else {
					return user_password_indexs_table, fmt.Errorf("failed to generate random to_usernames_table_x")
				}
			} else {
				return user_password_indexs_table, fmt.Errorf("failed to generate random key for user_password_indexs_table")
			}
		} else {
			return user_password_indexs_table, fmt.Errorf("failed to generate random to_passwords_table_x")
		}
	} else {
		return user_password_indexs_table, fmt.Errorf("failed to generate random from_initial_user_table")
	}
}

func fill_user(user_password_indexs_table user_password_indexs_table, userid []byte) (usernames_table, error) {
	log.Info("Entered fill_user() in auth / signup.go ...")
	var usernames_table usernames_table

	usernames_table.I_d_x = util.DeriveBytesFromHMAC(util.GetKeyFromEnv(user_password_indexs_table.to_usernames_table_y), user_password_indexs_table.to_usernames_table_x)
	if usernames_table.I_d_x != nil {
		usernames_table.username = userid
		if usernames_table.username != nil {
			log.Debug("Exited fill_user() ...")
			return usernames_table, nil
		} else {
			return usernames_table, fmt.Errorf("failed to generate random user_id")
		}
	} else {
		return usernames_table, fmt.Errorf("failed to generate random idx")
	}
}

func fill_password(user_password_indexs_table user_password_indexs_table, hashed_password []byte) (passwords_table, error) {
	log.Info("Entered fill_password() in auth / signup.go ...")
	var passwords_table passwords_table

	passwords_table.i_D_x = util.DeriveBytesFromHMAC(util.GetKeyFromEnv(user_password_indexs_table.to_passwords_table_y), user_password_indexs_table.to_passwords_table_x)
	if passwords_table.i_D_x != nil {
		passwords_table.password = hashed_password
		if passwords_table.password != nil {
			log.Debug("Exited fill_password() ...")
			return passwords_table, nil
		} else {
			return passwords_table, fmt.Errorf("failed to generate random password")
		}
	} else {
		return passwords_table, fmt.Errorf("failed to generate random idx")
	}
}
