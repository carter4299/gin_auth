package sessionmanager

import (
	"time"

	util "github.com/carter4299/gin_auth/util"
)

func AddSession(to_usernames []byte, to_usernames_key int) (string, error) {
	log.Debug("Entering AddSession() in sessionmanager / add.go ...")

	all_auth_tokens, active_auth_tokens, err := foo_add(to_usernames, to_usernames_key)
	if err != nil {
		log.Error(err)
		return "", err
	}

	is_added, err := add_user_session(all_auth_tokens, active_auth_tokens)
	if !is_added || err != nil {
		log.Error(err)
		return "", err
	}

	clear_bytes(to_usernames)
	clear_bytes(all_auth_tokens.i_D_X)
	clear_bytes(all_auth_tokens.to_usernames)
	all_auth_tokens.to_usernames_key = 0
	clear_bytes(active_auth_tokens.to_all_auth_tokens)
	active_auth_tokens.to_all_auth_tokens_key = 0
	active_auth_tokens.start_time = 0
	active_auth_tokens.close_time = 0

	log.Debug("Cleared Locals ...\tSuccessfully added session to the database.")



	return encrypt_session_ID(active_auth_tokens.private_auth_token)
}

func foo_add(to_usernames []byte, to_usernames_key int) (all_auth_tokens, active_auth_tokens, error) {
	log.Debug("Entering foo_add() in sessionmanager / add.go ...")

	var active_auth_tokens active_auth_tokens
	var all_auth_tokens all_auth_tokens
	var err error

	active_auth_tokens.private_auth_token, err = generated_session_ID()
	if err != nil {
		log.Error(err)
		return temp_fix(err)
	}
	active_auth_tokens.to_all_auth_tokens = util.GenerateLargeRand()
	active_auth_tokens.to_all_auth_tokens_key = util.GenerateSmallRand()
	active_auth_tokens.start_time = int32(time.Now().Unix())
	active_auth_tokens.close_time = int32(time.Now().Add(1 * time.Hour).Unix())


	all_auth_tokens.i_D_X = util.DeriveBytesFromHMAC(util.GetKeyFromEnv(active_auth_tokens.to_all_auth_tokens_key), active_auth_tokens.to_all_auth_tokens)
	all_auth_tokens.to_usernames = to_usernames
	all_auth_tokens.to_usernames_key = to_usernames_key

	return all_auth_tokens, active_auth_tokens, err

}

func temp_fix(err error) (all_auth_tokens, active_auth_tokens, error) {
	var active_auth_tokens active_auth_tokens

	active_auth_tokens.private_auth_token = nil
	active_auth_tokens.to_all_auth_tokens = nil
	active_auth_tokens.to_all_auth_tokens_key = 0
	active_auth_tokens.start_time = 0
	active_auth_tokens.close_time = 0

	var all_auth_tokens all_auth_tokens

	all_auth_tokens.i_D_X = nil
	all_auth_tokens.to_usernames = nil
	all_auth_tokens.to_usernames_key = 0

	return all_auth_tokens, active_auth_tokens, err
}