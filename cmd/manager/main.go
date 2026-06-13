package main

import (
	"fmt"
	"log"

	"github.com/galaxy-empire-team/bridge-api/pkg/registry"
	"github.com/galaxy-empire-team/event-manager/internal/app"
	"github.com/galaxy-empire-team/event-manager/internal/clients/bridgeapi"
	"github.com/galaxy-empire-team/event-manager/internal/config"
	"github.com/galaxy-empire-team/event-manager/internal/db"
	"github.com/galaxy-empire-team/event-manager/internal/repository"
	buildingservice "github.com/galaxy-empire-team/event-manager/internal/service/building"
	fleetconstructionservice "github.com/galaxy-empire-team/event-manager/internal/service/fleetconsturction"
	missionservice "github.com/galaxy-empire-team/event-manager/internal/service/mission"
	researchservice "github.com/galaxy-empire-team/event-manager/internal/service/research"
	researchstorage "github.com/galaxy-empire-team/event-manager/internal/storage/research"
	"github.com/galaxy-empire-team/event-manager/internal/storage/txmanager"
	"github.com/galaxy-empire-team/event-manager/pkg/worker"
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

	// initialize clients.
	bridgeAPIClient, err := bridgeapi.New(cfg.BridgeAPIClient)
	if err != nil {
		return fmt.Errorf("bridgeapi.New(): %w", err)
	}
	defer bridgeAPIClient.Close()

	// initialize manager that implemets storage methods inside transactions.
	txManager := txmanager.New(db)
	researchStorage := researchstorage.New(db)

	reg, err := registry.New(ctx, db.Pool)
	if err != nil {
		return fmt.Errorf("registry.New(): %w", err)
	}

	// initialize repositories
	repo := repository.New(researchStorage, reg)

	// initialize services. Use other binaries for other services as needed.
	buildingService := buildingservice.New(txManager, reg, app.ComponentLogger("buildingservice"))
	missionService := missionservice.New(txManager, bridgeAPIClient, repo, reg, app.ComponentLogger("missionservice"))
	researchService := researchservice.New(txManager, reg, app.ComponentLogger("researchservice"))
	fleetConstructionService := fleetconstructionservice.New(txManager, app.ComponentLogger("fleetConstructionService"))

	worker.StartWorker(
		ctx,
		cfg.BuildingWorker,
		buildingService,
		app.ComponentLogger("building_worker"),
	)

	worker.StartWorker(
		ctx,
		cfg.MissionWorker,
		missionService,
		app.ComponentLogger("mission_worker"),
	)

	worker.StartWorker(
		ctx,
		cfg.ResearchWorker,
		researchService,
		app.ComponentLogger("research_worker"),
	)

	worker.StartWorker(
		ctx,
		cfg.FleetConstructionWorker,
		fleetConstructionService,
		app.ComponentLogger("fleet_construction_worker"),
	)

	app.WaitShutdown(ctx)

	return nil
}
