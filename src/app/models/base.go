// model base setting
package models

import (
	"app/config"
	"database/sql"
	"encoding/json"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	explain "github.com/kyokomi/gorm-explain"
)

// DB connect database. setting logging.
func DB() *gorm.DB {
	db, err := gorm.Open(config.DB_TYPE, config.DB_URL)
	if err != nil {
		panic("failed to connect database")
	}

	logPath := config.PROJECT_PATH + "log/gorm.log"
	file, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	log.SetOutput(file)

	db.SingularTable(true)
	db.LogMode(true)
	db = db.Debug()
	db.SetLogger(log.New(file, "", 0))
	db.Callback().Query().Register("explain", explain.Callback)

	return db
}

type (
	NullString struct {
		sql.NullString
	}
)

func (s *NullString) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return json.Marshal(s.String)
	} else {
		return json.Marshal(nil)
	}
}

func (s *NullString) UnmarshalJSON(data []byte) error {
	var str string
	json.Unmarshal(data, &str)
	s.String = str
	s.Valid = str != ""
	return nil
}

func NewNullString(s string) NullString {
	return NullString{sql.NullString{s, s != ""}}
}
