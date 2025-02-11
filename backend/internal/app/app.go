package app

import (
	"log"

	"github.com/wojcikp/deps-dev-assignment/backend/internal/database"
	dependenciesloader "github.com/wojcikp/deps-dev-assignment/backend/internal/dependencies_loader"
)

type App struct {
	dependenciesLoader *dependenciesloader.Loader
	db                 *database.SQLiteDB
}

func NewApp(
	dependenciesLoader *dependenciesloader.Loader,
	db *database.SQLiteDB,
) *App {
	return &App{dependenciesLoader, db}
}

func (app App) Run() {
	if err := app.db.CreateTables(); err != nil {
		log.Fatalf("failed to create db tables due to an error: %v \n exiting...", err)
	}

	if err := app.dependenciesLoader.FetchDepsDevDependencies(); err != nil {
		log.Fatalf("failed to fetch deps.dev dependencies due to an error: %v \n exiting...", err)
	}

	detailedDependencies := app.dependenciesLoader.FetchDetailsForAllDependencies()

	for _, d := range detailedDependencies {
		log.Print(d.Scorecard.Date)
	}

	if err := app.db.LoadDependencies(detailedDependencies); err != nil {
		log.Fatalf("failed to load detailed dependencies into db due to an error: %v \n exiting...", err)
	}
}
