package api

import (
	"time"

	"github.com/yuugure-aikouka/kyoto-common/utils"
)

type Config struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	DBAddr       string
}

func LoadConfig() Config {
	return Config{
		Addr:         utils.GetEnvString("ADDR", ":8080"),
		DBAddr:       utils.GetEnvString("DB_ADDR", "postgres://root:secret@localhost:5432/kyoto_common"),
		ReadTimeout:  time.Duration(utils.GetEnvInt("SERVER_READ_TIMEOUT", 10)),
		WriteTimeout: time.Duration(utils.GetEnvInt("SERVER_WRITE_TIMEOUT", 30)),
		IdleTimeout:  time.Duration(utils.GetEnvInt("SERVER_IDLE_TIMEOUT", 60)),
	}
}
