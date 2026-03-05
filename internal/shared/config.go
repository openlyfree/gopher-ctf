package shared

import (
	"crypto/rand"
	"encoding/json"
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
	err := os.MkdirAll("ctf", 0o755)
	if err != nil {
		log.Panicln("Failed to create ctf directory: ", err)
	}
	f, err := os.ReadFile("ctf/config.json")
	if err != nil {
		GenConfig()
		return
	}
	err = json.Unmarshal(f, &Config)
	if err != nil {
		log.Panicln("Failed to parse config file: ", err)
	}
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

func GenConfig() {
	cfg := ConfigData{
		Key:      RandString(20),
		Password: RandString(10),
	}

	f, err := os.Create("ctf/config.json")
	if err != nil {
		log.Panicln("failed to create config file: ", err)
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
	log.Println("AdminPass: ", cfg.Password)
	written, err := f.Write(enc)
	if err != nil || written != len(enc) {
		log.Panicln("Failed to write config file")
	}
	Config = cfg
}
