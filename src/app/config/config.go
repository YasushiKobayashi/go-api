package config

import (
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

var conf = main()

// server config
const PORT = "5000"

var HOST = conf.HOST + ":"

var URL = conf.URL
var SITE_TITLE = conf.SITE_TITLE
var HASH_SALT = conf.HASH_SALT
var SLACK_WEBHOOKURL = conf.SLACK_WEBHOOKURL
var SLACK_FACEICON = conf.SLACK_FACEICON
var SLACK_CNANNEL = conf.SLACK_CNANNEL
var SLACK_USERNAME = conf.SLACK_USERNAME
var FRONT_URL = "http://" + conf.FRONT_DOMEIN + "/"
var ALLOW_ORIGINS = "http://" + conf.FRONT_DOMEIN

// db config
const DB_TYPE = "mysql"
const DB_OPTIONS_STRING = "gorm:table_options"
const DB_OPTIONS_INTERFACE = "ENGINE=InnoDB"

var DB_URL = conf.DB_USER + ":" + conf.DB_PASS + "@tcp(" + conf.DB_HOST + ":3306)/" + conf.DB_DATABASE + "?charset=utf8&parseTime=true&loc=Local"

// JWT
var JWT_EXP = time.Now().Add(time.Hour * 168).Unix()

const UPLOAD_PATH = "static/img/"

var GOPATH = os.Getenv("GOPATH")
var PROJECT_PATH = GOPATH + "/src/app/"
var UPLOAD_DIR = PROJECT_PATH + UPLOAD_PATH

type (
	ENV struct {
		Production Config
		Develop    Config
		Test       Config
	}

	Config struct {
		DB_DATABASE      string `toml:"DB_DATABASE"`
		DB_HOST          string `toml:"DB_HOST"`
		DB_USER          string `toml:"DB_USER"`
		DB_PASS          string `toml:"DB_PASS"`
		HOST             string `toml:"HOST"`
		URL              string `toml:"URL"`
		SITE_TITLE       string `toml:"SITE_TITLE"`
		HASH_SALT        string `toml:"HASH_SALT"`
		SLACK_WEBHOOKURL string `toml:"SLACK_WEBHOOKURL"`
		SLACK_FACEICON   string `toml:"SLACK_FACEICON"`
		SLACK_USERNAME   string `toml:"SLACK_USERNAME"`
		SLACK_CNANNEL    string `toml:"SLACK_CNANNEL"`
		FRONT_DOMEIN     string `toml:"FRONT_DOMEIN"`
	}
)

func main() Config {
	env := ENV{}
	confPath := PROJECT_PATH + "config/config.tml"
	_, err := toml.DecodeFile(confPath, &env)

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
