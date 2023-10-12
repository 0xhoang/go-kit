package config

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Env       string `json:"env"`
	Port      int    `json:"port"`
	SentryDSN string `json:"sentry_dsn"`
	//db
	Db string
}

func ReadConfigAndArg() *Config {
	err := godotenv.Load()
	if err != nil {
		//log.Fatal("Error loading .env file")
	}

	fileConfig := "config.json"
	data, err := os.ReadFile("./config/" + fileConfig)
	if err != nil {
		log.Fatalln(err)
	}

	var tempCfg *Config
	if data != nil {
		err = json.Unmarshal(data, &tempCfg)
		if err != nil {
			log.Fatalf("Unmarshal err %v", err.Error())
		}
	}

	tempCfg.Env = tempCfg.Env
	tempCfg.Db = os.Getenv("MYSQL_CONNECT")

	fmt.Println("============Config===============")
	fmt.Println("env =", tempCfg.Env)
	fmt.Println("fileConfig =", fileConfig)
	fmt.Println("===========================")

	return tempCfg
}
