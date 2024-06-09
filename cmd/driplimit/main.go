package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/i4n-co/driplimit/pkg/api"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	initAdminFlag := flag.Bool("init-admin", false, "initialize an admin service key")
	printDefaultsFlag := flag.Bool("print-defaults", false, "print the default configuration")
	configPathFlag := flag.String("config", "", "path to the configuration file")
	flag.Parse()

	cfg, err := loadConfig(ctx, *configPathFlag)
	if err != nil {
		fmt.Println("[ERROR] failed to load config: ", err)
		os.Exit(1)
	}

	if *initAdminFlag {
		err = initAdmin(ctx, cfg)
		if err != nil {
			cfg.Logger().Error("failed to initialize admin", "err", err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}

	if *printDefaultsFlag {
		cfg.PrintDefaults()
		os.Exit(0)
	}

	service, err := initService(ctx, cfg)
	if err != nil {
		cfg.Logger().Error("failed to initialize service", "err", err)
		os.Exit(1)
	}

	api := api.New(cfg, service)
	go func() {
		err := api.Listen(cfg.AddrPort())
		if err != nil {
			cfg.Logger().Error("failed to listen", "err", err)
		}
		stop()
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = api.ShutdownWithContext(ctx)
	if err != nil {
		cfg.Logger().Error("failed to shutdown api", "err", err)
		os.Exit(1)
	}
}
