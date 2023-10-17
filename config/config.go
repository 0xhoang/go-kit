package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type JobInfo struct {
	TimesRetry int `json:"times_retry"`
}

type Config struct {
	Env       string `json:"env"`
	Port      int    `json:"port"`
	SentryDSN string `json:"sentry_dsn"`
	//db
	Db      string   `json:"db"`
	JobInfo *JobInfo `json:"job_info"`
}

func ReadConfigAndArg() *Config {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	fileConfig := "config.json"

	data, err := os.ReadFile(basepath + "/" + fileConfig)
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

	fmt.Println("============Config===============")
	fmt.Println("env =", tempCfg.Env)
	fmt.Println("fileConfig =", fileConfig)
	fmt.Println("===========================")

	return tempCfg
}
