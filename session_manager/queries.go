package sessionmanager

import (
	"time"

	mini_base "github.com/carter4299/gin_auth/mini_base"
)

func add_user_session(all_auth_tokens all_auth_tokens, active_auth_tokens active_auth_tokens) (bool, error) {
	log.Debug("Entering add_user_session() in sessionmanager / queries.go authdb")
	db, err := mini_base.Open()
	if err != nil {
		log.Error(err)
		return false, err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO all_auth_tokens(i_D_X, to_usernames, to_usernames_key) VALUES(?, ?, ?)", all_auth_tokens.i_D_X, all_auth_tokens.to_usernames, all_auth_tokens.to_usernames_key)
	if err != nil {
		log.Error(err)
		return false, err
	}

	_, err = db.Exec("INSERT INTO active_auth_tokens(private_auth_token, to_all_auth_tokens, to_all_auth_tokens_key, start_time, close_time) VALUES(?, ?, ?, ?, ?)", active_auth_tokens.private_auth_token, active_auth_tokens.to_all_auth_tokens, active_auth_tokens.to_all_auth_tokens_key, active_auth_tokens.start_time, active_auth_tokens.close_time)
	if err != nil {
		log.Error(err)
		return false, err
	}

	log.Debug("Successfully inserted user session data into database")
	return true, nil

}

func val_session_time(prat []byte) (time.Time, error) {
	var expirationTime time.Time
	log.Debug("Entering val_session_time() in sessionmanager / queries.go authdb")
	db, err := mini_base.Open()
	if err != nil {
		log.Error(err)
		return expirationTime, err
	}
	defer db.Close()
	
	loc := db.QueryRow("SELECT close_time FROM active_auth_tokens WHERE private_auth_token = ?", prat)
	err = loc.Scan(&expirationTime)
	if err != nil {
		log.Error(err)
		return expirationTime, err
	}
	clear_bytes(prat)
	log.Debug("Cleared Locals ... Exited val_session_time() ...")

	return expirationTime, nil
}


func remove_user_session(session []byte) (bool, error) {
	log.Debug("Entering add_user_session() in sessionmanager / queries.go authdb")
	db, err := mini_base.Open()
	if err != nil {
		log.Error(err)
		return false, err
	}
	defer db.Close()


	_, err = db.Exec("DELETE FROM active_auth_tokens WHERE private_auth_token=? ", session)
	if err != nil {
		log.Error(err)
		return false, err
	}

	log.Debug("Successfully removed user session data from database")
	return true, nil

}
