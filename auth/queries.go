package auth

import (
	"database/sql"
	"fmt"

	mini_base "github.com/carter4299/gin_auth/mini_base"
	util "github.com/carter4299/gin_auth/util"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func store_signup(initial_user_table initial_user_table, user_password_indexs_table user_password_indexs_table, usernames_table usernames_table, passwords_table passwords_table) error {
	log.Debug("Entering store_login() in auth / queries.go ...")

	_db, err := mini_base.Open()
	if err != nil {
		log.Error(err)
		return fmt.Errorf("failed to open database")
	}
	defer _db.Close()

	err = add_initial_user_table(_db, initial_user_table)
	if err != nil {
		log.Error(err)
		return fmt.Errorf("failed to add initial_user_table to database")
	}

	err = add_user_password_indexs_table(_db, user_password_indexs_table)
	if err != nil {
		log.Error(err)
		return fmt.Errorf("failed to add user_password_indexs_table to database")
	}

	err = add_user(_db, usernames_table)
	if err != nil {
		log.Error(err)
		return fmt.Errorf("failed to add user to database")
	}

	err = add_password(_db, passwords_table)
	if err != nil {
		log.Error(err)
		return fmt.Errorf("failed to add password to database")
	}

	log.Debug("Successfully added login data to database.")

	return nil
}

func check_login(userid []byte, password []byte) ([]byte, int, error) {
	log.Debug("Entering check_login() in auth / queries.go ...")
	_db, err := mini_base.Open()
	if err != nil {
		log.Error(err)
		return nil, 0, fmt.Errorf("failed to open database")
	}
	defer _db.Close()

	user_password_indexs_table, err := get_user_password_indexs_table(_db, userid)
	if err != nil {
		log.Error(err)
		return nil, 0, fmt.Errorf("failed to get user_password_indexs_table from database")
	}

	hashed_password, err := get_password(_db, user_password_indexs_table.to_passwords_table_x, user_password_indexs_table.to_passwords_table_y)
	if err != nil {
		log.Error(err)
		return nil, 0, fmt.Errorf("failed to get password from database")
	}

	err = bcrypt.CompareHashAndPassword((hashed_password), (password))
	if err != nil {
		log.Error(err)
		return nil, 0, fmt.Errorf("failed to get password from database")
	}

	x := user_password_indexs_table.to_usernames_table_x
	y := user_password_indexs_table.to_usernames_table_y

	clear_user_password_indexs_table(user_password_indexs_table)
	clear_bytes(hashed_password)
	clear_bytes(password)
	clear_bytes(userid)
	log.Debug("Cleared Locals ... Exited check_login() ..")
	return x, y, nil
}

func add_initial_user_table(db *sql.DB, initial_user_table initial_user_table) error {
	log.Debug("Entering add_initial_user_table() in auth / queries.go ...")

	_, err := db.Exec("INSERT INTO initial_user_table(hashed_username, to_user_password_indexs_table_x, to_user_password_indexs_table_y) VALUES(?, ?, ?)", initial_user_table.hashed_username, initial_user_table.to_user_password_indexs_table_x, initial_user_table.to_user_password_indexs_table_y)
	if err != nil {
		log.Error(err)
		return fmt.Errorf("failed to insert add_initial_user_table data into database")
	}

	log.Debug("Successfully inserted add_initial_user_table data into database.")
	return nil
}

func add_user_password_indexs_table(db *sql.DB, user_password_indexs_table user_password_indexs_table) error {
	log.Debug("Entering add_user_password_indexs_table() in auth / queries.go ...")

	_, err := db.Exec("INSERT INTO user_password_indexs_table(i_d_X, to_passwords_table_x, to_passwords_table_y, to_usernames_table_x, to_usernames_table_y) VALUES(?, ?, ?, ?, ?)", user_password_indexs_table.i_d_X, user_password_indexs_table.to_passwords_table_x, user_password_indexs_table.to_passwords_table_y, user_password_indexs_table.to_usernames_table_x, user_password_indexs_table.to_usernames_table_y)
	if err != nil {
		log.Error(err)
		return fmt.Errorf("failed to insert add_user_password_indexs_table data into database")
	}
	log.Info("Successfully inserted user data into database.")
	return nil
}

func get_user_password_indexs_table(db *sql.DB, userid []byte) (user_password_indexs_table, error) {
	log.Debug("Entering get_user_password_indexs_table() in auth / queries.go ...")

	var initial_user_table initial_user_table
	var user_password_indexs_table user_password_indexs_table

	initial_user_table.hashed_username = userid
	loc := db.QueryRow("SELECT to_user_password_indexs_table_x, to_user_password_indexs_table_y FROM initial_user_table WHERE hashed_username=?", userid)
	err := loc.Scan(&initial_user_table.to_user_password_indexs_table_x, &initial_user_table.to_user_password_indexs_table_y)
	if err != nil {
		log.Error(err)
		return user_password_indexs_table, fmt.Errorf("failed to get initial_user_table data from database")
	}

	row := db.QueryRow("SELECT i_d_X, to_passwords_table_x, to_passwords_table_y, to_usernames_table_x, to_usernames_table_y FROM user_password_indexs_table WHERE i_d_X=?", util.DeriveBytesFromHMAC(util.GetKeyFromEnv(initial_user_table.to_user_password_indexs_table_y), initial_user_table.to_user_password_indexs_table_x))
	err = row.Scan(&user_password_indexs_table.i_d_X, &user_password_indexs_table.to_passwords_table_x, &user_password_indexs_table.to_passwords_table_y, &user_password_indexs_table.to_usernames_table_x, &user_password_indexs_table.to_usernames_table_y)
	if err != nil {
		log.Error(err)
		return user_password_indexs_table, fmt.Errorf("failed to get user_password_indexs_table data from database")
	}
	log.Info("Successfully retrieved user_password_indexs_table data from database.")
	return user_password_indexs_table, nil
}

func add_user(db *sql.DB, usernames_table usernames_table) error {
	log.Debug("Entering add_user() in auth / queries.go ...")

	_, err := db.Exec("INSERT INTO usernames_table(I_d_x, username) VALUES(?, ?)", usernames_table.I_d_x, usernames_table.username)
	if err != nil {
		log.Error(err)
		return fmt.Errorf("failed to insert add_user data into database")
	}
	log.Info("Successfully inserted user data into database.")
	return nil
}

func add_password(db *sql.DB, passwords_table passwords_table) error {
	log.Debug("Entering add_password() in auth / queries.go ...")

	_, err := db.Exec("INSERT INTO passwords_table(i_D_x, password) VALUES(?, ?)", passwords_table.i_D_x, passwords_table.password)
	if err != nil {
		log.Error(err)
		return fmt.Errorf("failed to insert add_password data into database")
	}
	log.Info("Successfully inserted password data into database.")
	return nil
}

func get_password(db *sql.DB, x []byte, y int) ([]byte, error) {
	log.Debug("Entering get_papassword() in auth / queries.go ...")

	var password []byte
	row := db.QueryRow("SELECT password FROM passwords_table WHERE i_D_x=?", util.DeriveBytesFromHMAC(util.GetKeyFromEnv(y), x))
	err := row.Scan(&password)
	if err != nil {
		log.Error(err)
		return password, fmt.Errorf("failed to get password data from database")
	}
	log.Info("Successfully retrieved password data from database.")
	return password, nil
}
