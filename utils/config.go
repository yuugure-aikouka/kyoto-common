package utils

import "time"

type Config struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	DBAddr       string
}

func LoadConfig() Config {
	return Config{
		Addr:         GetEnvString("ADDR", ":8080"),
		DBAddr:       GetEnvString("DB_ADDR", "postgres://root:secret@localhost:5432/kyoto_common"),
		ReadTimeout:  time.Duration(GetEnvInt("SERVER_READ_TIMEOUT", 10)),
		WriteTimeout: time.Duration(GetEnvInt("SERVER_WRITE_TIMEOUT", 30)),
		IdleTimeout:  time.Duration(GetEnvInt("SERVER_IDLE_TIMEOUT", 60)),
	}
}
