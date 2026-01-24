package main

import (
	"fmt"
	"log"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"github.com/galaxy-empire-team/event-manager/internal/app"
	"github.com/galaxy-empire-team/event-manager/internal/config"
	"github.com/galaxy-empire-team/event-manager/internal/db"
	buildingservice "github.com/galaxy-empire-team/event-manager/internal/service/building"
	"github.com/galaxy-empire-team/event-manager/internal/storage/txmanager"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("config.New(): %w", err)
	}

	ctx, app, err := app.New(cfg.App)
	if err != nil {
		return fmt.Errorf("app.New(): %w", err)
	}

	// initialize pgx infra
	db, err := db.New(ctx, cfg.PgConn)
	if err != nil {
		return fmt.Errorf("db.New(): %w", err)
	}
	defer db.Close()

	// initialize manager that implemets storage methods inside transactions
	txManager := txmanager.New(db)

	// initialize services
	buildingService := buildingservice.New(txManager, app.ComponentLogger("buildingservice"))

	cron := cron.New()

	cron.AddFunc("@every 1s", func() {
		err := buildingService.UpdateBuildings(ctx)
		if err != nil {
			app.Logger().Error("buildingService.UpdateBuildings(): %w", zap.Error(err))
		}
	})

	cron.Start()
	defer cron.Stop()

	app.WaitShutdown(ctx)

	return nil
}
