package minibase

import (
	"database/sql"
	"fmt"
	"os"

	dev_logger "github.com/carter4299/gin_auth/my_server_config"
	_ "github.com/mattn/go-sqlite3"
)

var log = dev_logger.DevLogger()

func Open() (*sql.DB, error) {
	log.Info("Opened connection to db in Open() in minibase / queries.go ...")
	return sql.Open("sqlite3", fmt.Sprintf("%s?_auth&_auth_user=%s&_auth_pass=%s&_auth_crypt=sha256", "./mini_base/My_s_a_t.db", os.Getenv("AUTH_DB_USER"), os.Getenv("AUTH_DB_PASS")))
}
