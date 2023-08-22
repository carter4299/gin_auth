package auth

import (
	"net/http"

	dev_logger "github.com/carter4299/gin_auth/my_server_config"
)

var log = dev_logger.DevLogger()

type initial_user_table struct {
	hashed_username     []byte
	to_user_password_indexs_table_x  []byte
	to_user_password_indexs_table_y int
}

type user_password_indexs_table struct {
	i_d_X []byte
	to_passwords_table_x        []byte
	to_passwords_table_y        int
	to_usernames_table_x        []byte
	to_usernames_table_y        int
}

type usernames_table struct {
	I_d_x     []byte
	username []byte
}

type passwords_table struct {
	i_D_x      []byte
	password []byte
}

func convert_to_array(input [32]byte) []byte {
	return input[:]
}

func clear_bytes(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

func clear_initial_user_table(initial_user_table initial_user_table) {
	clear_bytes(initial_user_table.hashed_username)
	clear_bytes(initial_user_table.to_user_password_indexs_table_x)
	initial_user_table.to_user_password_indexs_table_y = 0
}

func clear_user_password_indexs_table(user_password_indexs_table user_password_indexs_table) {
	clear_bytes(user_password_indexs_table.i_d_X)
	clear_bytes(user_password_indexs_table.to_passwords_table_x)
	user_password_indexs_table.to_passwords_table_y = 0
	clear_bytes(user_password_indexs_table.to_usernames_table_x)
	user_password_indexs_table.to_usernames_table_y = 0
}

func clear_usernames_table(usernames_table usernames_table) {
	clear_bytes(usernames_table.I_d_x)
	clear_bytes(usernames_table.username)
}

func clear_passwords_table(passwords_table passwords_table) {
	clear_bytes(passwords_table.i_D_x)
	clear_bytes(passwords_table.password)
}

func set_cookie(auth_token string) http.Cookie {
	maxAge := 86400
	cookie := http.Cookie{
		Name:     "auth_token",
		Value:    auth_token,
		MaxAge:   maxAge,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		// Secure:   true,
	}
	return cookie
}

