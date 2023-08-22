package startup

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var log = s_log()

func open_db() (*sql.DB, error) {
	is_it, b := os.LookupEnv("AUTH_DB_USER")
	_is_it, _b := os.LookupEnv("AUTH_DB_PASS")
	if !b || is_it == "" || !_b || _is_it == "" {
		log.Error("AUTH_DB_USER not found in environment ... EXITING")
		os.Exit(1)
	}
	return sql.Open("sqlite3", fmt.Sprintf("%s?_auth&_auth_user=%s&_auth_pass=%s&_auth_crypt=sha256", "./mini_base/My_s_a_t.db", os.Getenv("AUTH_DB_USER"), os.Getenv("AUTH_DB_PASS")))
}

func user_db_init() (bool, error) {
	db, err := open_db()
	if err != nil {
		log.Error("Database could not be opened ... EXITING")
		os.Exit(1)
	}
	defer db.Close()

	init_idx_t_q := `
	CREATE TABLE IF NOT EXISTS initial_user_table(
		hashed_username BLOB PRIMARY KEY NOT NULL,
		to_user_password_indexs_table_x BLOB NOT NULL, 
		to_user_password_indexs_table_y INTEGER NOT NULL
	)`
	_, err = db.Exec(init_idx_t_q)
	if err != nil {
		return false, err
	}

	red_idx_t_q := `
	CREATE TABLE IF NOT EXISTS user_password_indexs_table(
		i_d_X BLOB PRIMARY KEY NOT NULL,
		to_passwords_table_x BLOB NOT NULL, 
		to_passwords_table_y INTEGER NOT NULL,
		to_usernames_table_x BLOB NOT NULL, 
		to_usernames_table_y INTEGER NOT NULL
	)`
	_, err = db.Exec(red_idx_t_q)
	if err != nil {
		return false, err
	}

	u_t_q := `
	CREATE TABLE IF NOT EXISTS passwords_table(
		i_D_x BLOB PRIMARY KEY NOT NULL,
		password BLOB NOT NULL
	)`
	_, err = db.Exec(u_t_q)
	if err != nil {
		return false, err
	}

	p_t_q := `
	CREATE TABLE IF NOT EXISTS usernames_table(
		I_d_x BLOB PRIMARY KEY NOT NULL,
		username BLOB NOT NULL
	)`
	_, err = db.Exec(p_t_q)
	if err != nil {
		return false, err
	}

	s_t_q := `
	CREATE TABLE IF NOT EXISTS all_auth_tokens(
		i_D_X BLOB PRIMARY KEY NOT NULL,
		to_usernames BLOB NOT NULL,
		to_usernames_key INTEGER NOT NULL
	)`
	_, err = db.Exec(s_t_q)
	if err != nil {
		return false, err
	}

	activePoolTableQuery := `
	CREATE TABLE IF NOT EXISTS active_auth_tokens(
		private_auth_token BLOB PRIMARY KEY NOT NULL,
		to_all_auth_tokens BLOB NOT NULL,
		to_all_auth_tokens_key INTEGER NOT NULL,
		start_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
		close_time TIMESTAMP DEFAULT (datetime(CURRENT_TIMESTAMP, '+1 hour'))

	)`
	_, err = db.Exec(activePoolTableQuery)
	if err != nil {
		return false, err
	}

	db.Close()
	return true, err
}

func encrypt_session_ID(sessionID []byte) (string, error) {

	rawKey, err := base64.StdEncoding.DecodeString(os.Getenv("AUTH_ENC_KEY"))
	if err != nil {
		log.Error(err)
		return "", err
	}

	block, err := aes.NewCipher(rawKey)
	if err != nil {
		log.Error(err)
		return "", err
	}

	nonce := make([]byte, 12)
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		log.Error(err)
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Error(err)
		return "", err
	}

	ciphertext := aesgcm.Seal(nil, nonce, sessionID, nil)
	encodedEncryptedData := base64.RawURLEncoding.EncodeToString(append(nonce, ciphertext...))

	return encodedEncryptedData, nil
}

