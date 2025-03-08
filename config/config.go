package config

import "time"

type Config struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	DBAddr       string
	AppBaseUrl   string
	SMTPHost     string
	SMTPUsername string
	SMTPPassword string
	EmailSender  string
}

func LoadConfig() Config {
	return Config{
		Addr:         GetEnvString("ADDR", ":8080"),
		DBAddr:       GetEnvString("DB_ADDR", "postgres://root:secret@localhost:5432/kyoto_common"),
		ReadTimeout:  time.Duration(GetEnvInt("SERVER_READ_TIMEOUT", 10)),
		WriteTimeout: time.Duration(GetEnvInt("SERVER_WRITE_TIMEOUT", 30)),
		IdleTimeout:  time.Duration(GetEnvInt("SERVER_IDLE_TIMEOUT", 60)),
		AppBaseUrl:   GetEnvString("APP_BASE_URL", ""),
		SMTPHost:     GetEnvString("SMTP_HOST", ""),
		SMTPUsername: GetEnvString("SMTP_USERNAME", ""),
		SMTPPassword: GetEnvString("SMTP_PASSWORD", ""),
		EmailSender:  GetEnvString("EMAIL_SENDER", ""),
	}
}
