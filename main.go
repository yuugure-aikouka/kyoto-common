package main

import (
	"log"

	"github.com/yuugure-aikouka/kyoto-common/api"
	"github.com/yuugure-aikouka/kyoto-common/utils"
)

func main() {
	cfg := api.Config{
		Addr: utils.GetEnvString("ADDR", ":8080"),
	}

	srv := api.NewServer(cfg)

	log.Fatal(srv.Start())
}
