package shared

import (
	"crypto/rand"
	"log"
	"math/big"
	"os"

	"github.com/joho/godotenv"
)

type ConfigData struct {
	Key         string
	Password    string
	Name        string
	Description string
}

var Config ConfigData

func InitConfig() {
	_ = godotenv.Load()

	Config = ConfigData{
		Key:         envOrDefault("CTF_KEY", ""),
		Password:    envOrDefault("CTF_PASSWORD", ""),
		Name:        envOrDefault("CTF_NAME", "Gopher CTF"),
		Description: envOrDefault("CTF_DESCRIPTION", "A CTF platform built with Go"),
	}

	if Config.Key == "" {
		Config.Key = RandString(20)
		log.Println("CTF_KEY not set, generated random key")
	}
	if Config.Password == "" {
		Config.Password = RandString(10)
		log.Println("CTF_PASSWORD not set, generated random password")
		log.Println("AdminPass: ", Config.Password)
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func RandString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	res := make([]byte, length)
	for i := range res {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			log.Panicln(err)
		}
		res[i] = charset[num.Int64()]
	}
	return string(res)
}
