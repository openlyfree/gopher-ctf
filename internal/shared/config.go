package shared

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
)

type ConfigData struct {
	Key      string `json:"key"`
	Password string `json:"password"`
}

var Config ConfigData

func InitConfig() {
	f, err := os.ReadFile("config.json")
	if err != nil {
		GenConfig()
		return
	}
	err = json.Unmarshal(f, &Config)
	if err != nil {
		log.Fatal("Failed to parse config file:", err)
	}
}
func RandString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	res := make([]byte, length)
	for i := range res {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			log.Fatal(err)
		}
		res[i] = charset[num.Int64()]
	}
	return string(res)
}
func GenConfig() {
	cfg := ConfigData{
		Key:      RandString(20),
		Password: RandString(10),
	}

	f, err := os.Create("config.json")
	if err != nil {
		log.Fatalf("failed to create config file: %v", err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			return
		}
	}(f)

	enc, err := json.Marshal(cfg)
	if err != nil {
		return
	}
	fmt.Println("AdminPass: ", cfg.Password)
	written, err := f.Write(enc)
	if err != nil || written != len(enc) {
		log.Fatal("Failed to write config file")
	}
	Config = cfg
}
