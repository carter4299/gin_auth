package sessionmanager

import dev_logger "github.com/carter4299/gin_auth/my_server_config"

var log = dev_logger.DevLogger()

func clear_bytes(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

type all_auth_tokens struct {
	i_D_X        []byte
	to_usernames    []byte
	to_usernames_key     int
}

type active_auth_tokens struct {
	private_auth_token []byte
	to_all_auth_tokens  []byte
	to_all_auth_tokens_key  int
	start_time int32
	close_time int32
}
