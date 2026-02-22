package main

import (
	"fmt"
	"log"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"github.com/galaxy-empire-team/bridge-api/pkg/registry"
	"github.com/galaxy-empire-team/event-manager/internal/app"
	"github.com/galaxy-empire-team/event-manager/internal/config"
	"github.com/galaxy-empire-team/event-manager/internal/db"
	buildingservice "github.com/galaxy-empire-team/event-manager/internal/service/building"
	missionservice "github.com/galaxy-empire-team/event-manager/internal/service/mission"
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

	// initialize pgx infra.
	db, err := db.New(ctx, cfg.PgConn)
	if err != nil {
		return fmt.Errorf("db.New(): %w", err)
	}
	defer db.Close()

	// initialize manager that implemets storage methods inside transactions.
	txManager := txmanager.New(db)
	
	reg, err := registry.New(ctx, db.Pool)
	if err != nil {
		return fmt.Errorf("registry.New(): %w", err)
	}

	// initialize services. Use other binaries for other services as needed.
	buildingService := buildingservice.New(txManager, reg, app.ComponentLogger("buildingservice"))
	missionService := missionservice.New(txManager, reg, app.ComponentLogger("missionservice"))

	cron := cron.New()

	_, err = cron.AddFunc("@every 1s", func() {
		if err := buildingService.HandleBuilds(ctx); err != nil {
			app.Logger().Error("failed to update buildings", zap.Error(err))
		}

		if err := missionService.HandleMissions(ctx); err != nil {
			app.Logger().Error("failed to handle missions", zap.Error(err))
		}
	})
	if err != nil {
		return fmt.Errorf("cron.AddFunc(): %w", err)
	}

	cron.Start()
	defer cron.Stop()

	app.WaitShutdown(ctx)

	return nil
}
