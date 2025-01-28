package main

import (
	"log"
	"time"

	"github.com/yuugure-aikouka/kyoto-common/api"
	"github.com/yuugure-aikouka/kyoto-common/utils"
)

func main() {
	cfg := api.Config{
		Addr:         utils.GetEnvString("ADDR", ":8080"),
		ReadTimeout:  time.Duration(utils.GetEnvInt("SERVER_READ_TIMEOUT", 10)),
		WriteTimeout: time.Duration(utils.GetEnvInt("SERVER_WRITE_TIMEOUT", 30)),
		IdleTimeout:  time.Duration(utils.GetEnvInt("SERVER_IDLE_TIMEOUT", 60)),
	}

	srv := api.NewServer(cfg)

	log.Fatal(srv.Start())
}
