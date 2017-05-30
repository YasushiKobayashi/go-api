package config

import (
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

// server config
const PORT = "5000"
const HOST = "localhost:"

var conf = main()
var URL = conf.URL

// db config
const DB_TYPE = "mysql"
const DB_OPTIONS_STRING = "gorm:table_options"
const DB_OPTIONS_INTERFACE = "ENGINE=InnoDB"

var DB_URL = conf.DB_USER + ":@" + conf.DB_PASS + "/" + conf.DB_DATABASE + "?charset=utf8&parseTime=true&loc=Local"

// HASH_SALT
const HASH_SALT = "bc9765d90f8afa361a6ad379ee74396a40a6ebf7b95908f4ba2e71b39aa7537f"

// JWT
var JWT_EXP = time.Now().Add(time.Hour * 168).Unix()

const UPLOAD_DIR = "/static/img/"

type (
	ENV struct {
		Production Config
		Develop    Config
		Test       Config
	}

	Config struct {
		DB_DATABASE string `toml:"DB_DATABASE"`
		DB_HOST     string `toml:"DB_HOST"`
		DB_USER     string `toml:"DB_USER"`
		DB_PASS     string `toml:",DB_PASS"`
		URL         string `toml:",URL"`
	}
)

func main() Config {
	env := ENV{}
	_, err := toml.DecodeFile("config/config.tml", &env)
	if err != nil {
		panic(err)
	}

	res := Config{}
	osEnv := os.Getenv("DOCUMENT_ENV")
	if osEnv == "develop" {
		res = env.Develop
	} else if osEnv == "test" {
		res = env.Test
	} else {
		res = env.Production
	}
	return res
}
