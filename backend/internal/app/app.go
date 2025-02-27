package app

import (
	"log"

	"github.com/wojcikp/deps-dev-assignment/backend/internal/api"
	"github.com/wojcikp/deps-dev-assignment/backend/internal/database"
	dependenciesloader "github.com/wojcikp/deps-dev-assignment/backend/internal/dependencies_loader"
)

type App struct {
	dependenciesLoader *dependenciesloader.Loader
	db                 *database.SQLiteDB
	api                *api.Api
}

func NewApp(
	dependenciesLoader *dependenciesloader.Loader,
	db *database.SQLiteDB,
	api *api.Api,
) *App {
	return &App{dependenciesLoader, db, api}
}

func (app App) Run() {
	if err := app.db.CreateTables(); err != nil {
		log.Fatalf("failed to create db tables due to an error: %v \n exiting...", err)
	}

	if err := app.dependenciesLoader.FetchDepsDevDependencies(); err != nil {
		log.Fatalf("failed to fetch deps.dev dependencies due to an error: %v \n exiting...", err)
	}

	if err := app.db.LoadDependencies(app.dependenciesLoader.Dependencies.Nodes); err != nil {
		log.Fatalf("failed to load version keys into db due to an error: %v \n exiting...", err)
	}

	detailedDependencies := app.dependenciesLoader.FetchDetailsForAllDependencies()

	if err := app.db.LoadDetailedDependencies(detailedDependencies); err != nil {
		log.Fatalf("failed to load detailed dependencies into db due to an error: %v \n exiting...", err)
	}

	app.api.Run()
}
