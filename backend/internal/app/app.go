package app

import (
	"log"

	dependenciesloader "github.com/wojcikp/deps-dev-assignment/backend/internal/dependencies_loader"
)

type App struct {
	dependenciesLoader *dependenciesloader.Loader
}

func NewApp(
	dependenciesLoader *dependenciesloader.Loader,
) *App {
	return &App{dependenciesLoader}
}

func (app App) Run() {
	if err := app.dependenciesLoader.FetchDepsDevDependencies(); err != nil {
		log.Fatal("failed to fetch deps.dev dependencies")
	}

	detailedDependencies := app.dependenciesLoader.FetchDetailsForAllDependencies()
	log.Print(detailedDependencies)
}
