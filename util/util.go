package util

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"os"

	dev_logger "github.com/carter4299/gin_auth/my_server_config"
)

var log = dev_logger.DevLogger()

func GenerateLargeRand() []byte {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func GenerateSmallRand() int {
	b := make([]byte, 1)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return int(b[0]%6) + 1
}

func DeriveBytesFromHMAC(keyBytes []byte, data []byte) []byte {
	h := hmac.New(sha256.New, keyBytes)
	h.Write(data)
	return h.Sum(nil)[:4]
}

func GetKeyFromEnv(smallRand int) []byte {
	keyName := fmt.Sprintf("INDEX_KEY%d", smallRand)
	return []byte(os.Getenv(keyName))
}