func generate_sessionID() ([]byte, error) {
	log.Debug("Entering GenerateSessionID ... Creating new random session ID...")
	sessionID := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, sessionID)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debug(fmt.Sprint("Returning generated session ID: ", len(sessionID)))
	return sessionID, err
}

func generateSmallRand() int {
	b := make([]byte, 1)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return int(b[0]%6) + 1
}

func getKeyFromEnv(smallRand int) []byte {
	keyName := fmt.Sprintf("INDEX_KEY%d", smallRand)
	return []byte(os.Getenv(keyName))
}

func generateLargeRand() []byte {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func deriveBytesFromHMAC(keyBytes []byte, data []byte) []byte {
	h := hmac.New(sha256.New, keyBytes)
	h.Write(data)
	return h.Sum(nil)[:4] // Returning only the first 4 bytes
}

func base64EncodeString(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

func convert_to_array(input [32]byte) []byte {
	return input[:]
}

type dummy_user struct {
	username string
	password string
}

func ret_dummy() dummy_user {
	var dummy_user dummy_user
	dummy_user.username = "1_pochahontas_1"
	dummy_user.password = "7N3w!guy_@n_th3_bl@ck9"

	return dummy_user
}

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

type all_auth_tokens struct {
	i_D_X        []byte
	to_usernames     []byte
	to_usernames_key     int
}
type active_auth_tokens struct {
	private_auth_token    []byte
	to_all_auth_tokens  []byte
	to_all_auth_tokens_key  int
	start_time int32
	close_time int32

}

func InitDB() {
	user_db_init()

	db, err := open_db()
	if err != nil {
		log.Error(err)
		return
	}
	defer db.Close()

	dummy_user := ret_dummy()

	/*------------------------------------------------------------init_idx Table-----------------------------------------------------------------------*/

	log.Info("Assigning Variables for initial_user_table Table.... ")
	var initial_user_table initial_user_table
	
	_u := []byte(base64EncodeString(dummy_user.username))
	initial_user_table.hashed_username = convert_to_array(sha256.Sum256([]byte(_u)))
	initial_user_table.to_user_password_indexs_table_x = generateLargeRand()
	initial_user_table.to_user_password_indexs_table_y = generateSmallRand()

	log.Info("initial_user_table:", initial_user_table)

	/*------------------------------------------------------------red_idx Table-----------------------------------------------------------------------*/

	log.Info("Assigning Variables for user_password_indexs_table Table.... ")
	var user_password_indexs_table user_password_indexs_table

	user_password_indexs_table.i_d_X = deriveBytesFromHMAC(getKeyFromEnv(initial_user_table.to_user_password_indexs_table_y), initial_user_table.to_user_password_indexs_table_x)
	user_password_indexs_table.to_passwords_table_x = generateLargeRand()
	user_password_indexs_table.to_passwords_table_y = generateSmallRand()
	user_password_indexs_table.to_usernames_table_x = generateLargeRand()
	user_password_indexs_table.to_usernames_table_y = generateSmallRand()

	log.Info("user_password_indexs_table:", user_password_indexs_table)

	/*------------------------------------------------------Password Table------------------------------------------------------------------------------*/

	log.Info("Assigning Variables for passwords_table Table.... ")
	var passwords_table passwords_table

	passwords_table.i_D_x = deriveBytesFromHMAC(getKeyFromEnv(user_password_indexs_table.to_passwords_table_y), user_password_indexs_table.to_passwords_table_x)

	passwords_table.password, err = bcrypt.GenerateFromPassword([]byte(dummy_user.password), bcrypt.DefaultCost)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("passwords_table:", passwords_table)

	/*------------------------------------------------------User Table------------------------------------------------------------------------------*/

	log.Info("Assigning Variables for usernames_table Table.... ")
	var usernames_table usernames_table

	usernames_table.I_d_x = deriveBytesFromHMAC(getKeyFromEnv(user_password_indexs_table.to_usernames_table_y), user_password_indexs_table.to_usernames_table_x)
	usernames_table.username = _u

	log.Info("usernames_table:", usernames_table)

	/*------------------------------------------------------Session Table------------------------------------------------------------------------------*/

	log.Info("Assigning Variables forall_auth_tokens Table.... ")
	var all_auth_tokens all_auth_tokens

	all_auth_tokens.to_usernames = user_password_indexs_table.to_usernames_table_x
	all_auth_tokens.to_usernames_key = user_password_indexs_table.to_usernames_table_y
	
	log.Info("all_auth_tokens:", all_auth_tokens)
	
	/*------------------------------------------------------Active Pool Table------------------------------------------------------------------------------*/

	log.Info("Assigning Variables for Active Pool Table.... ")
	var active_auth_tokens active_auth_tokens
	var _ error

	active_auth_tokens.private_auth_token, _ = generate_sessionID()
	active_auth_tokens.to_all_auth_tokens = generateLargeRand()
	active_auth_tokens.to_all_auth_tokens_key = generateSmallRand()
	active_auth_tokens.start_time = int32(time.Now().Unix())
	active_auth_tokens.close_time = int32(time.Now().Add(1 * time.Hour).Unix())


	all_auth_tokens.i_D_X = deriveBytesFromHMAC(getKeyFromEnv(active_auth_tokens.to_all_auth_tokens_key), active_auth_tokens.to_all_auth_tokens)

	log.Info("Active Pool:", active_auth_tokens)

	/*------------------------------------------------------Inserting into Database------------------------------------------------------------------------------*/

	log.Info("Inserting into Database.... ")
	/*-----------------------------init_idx-------------------------*/
	_, err = db.Exec("INSERT INTO initial_user_table(hashed_username, to_user_password_indexs_table_x, to_user_password_indexs_table_y) VALUES(?, ?, ?)", initial_user_table.hashed_username, initial_user_table.to_user_password_indexs_table_x, initial_user_table.to_user_password_indexs_table_y)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("Successfully inserted users data into database.")
	/*-----------------------------red_idx-------------------------*/
	
	_, err = db.Exec("INSERT INTO user_password_indexs_table(i_d_X, to_passwords_table_x, to_passwords_table_y, to_usernames_table_x, to_usernames_table_y) VALUES(?, ?, ?, ?, ?)", user_password_indexs_table.i_d_X, user_password_indexs_table.to_passwords_table_x, user_password_indexs_table.to_passwords_table_y, user_password_indexs_table.to_usernames_table_x, user_password_indexs_table.to_usernames_table_y)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("Successfully inserted user data into database.")

	/*-----------------------------USER-------------------------*/
	_, err = db.Exec("INSERT INTO usernames_table(I_d_x, username) VALUES(?, ?)", usernames_table.I_d_x, usernames_table.username)
	if err != nil {
		log.Error(err)
		return
	}
	/*-----------------------------PASSWORDS-------------------------*/
	_, err = db.Exec("INSERT INTO passwords_table(i_D_x, password) VALUES(?, ?)", passwords_table.i_D_x, passwords_table.password)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("Successfully inserted password data into database.")

	/*-----------------------------SESSIONS-------------------------*/

	_, err = db.Exec("INSERT INTO all_auth_tokens(i_D_X, to_usernames, to_usernames_key) VALUES(?, ?, ?)", all_auth_tokens.i_D_X, all_auth_tokens.to_usernames, all_auth_tokens.to_usernames_key)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("Successfully inserted user instances data into database.")
	/*-----------------------------ACTIVE POOL-------------------------*/
	_, err = db.Exec("INSERT INTO active_auth_tokens(private_auth_token, to_all_auth_tokens, to_all_auth_tokens_key, start_time, close_time) VALUES(?, ?, ?, ?, ?)", active_auth_tokens.private_auth_token, active_auth_tokens.to_all_auth_tokens, active_auth_tokens.to_all_auth_tokens_key, active_auth_tokens.start_time, active_auth_tokens.close_time)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Successfully inserted all data into database. \n")

}