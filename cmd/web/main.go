package main

import (
	"fmt"
	"log"

	"github.com/fgeck/gotth-sqlite/internal/service/config"
	"github.com/fgeck/gotth-sqlite/internal/web"
	"github.com/labstack/echo/v4"
)

func main() {
	cfgLoader := config.NewLoader()
	cfg, err := cfgLoader.LoadConfig("")
	if err != nil {
		panic(err)
	}

	e := echo.New()

	web.InitServer(e, cfg)

	log.Fatal(e.Start(fmt.Sprintf("%s:%s", cfg.App.Host, cfg.App.Port)))
}
